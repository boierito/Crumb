#define _GNU_SOURCE
#include "../SharedDecrypt/CrumbCacheLease.h"
#include <assert.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>
static const char *first="C4A43AD1-BF62-4F81-8628-006A00870276";
static const char *second="9429AAC9-FF83-4496-98AD-E15A276CC274";
static void put(const char *path,const char *value){FILE *f=fopen(path,"w");assert(f);assert(fputs(value,f)>=0);fclose(f);}
static void has(const char *path,const char *value){char b[128]={0};FILE *f=fopen(path,"r");assert(f);assert(fread(b,1,sizeof(b)-1,f)==strlen(value));fclose(f);assert(!strcmp(b,value));}
static void append(char *out,size_t size,const char *base,const char *suffix){assert(snprintf(out,size,"%s/%s",base,suffix)<(int)size);}
int main(void){
 char temp[]="/tmp/crumb-lease-XXXXXX";assert(mkdtemp(temp));
 char root[PATH_MAX],library[PATH_MAX],cache[PATH_MAX],stage[PATH_MAX],other[PATH_MAX],target[PATH_MAX],file[PATH_MAX];
 append(root,sizeof(root),temp,"group");assert(mkdir(root,0700)==0);
 append(library,sizeof(library),root,"Library");assert(mkdir(library,0700)==0);
 append(cache,sizeof(cache),library,"Caches");assert(mkdir(cache,0750)==0);
 append(stage,sizeof(stage),temp,"private-link");append(other,sizeof(other),temp,"other-link");
 append(target,sizeof(target),temp,"private-target");assert(mkdir(target,0700)==0);
 append(file,sizeof(file),cache,"existing");
 CrumbCacheLease lease;struct stat original,restored;assert(stat(cache,&original)==0);
 // A normal transaction restores the original inode and permissions.
 assert(crumb_cache_begin(&lease,root,target,stage,first)==0);
 assert(stat(cache,&restored)==0 && restored.st_ino==original.st_ino);
 assert(crumb_cache_activate(&lease)==0 && lstat(cache,&restored)==0 && S_ISLNK(restored.st_mode));
 assert(crumb_cache_restore(&lease)==0 && stat(cache,&restored)==0 && restored.st_ino==original.st_ino && (restored.st_mode&0777)==(original.st_mode&0777));
 // Existing shared data is never removed or moved by a rejected preflight.
 put(file,"other-app-data");assert(crumb_cache_inspect(root)==-1 && errno==ENOTEMPTY);
 assert(crumb_cache_begin(&lease,root,target,stage,first)==-1);has(file,"other-app-data");assert(unlink(file)==0);
 // A concurrent writer after selection is preserved by rollback.
 assert(crumb_cache_begin(&lease,root,target,stage,first)==0);put(file,"arrived-after-selection");
 assert(crumb_cache_activate(&lease)==-1 && errno==ENOTEMPTY);assert(crumb_cache_restore(&lease)==0);has(file,"arrived-after-selection");assert(unlink(file)==0);
 // A writer holding an original directory descriptor keeps its data as well.
 int cachefd=open(cache,O_RDONLY|O_DIRECTORY);assert(cachefd>=0);
 assert(crumb_cache_begin(&lease,root,target,stage,first)==0 && crumb_cache_activate(&lease)==0);
 int newfile=openat(cachefd,"existing",O_WRONLY|O_CREAT|O_EXCL,0600);assert(newfile>=0);assert(write(newfile,"late",4)==4);close(newfile);close(cachefd);
 assert(crumb_cache_restore(&lease)==0);has(file,"late");assert(unlink(file)==0);
 // The ledger can be replayed after process death releases the lock.
 assert(crumb_cache_begin(&lease,root,target,stage,first)==0 && crumb_cache_activate(&lease)==0);
 CrumbCacheLease recovered=lease;close(lease.lock_fd);recovered.lock_fd=-1;
 assert(crumb_cache_restore(&recovered)==0 && stat(cache,&restored)==0 && restored.st_ino==original.st_ino);
 // Cooperating simultaneous jobs cannot both alter the shared Caches name.
 assert(crumb_cache_begin(&lease,root,target,stage,first)==0);
 CrumbCacheLease blocked;assert(crumb_cache_begin(&blocked,root,target,other,second)==-1);
 assert(crumb_cache_restore(&lease)==0);
 // A foreign replacement is retained, while the original backup waits for recovery.
 assert(crumb_cache_begin(&lease,root,target,stage,first)==0 && crumb_cache_activate(&lease)==0);
 assert(unlink(cache)==0 && mkdir(cache,0700)==0);put(file,"foreign-replacement");
 assert(crumb_cache_restore(&lease)==-1 && errno==EEXIST);has(file,"foreign-replacement");
 assert(unlink(file)==0 && rmdir(cache)==0);assert(crumb_cache_restore(&lease)==0);
 assert(stat(cache,&restored)==0 && restored.st_ino==original.st_ino);
 // An unrelated link and an existing private staging file are rejected intact.
 assert(rmdir(cache)==0 && symlink(target,cache)==0);
 assert(crumb_cache_begin(&lease,root,target,stage,first)==-1);assert(lstat(cache,&restored)==0 && S_ISLNK(restored.st_mode));assert(unlink(cache)==0);
 put(stage,"unrelated-private-file");assert(crumb_cache_begin(&lease,root,target,stage,first)==-1);has(stage,"unrelated-private-file");assert(unlink(stage)==0);
 // If Caches did not exist, restoring the lease does not manufacture a directory.
 assert(crumb_cache_begin(&lease,root,target,stage,first)==0 && !lease.had_caches);
 assert(crumb_cache_activate(&lease)==0 && crumb_cache_restore(&lease)==0);
 assert(lstat(cache,&restored)==-1 && errno==ENOENT);
 char lock[PATH_MAX];append(lock,sizeof(lock),root,".CrumbWorseQuery.lock");assert(unlink(lock)==0);
 assert(rmdir(library)==0 && rmdir(root)==0 && rmdir(target)==0 && rmdir(temp)==0);
 puts("cache lease: inode/mode preservation, occupied groups, concurrent writes, crash replay, locking, foreign replacements and non-overwrite passed");
}
