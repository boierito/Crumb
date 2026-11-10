#include "../SharedDecrypt/CrumbDecrypt.h"
#include <assert.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

typedef struct {int remaps,reads,fail_page,fail_remap,cancel_after,retry;bool nochange;size_t pg;const char *source;} Fixture;
static int remap(void *p,void *addr,size_t n,uint32_t cpu,uint32_t sub){
  (void)addr;(void)n;(void)cpu;(void)sub;Fixture *f=p;f->remaps++;
  if(f->fail_remap || (f->retry && f->remaps==1)){errno=EPERM;return -1;}return 0;
}
static int signature(void *p,int fd,const uint8_t *b,size_t n){(void)p;(void)fd;(void)b;(void)n;return 0;}
static int read_page(void *p,const void *addr,void *out,size_t n,size_t *got){
  Fixture *f=p;f->reads++;if(f->reads==f->fail_page){*got=0;return 5;}
  if(f->nochange)memcpy(out,addr,n);
  else for(size_t i=0;i<n;i+=4){uint32_t ret=0xd65f03c0;memcpy((uint8_t *)out+i,&ret,n-i<4?n-i:4);}
  *got=n;return 0;
}
static bool cancelled(void *p){Fixture *f=p;return f->cancel_after && f->reads>=f->cancel_after;}
static void set32(uint8_t *p,size_t offset,uint32_t x){memcpy(p+offset,&x,4);}
static void make_image(const char *path,size_t pg){
  size_t size=pg*3;uint8_t *b=calloc(1,size);assert(b);
  set32(b,0,0xfeedfacf);set32(b,4,0x0100000c);set32(b,16,1);set32(b,20,24);
  set32(b,32,0x2c);set32(b,36,24);set32(b,40,(uint32_t)pg);set32(b,44,(uint32_t)(2*pg));set32(b,48,1);
  for(size_t i=pg;i<size;i++)b[i]=(uint8_t)i;
  FILE *file=fopen(path,"wb");assert(file);assert(fwrite(b,1,size,file)==size);fclose(file);free(b);assert(chmod(path,0755)==0);
}
int main(void){
  char directory[]="/tmp/crumb-reader-XXXXXX";assert(mkdtemp(directory));char src[512],dst[512],link[512];
  snprintf(src,sizeof(src),"%s/source",directory);snprintf(dst,sizeof(dst),"%s/output",directory);snprintf(link,sizeof(link),"%s/link",directory);
  size_t pg=(size_t)getpagesize();make_image(src,pg);
  CrumbDecryptMetrics m;
  for(int mode=0;mode<6;mode++){
    Fixture f={.pg=pg,.fail_page=mode==1?2:0,.fail_remap=mode==2,.cancel_after=mode==3?1:0,.nochange=mode==4,.retry=mode==5};
    CrumbDecryptHooks h={.context=&f,.cancelled=cancelled,.remap=remap,.read_page=read_page,.signature=signature};
    int result=crumb_decrypt_image(src,dst,&m,&h);
    if(mode==0 || mode==5){
      assert(result==1 && m.pages_read==2 && m.pages_expected==2 && m.changed_pages==2);
      struct stat attributes;assert(stat(dst,&attributes)==0 && (attributes.st_mode & 0777)==0755);
      int fd=open(dst,O_RDONLY);assert(fd>=0);uint32_t cid=1;assert(pread(fd,&cid,4,48)==4 && cid==0);close(fd);unlink(dst);
      if(mode==5)assert(f.remaps==2 && m.signature_attempted);
    }else{assert(result==-1 && access(dst,F_OK)!=0);if(mode==1)assert(m.page_result==5 && m.pages_read==1);}
    int fd=open(src,O_RDONLY);uint32_t cid=0;assert(pread(fd,&cid,4,48)==4 && cid==1);close(fd);
  }
  assert(symlink(src,link)==0);assert(crumb_decrypt_image(link,dst,&m,NULL)==-1 && m.stage==CRUMB_OPEN);
  unlink(link);
  FILE *file=fopen(dst,"wb");assert(file);fputs("existing",file);fclose(file);
  assert(crumb_decrypt_image(src,dst,&m,NULL)==-1 && m.error_code==EEXIST);
  char bytes[9]={0};file=fopen(dst,"rb");assert(fread(bytes,1,8,file)==8);fclose(file);assert(!strcmp(bytes,"existing"));unlink(dst);
  file=fopen(src,"wb");fputs("resource",file);fclose(file);
  assert(crumb_decrypt_image(src,dst,&m,NULL)==0 && !m.remap_attempted);unlink(dst);
  make_image(src,pg);int fd=open(src,O_WRONLY);uint32_t invalid=0xffffffff;assert(pwrite(fd,&invalid,4,40)==4);close(fd);
  assert(crumb_decrypt_image(src,dst,&m,NULL)==-1 && m.stage==CRUMB_PARSE);
  unlink(src);rmdir(directory);puts("shared reader: success, retry, partial read, remap failure, cancellation, unchanged bytes, symlink, non-overwrite and malformed input passed (synthetic callbacks)");
}
