#include "CrumbLSStore.h"
#include <stdbool.h>
#include <string.h>
#include <errno.h>
#include <stdio.h>
static uint32_t ls_u32(const unsigned char *b){uint32_t x;memcpy(&x,b,4);return x;}
static size_t ls_string(const unsigned char *b,size_t n,size_t p){
  if(n-p<9)return 0;uint32_t len=ls_u32(b+p+4);
  if(!len||len>=2048||len>=n-p-8||(len%4 && b[p+8+len])||memchr(b+p+8,0,len))return 0;
  return len;
}
static bool ls_has_id(const uint32_t *ids,size_t count,uint32_t value){
  for(size_t i=0;i<count;i++)if(ids[i]==value)return true;return false;
}
static bool ls_path(const unsigned char *s,size_t n,char *out){
  static const char *roots[]={"/private/var/containers/Bundle/Application/","/var/containers/Bundle/Application/"};
  size_t prefix=0;for(size_t i=0;i<2;i++){size_t k=strlen(roots[i]);if(n>k&&!memcmp(s,roots[i],k)){prefix=k;break;}}
  if(!prefix||n>=PATH_MAX||n-prefix<42)return false;
  const unsigned char *p=s+prefix;
  for(size_t i=0;i<36;i++){
    if(i==8||i==13||i==18||i==23){if(p[i]!='-')return false;}
    else if(!((p[i]>='0'&&p[i]<='9')||(p[i]>='a'&&p[i]<='f')||(p[i]>='A'&&p[i]<='F')))return false;
  }
  if(p[36]!='/')return false;
  size_t end=n;if(s[end-1]=='/')end--;
  if(end<=prefix+37+4||memcmp(s+end-4,".app",4))return false;
  for(size_t i=prefix+37;i<end;i++)if(s[i]=='/'||s[i]=='\\'||s[i]<32)return false;
  if(end-prefix-37<5)return false;
  memcpy(out,s,end);out[end]=0;return true;
}
int crumb_ls_bundle_paths(const void *data,size_t n,const char *bundle,
                          char out[CRUMB_LS_MAX_PATHS][PATH_MAX],CrumbLSMetrics *m){
  memset(m,0,sizeof(*m));memset(out,0,CRUMB_LS_MAX_PATHS*PATH_MAX);
  const unsigned char *b=data;size_t bundle_len=bundle?strlen(bundle):0;
  if(!b||n<0x258||n>64u*1024u*1024u||!bundle_len||bundle_len>255||memcmp(b,"bdsl",4)){m->error=EINVAL;return -1;}
  uint32_t ids[32]={0}, refs[32]={0};size_t id_count=0,ref_count=0;
  for(size_t p=0;p+8<n;p+=4){
    uint32_t sid=ls_u32(b+p);if(!sid||sid>=0x1000000)continue;
    size_t len=ls_string(b,n,p);if(len!=bundle_len||memcmp(b+p+8,bundle,len))continue;
    if(!ls_has_id(ids,id_count,sid)){if(id_count==32){m->error=E2BIG;return -1;}ids[id_count++]=sid;}
  }
  for(size_t p=0;p+8+0x250<=n;p+=4){
    if(ls_u32(b+p+4)!=0x250)continue;m->records++;
    if(m->records>100000){m->error=E2BIG;return -1;}
    uint32_t id=ls_u32(b+p+8+0x0c),path=ls_u32(b+p+8);
    if(!id||(id&3)||!ls_has_id(ids,id_count,id>>2))continue;
    m->target_records++;
    if(!path||(path&3)||(path>>2)>=0x1000000)continue;
    if(!ls_has_id(refs,ref_count,path>>2)){if(ref_count==32){m->error=E2BIG;return -1;}refs[ref_count++]=path>>2;}
  }
  // Overlapping table IDs are accepted only when the string itself has the
  // installed-bundle shape. Unrelated table strings never become file paths.
  for(size_t p=0;p+8<n;p+=4){
    uint32_t sid=ls_u32(b+p);if(!ls_has_id(refs,ref_count,sid))continue;
    size_t len=ls_string(b,n,p);char candidate[PATH_MAX];
    if(!len||!ls_path(b+p+8,len,candidate))continue;
    bool found=false;for(size_t i=0;i<m->path_count;i++)if(!strcmp(out[i],candidate))found=true;
    if(found)continue;
    if(m->path_count==CRUMB_LS_MAX_PATHS){memset(out,0,CRUMB_LS_MAX_PATHS*PATH_MAX);m->error=E2BIG;return -1;}
    snprintf(out[m->path_count++],PATH_MAX,"%s",candidate);
  }
  return (int)m->path_count;
}
