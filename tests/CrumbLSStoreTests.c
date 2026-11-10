#include "../SharedDecrypt/CrumbLSStore.h"
#include <assert.h>
#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
static void u32(unsigned char *b,size_t at,uint32_t x){memcpy(b+at,&x,4);}
static size_t string(unsigned char *b,size_t at,uint32_t id,const char *value){size_t n=strlen(value);u32(b,at,id);u32(b,at+4,(uint32_t)n);memcpy(b+at+8,value,n);return (at+8+n+4)&~(size_t)3;}
int main(void){
 unsigned char *b=calloc(1,16384);assert(b);memcpy(b,"bdsl",4);
 u32(b,260,0x250);u32(b,264,20<<2);u32(b,276,10<<2);
 const char *target="/private/var/containers/Bundle/Application/01234567-89AB-CDEF-0123-456789ABCDEF/Test.app";
 size_t at=string(b,1024,10,"test.selected");at=string(b,at,20,target);string(b,at,30,"test.unrelated");
 char paths[CRUMB_LS_MAX_PATHS][PATH_MAX];CrumbLSMetrics m;
 assert(crumb_ls_bundle_paths(b,16384,"test.selected",paths,&m)==1 && !strcmp(paths[0],target));
 assert(crumb_ls_bundle_paths(b,16384,"test.unrelated",paths,&m)==0);
 // An unrelated record cannot give its path to the selected bundle.
 u32(b,276,30<<2);assert(crumb_ls_bundle_paths(b,16384,"test.selected",paths,&m)==0);u32(b,276,10<<2);
 // Truncation, unknown format, oversized strings and traversal remain candidates-free.
 assert(crumb_ls_bundle_paths(b,200,"test.selected",paths,&m)==-1);
 b[0]='x';assert(crumb_ls_bundle_paths(b,16384,"test.selected",paths,&m)==-1);b[0]='b';
 memset(b+1024,0,15360);at=string(b,1024,10,"test.selected");
 string(b,at,20,"/private/var/containers/Bundle/Application/01234567-89AB-CDEF-0123-456789ABCDEF/../Escape.app");
 assert(crumb_ls_bundle_paths(b,16384,"test.selected",paths,&m)==0);
 // Multiple stale records are returned as distinct candidates; live identity
 // validation remains mandatory rather than trusting the first string.
 memset(b+1024,0,15360);at=string(b,1024,10,"test.selected");at=string(b,at,20,target);
 const char *other="/var/containers/Bundle/Application/FEDCBA98-7654-3210-FEDC-BA9876543210/Test.app";
 string(b,at,21,other);u32(b,900,0x250);u32(b,904,21<<2);u32(b,916,10<<2);
 assert(crumb_ls_bundle_paths(b,16384,"test.selected",paths,&m)==2);
 free(b);puts("LS store parser: target-only records, stale candidates, malformed data, path bounds and traversal checks passed (synthetic store)");
}
