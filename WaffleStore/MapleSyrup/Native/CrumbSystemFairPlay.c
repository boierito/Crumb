#if defined(__APPLE__)
#define _DARWIN_C_SOURCE 1
#else
#define _POSIX_C_SOURCE 200809L
#endif
#include "CrumbSystemFairPlay.h"
#include <errno.h>
#include <fcntl.h>
#include <limits.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>
#if defined(__APPLE__) && defined(__aarch64__)
#include <dlfcn.h>
#include <CoreFoundation/CoreFoundation.h>
#include <mach/mach.h>
#include <sys/mman.h>

// Read only entitlements already granted to THIS process by iOS. This never
// fabricates a group identifier or asks for another application's container.
static CFTypeRef self_entitlement(CFStringRef name) {
    typedef CFTypeRef (*CreateTask)(CFAllocatorRef);
    typedef CFTypeRef (*ReadEntitlement)(CFTypeRef,CFStringRef,CFErrorRef*);
    CreateTask create=(CreateTask)dlsym(RTLD_DEFAULT,"SecTaskCreateFromSelf");
    ReadEntitlement read=(ReadEntitlement)dlsym(RTLD_DEFAULT,"SecTaskCopyValueForEntitlement");
    if (!create || !read) return NULL;
    CFTypeRef task=create(kCFAllocatorDefault);if(!task)return NULL;
    CFTypeRef result=read(task,name,NULL);CFRelease(task);return result;
}
static uint32_t little32(const unsigned char *p) {
    return (uint32_t)p[0] | (uint32_t)p[1]<<8 | (uint32_t)p[2]<<16 | (uint32_t)p[3]<<24;
}
static int transfer(int fd, void *buffer, size_t length, uint64_t offset, int writing) {
    unsigned char *p = buffer;
    while (length) {
        ssize_t n = writing ? pwrite(fd,p,length,(off_t)offset) : pread(fd,p,length,(off_t)offset);
        if (n < 0 && errno == EINTR) continue;
        if (n <= 0) return -1;
        p += n; offset += (uint64_t)n; length -= (size_t)n;
    }
    return 0;
}
static int native_remap(void *address, size_t size, uint32_t cpu, uint32_t subtype) {
    typedef int (*Remap)(void *,size_t,uint32_t,uint32_t,uint32_t);
    Remap remap = (Remap)dlsym(RTLD_DEFAULT,"mremap_encrypted");
    if (remap) return remap(address,size,2,cpu,subtype);
    register uint64_t a __asm__("x0") = (uint64_t)address;
    register uint64_t b __asm__("x1") = size;
    register uint64_t c __asm__("x2") = 2;
    register uint64_t d __asm__("x3") = cpu;
    register uint64_t e __asm__("x4") = subtype;
    register uint64_t number __asm__("x16") = 489;
    uint64_t failed;
    __asm__ volatile("svc #0x80\n\tcset %w[failed], cs"
        : "+r"(a), [failed] "=r"(failed)
        : "r"(b),"r"(c),"r"(d),"r"(e),"r"(number) : "memory","cc");
    if (failed) { errno = (int)a; return -1; }
    return (int)a;
}
#endif

int CrumbCurrentSigningTeam(char *output,size_t length) {
    if (!output || length<2) return -1;
    output[0]=0;
#if defined(__APPLE__) && defined(__aarch64__)
    CFTypeRef team=self_entitlement(CFSTR("com.apple.developer.team-identifier"));
    if(!team)return -1;
    int ok=CFGetTypeID(team)==CFStringGetTypeID()&&CFStringGetCString((CFStringRef)team,output,(CFIndex)length,kCFStringEncodingUTF8);
    CFRelease(team);return ok&&output[0]?0:-1;
#else
    return -1;
#endif
}

int CrumbCurrentAppGroup(size_t index,char *output,size_t length) {
    if (!output || length<2 || index>=32) return -1;
    output[0]=0;
#if defined(__APPLE__) && defined(__aarch64__)
    CFTypeRef groups=self_entitlement(CFSTR("com.apple.security.application-groups"));
    if(!groups)return 1;
    if(CFGetTypeID(groups)!=CFArrayGetTypeID()){CFRelease(groups);return -1;}
    if(index>=(size_t)CFArrayGetCount((CFArrayRef)groups)){CFRelease(groups);return 1;}
    CFTypeRef group=CFArrayGetValueAtIndex((CFArrayRef)groups,(CFIndex)index);
    int ok=group && CFGetTypeID(group)==CFStringGetTypeID() &&
           CFStringGetCString((CFStringRef)group,output,(CFIndex)length,kCFStringEncodingUTF8);
    CFRelease(groups);return ok&&output[0]?0:-1;
#else
    (void)index;return 1;
#endif
}

int CrumbSystemDecryptRange(const char *source,const char *destination,
                           CrumbSystemRange range,CrumbSystemRangeResult *result) {
    if (!result) return -1;
    memset(result,0,sizeof(*result));result->remap_result = -1;
#if !defined(__APPLE__) || !defined(__aarch64__)
    (void)source;(void)destination;(void)range;
    result->remap_errno = ENOTSUP;return -1;
#else
    int input = -1, output = -1, status = -1;
    unsigned char *mapping = MAP_FAILED, *plain = NULL, *cipher = NULL;
    size_t mapped_size = 0;
    long page_value = sysconf(_SC_PAGESIZE);
    if (!source || !destination || page_value<=0) { result->remap_errno=EINVAL;return -1; }
    uint64_t page = (uint64_t)page_value;
    if (!range.crypt_size || range.cpu_type!=0x100000c || range.slice_offset%page || range.crypt_offset%page ||
        range.crypt_offset>range.slice_size || range.crypt_size>range.slice_size-range.crypt_offset ||
        range.slice_size>INT64_MAX || range.slice_offset>INT64_MAX-range.slice_size || range.cryptid_field<range.slice_offset+48 ||
        range.cryptid_field>range.slice_offset+range.slice_size-4 || range.slice_size>SIZE_MAX-page) {
        result->remap_errno=EINVAL;return -1;
    }
    input=open(source,O_RDONLY|O_NOFOLLOW);output=open(destination,O_RDWR|O_NOFOLLOW);
    if(input<0||output<0){result->remap_errno=errno;goto done;}
    struct stat original,copy;
    if(fstat(input,&original)||fstat(output,&copy)){result->remap_errno=errno;goto done;}
    if(!S_ISREG(original.st_mode)||!S_ISREG(copy.st_mode)||original.st_size!=copy.st_size||
       (original.st_dev==copy.st_dev&&original.st_ino==copy.st_ino)||
       range.slice_offset+range.slice_size>(uint64_t)original.st_size){result->remap_errno=EINVAL;goto done;}
    unsigned char header[32],command[24];
    if(transfer(input,header,sizeof(header),range.slice_offset,0)||transfer(input,command,sizeof(command),range.cryptid_field-16,0)){
        result->read_error=EIO;goto done;
    }
    int encryption_command=(little32(command)==0x2c&&little32(command+4)==24)||
                           (little32(command)==0x21&&little32(command+4)==20);
    if(little32(header)!=0xfeedfacf||little32(header+4)!=range.cpu_type||little32(header+8)!=range.cpu_subtype||
       !encryption_command||little32(command+8)!=range.crypt_offset||
       little32(command+12)!=range.crypt_size||!little32(command+16)){
        result->remap_errno=EINVAL;goto done;
    }
    mapped_size=(size_t)(range.crypt_offset+range.crypt_size+page-1)/page*(size_t)page;
    result->mapping_attempted=1;
    mapping=mmap(NULL,mapped_size,PROT_READ,MAP_PRIVATE,input,(off_t)range.slice_offset);
    if(mapping==MAP_FAILED){result->mapping_errno=errno;result->remap_errno=errno;goto done;}
    result->remap_attempted=1;
    errno=0;result->remap_result=native_remap(mapping+range.crypt_offset,(size_t)range.crypt_size,range.cpu_type,range.cpu_subtype);
    result->remap_errno=result->remap_result==0?0:errno;
    if(result->remap_result!=0)goto done;
    result->pages_total=(range.crypt_size+page-1)/page;
    plain=malloc((size_t)page);cipher=malloc((size_t)page);
    if(!plain||!cipher){result->read_error=ENOMEM;goto done;}
    for(uint64_t offset=0;offset<range.crypt_size;offset+=page){
        size_t count=(size_t)(range.crypt_size-offset<page?range.crypt_size-offset:page);
        vm_size_t read=0;
        kern_return_t kr=vm_read_overwrite(mach_task_self(),(vm_address_t)(mapping+range.crypt_offset+offset),
                                          (vm_size_t)count,(vm_address_t)plain,&read);
        if(kr!=KERN_SUCCESS||read!=count){result->read_error=kr==KERN_SUCCESS?EIO:kr;goto done;}
        result->pages_read++;
        if(transfer(input,cipher,count,range.slice_offset+range.crypt_offset+offset,0)){result->read_error=EIO;goto done;}
        if(!memcmp(plain,cipher,count)){result->read_error=EINVAL;goto done;}
        result->pages_changed++;
        if(transfer(output,plain,count,range.slice_offset+range.crypt_offset+offset,1)){result->read_error=EIO;goto done;}
    }
    unsigned char zero[4]={0};
    if(transfer(output,zero,4,range.cryptid_field,1)||fsync(output)){result->read_error=EIO;goto done;}
    status=0;
done:
    if(plain){volatile unsigned char *v=plain;for(size_t i=0;i<(size_t)page;i++)v[i]=0;free(plain);}
    if(cipher){volatile unsigned char *v=cipher;for(size_t i=0;i<(size_t)page;i++)v[i]=0;free(cipher);}
    if(mapping!=MAP_FAILED)munmap(mapping,mapped_size);
    if(input>=0)close(input);if(output>=0)close(output);
    return status;
#endif
}
