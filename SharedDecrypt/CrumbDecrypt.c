#include "CrumbDecrypt.h"
#include "../DecryptHelper/CrumbMachO.h"
#include <errno.h>
#include <fcntl.h>
#include <math.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <sys/stat.h>
#include <unistd.h>
#ifdef __APPLE__
#include <dlfcn.h>
#include <mach/mach.h>
#endif
static bool checkpoint(CrumbDecryptMetrics *m, const CrumbDecryptHooks *h, int stage) {
  m->stage=stage;
  if(h && h->checkpoint) h->checkpoint(h->context,m);
  m->cancelled=h && h->cancelled && h->cancelled(h->context);
  if(m->cancelled) m->error_code=ECANCELED;
  return !m->cancelled;
}
static int native_remap(void *ctx,void *address,size_t length,uint32_t cpu,uint32_t subtype) {
  (void)ctx;
#if defined(__APPLE__) && defined(__arm64__)
  typedef int (*Remap)(void *,size_t,uint32_t,uint32_t,uint32_t);
  Remap fn=(Remap)dlsym(RTLD_DEFAULT,"mremap_encrypted");
  if(fn) return fn(address,length,2,cpu,subtype);
  register uint64_t x0 __asm__("x0")=(uint64_t)address;
  register uint64_t x1 __asm__("x1")=length;
  register uint64_t x2 __asm__("x2")=2;
  register uint64_t x3 __asm__("x3")=cpu;
  register uint64_t x4 __asm__("x4")=subtype;
  register uint64_t x16 __asm__("x16")=489;
  uint64_t failed;
  __asm__ volatile("svc #0x80\n\tcset %[failed], cs" : "+r"(x0),"+r"(x1),[failed] "=r"(failed)
                   : "r"(x2),"r"(x3),"r"(x4),"r"(x16) : "memory","cc");
  if(failed) { errno=(int)x0; return -1; } return (int)x0;
#else
  (void)address;(void)length;(void)cpu;(void)subtype;errno=ENOTSUP;return -1;
#endif
}
static int native_read(void *ctx,const void *address,void *out,size_t len,size_t *got) {
  (void)ctx;
#ifdef __APPLE__
  vm_size_t count=0;
  int kr=vm_read_overwrite(mach_task_self(),(vm_address_t)address,len,(vm_address_t)out,&count);
  *got=count;return kr;
#else
  (void)address;(void)out;(void)len;*got=0;return ENOTSUP;
#endif
}
static int native_signature(void *ctx,int fd,const uint8_t *map,size_t size) {
  (void)ctx;
#ifdef __APPLE__
  uint32_t count;memcpy(&count,map+16,4);size_t pos=32;
  for(uint32_t i=0;i<count;i++) {
    uint32_t command,length;memcpy(&command,map+pos,4);memcpy(&length,map+pos+4,4);
    if(command==0x1d && length>=16) {
      uint32_t offset,bytes;memcpy(&offset,map+pos+8,4);memcpy(&bytes,map+pos+12,4);
      if(offset>size || bytes>size-offset || !bytes) {errno=EINVAL;return -1;}
      fsignatures_t sig={.fs_file_start=0,.fs_blob_start=(void *)(uintptr_t)offset,.fs_blob_size=bytes};
      return fcntl(fd,F_ADDFILESIGS_RETURN,&sig);
    }
    pos+=length;
  }
#else
  (void)fd;(void)map;(void)size;
#endif
  errno=ENOENT;return -1;
}
static uint64_t words(const uint8_t *b,size_t n) {
  uint64_t count=0;for(size_t i=0;i+4<=n;i+=4) {
    uint32_t w;memcpy(&w,b+i,4);
    count+=w==0xd65f03c0 || w==0xd503237f || w==0xd65f0fff;
  }return count;
}
static double entropy(const uint64_t *h,size_t n) {
  double e=0;for(int i=0;i<256;i++) if(h[i]) {double p=(double)h[i]/n;e-=p*log2(p);}return e;
}
static int full_pwrite(int fd,const void *buf,size_t n,off_t off) {
  size_t p=0;while(p<n) {ssize_t w=pwrite(fd,(const uint8_t *)buf+p,n-p,off+p);
    if(w<0 && errno==EINTR) continue;if(w<=0) {if(!w)errno=EIO;return -1;}p+=(size_t)w;}return 0;
}
int crumb_decrypt_image(const char *src,const char *dst,CrumbDecryptMetrics *m,const CrumbDecryptHooks *h) {
  memset(m,0,sizeof(*m));m->remap_result=-1;m->signature_result=-1;m->first_remap_result=-1;
  int fd=-1,out=-1,result=-1;bool created=false;uint8_t *map=MAP_FAILED,*plain=NULL,*cipher=NULL;size_t size=0;
  if(!checkpoint(m,h,CRUMB_OPEN))goto done;
  fd=open(src,O_RDONLY|O_NOFOLLOW|O_CLOEXEC);if(fd<0){m->error_code=errno;goto done;}
  if(!checkpoint(m,h,CRUMB_INSPECT))goto done;
  struct stat st;if(fstat(fd,&st)) {m->error_code=errno;goto done;}
  if(!S_ISREG(st.st_mode)||st.st_size<0||st.st_size>512LL*1024*1024){m->error_code=EFBIG;goto done;}
  size=(size_t)st.st_size;
  if(size) {
    if(!checkpoint(m,h,CRUMB_MAP))goto done;
    map=mmap(NULL,size,PROT_READ,MAP_PRIVATE,fd,0);if(map==MAP_FAILED){m->error_code=errno;goto done;}
  }
  if(!checkpoint(m,h,CRUMB_PARSE))goto done;
  CrumbEncryption enc={0};int kind=size?crumb_macho(map,size,&enc):0;
  if(kind<0){m->error_code=ENOEXEC;goto done;}m->encrypted=kind==1;
  size_t pg=(size_t)getpagesize();
  if(kind && enc.offset%pg){m->error_code=EINVAL;goto done;}
  if(!checkpoint(m,h,CRUMB_COPY))goto done;
  out=open(dst,O_WRONLY|O_CREAT|O_EXCL|O_NOFOLLOW|O_CLOEXEC,0700);
  if(out<0){m->error_code=errno;goto done;}created=true;
  for(size_t p=0;p<size;) {
    if(h && h->cancelled && h->cancelled(h->context)){m->cancelled=1;m->error_code=ECANCELED;goto done;}
    size_t len=size-p<65536?size-p:65536;
    if(full_pwrite(out,map+p,len,(off_t)p)){m->write_errno=errno;m->error_code=errno;goto done;}p+=len;
  }
  if(!kind){if(fchmod(out,st.st_mode & 0777)||fsync(out)){m->error_code=errno;goto done;}result=0;goto done;}
  m->pages_expected=(enc.length+pg-1)/pg;
  // Signature registration is a measured retry, not a prerequisite. This keeps
  // the no-registration path that has documented successful installed controls.
  for(int attempt=0;attempt<2;attempt++) {
    if(attempt) {
      if(!checkpoint(m,h,CRUMB_SIGNATURE))goto done;
      m->signature_attempted=true;errno=0;
      m->signature_result=(h && h->signature?h->signature:native_signature)(h?h->context:NULL,fd,map,size);
      m->signature_errno=m->signature_result?errno:0;
      if(m->signature_result)break;
      munmap(map,size);map=mmap(NULL,size,PROT_READ,MAP_PRIVATE,fd,0);
      if(map==MAP_FAILED){m->error_code=errno;goto done;}
    }
    if(!checkpoint(m,h,CRUMB_REMAP))goto done;
    m->remap_attempted=true;errno=0;
    m->remap_result=(h && h->remap?h->remap:native_remap)(h?h->context:NULL,map+enc.offset,enc.length,enc.cpu,enc.subtype);
    m->remap_errno=m->remap_result?errno:0;
    m->remap_attempts++;
    if(attempt==0){m->first_remap_result=m->remap_result;m->first_remap_errno=m->remap_errno;}
    if(!m->remap_result)break;
  }
  if(m->remap_result){m->error_code=m->remap_errno;goto done;}
  plain=malloc(pg);cipher=malloc(pg);if(!plain||!cipher){m->error_code=ENOMEM;goto done;}
  uint64_t before[256]={0},after[256]={0};
  for(size_t p=0;p<enc.length;p+=pg) {
    if(!checkpoint(m,h,CRUMB_READ))goto done;
    size_t len=enc.length-p<pg?enc.length-p:pg,got=0;
    m->page_result=(h && h->read_page?h->read_page:native_read)(h?h->context:NULL,map+enc.offset+p,plain,len,&got);
    if(m->page_result||got!=len){m->first_fail_offset=p;m->error_code=EIO;goto done;}
    errno=0;ssize_t n=pread(fd,cipher,len,enc.offset+p);
    if(n!=(ssize_t)len){m->read_errno=n<0?errno:EIO;m->error_code=m->read_errno;m->first_fail_offset=p;goto done;}
    m->pages_read++;m->changed_pages+=memcmp(cipher,plain,len)!=0;
    m->oldwords+=words(cipher,len);m->newwords+=words(plain,len);
    for(size_t i=0;i<len;i++){before[cipher[i]]++;after[plain[i]]++;}
    if(full_pwrite(out,plain,len,enc.offset+p)){m->write_errno=errno;m->error_code=errno;goto done;}
  }
  if(!checkpoint(m,h,CRUMB_VERIFY))goto done;
  m->entropy_before=entropy(before,enc.length);m->entropy_after=entropy(after,enc.length);
  if(!m->changed_pages||m->entropy_before<=7.9||m->entropy_after>=7.5||m->newwords<=m->oldwords){m->error_code=EBADMSG;goto done;}
  if(!checkpoint(m,h,CRUMB_PATCH))goto done;
  uint32_t zero=0;if(full_pwrite(out,&zero,4,(off_t)enc.command)||fchmod(out,st.st_mode & 0777)||fsync(out)){m->error_code=errno;goto done;}
  result=1;
done:
  free(plain);free(cipher);if(map!=MAP_FAILED)munmap(map,size);if(fd>=0)close(fd);if(out>=0)close(out);
  if(result<0 && created)unlink(dst);
  if(result>=0 && !checkpoint(m,h,CRUMB_COMPLETE)){if(created)unlink(dst);result=-1;}
  return result;
}
