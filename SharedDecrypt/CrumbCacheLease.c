#define _GNU_SOURCE
#include "CrumbCacheLease.h"
#include <dirent.h>
#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <string.h>
#include <sys/file.h>
#include <sys/stat.h>
#include <unistd.h>
#ifdef __linux__
#include <sys/syscall.h>
#endif
static int lease_paths(const CrumbCacheLease *l,char *caches,char *backup,char *quarantine) {
  if(strlen(l->identifier)!=36 || !l->root[0] || !l->staging[0] || !l->target[0]){errno=EINVAL;return -1;}
  for(size_t i=0;i<36;i++) {
    char c=l->identifier[i];
    if(!((c>='0'&&c<='9')||(c>='a'&&c<='f')||(c>='A'&&c<='F')||c=='-')){errno=EINVAL;return -1;}
  }
  if(snprintf(caches,PATH_MAX,"%s/Library/Caches",l->root)>=PATH_MAX ||
     snprintf(backup,PATH_MAX,"%s/Library/.CrumbProbe-Caches-%s",l->root,l->identifier)>=PATH_MAX ||
     snprintf(quarantine,PATH_MAX,"%s.quarantine",l->staging)>=PATH_MAX){errno=ENAMETOOLONG;return -1;}
  return 0;
}
static int exclusive_move(const char *from,const char *to) {
#ifdef __APPLE__
  return renamex_np(from,to,RENAME_EXCL);
#elif defined(__linux__)
  return (int)syscall(SYS_renameat2,AT_FDCWD,from,AT_FDCWD,to,1 /* RENAME_NOREPLACE */);
#else
  (void)from;(void)to;errno=ENOTSUP;return -1;
#endif
}
static int empty_directory(const char *path) {
  int fd=open(path,O_RDONLY|O_DIRECTORY|O_NOFOLLOW|O_CLOEXEC);
  if(fd<0)return -1;
  DIR *directory=fdopendir(fd);if(!directory){close(fd);return -1;}
  int code=0;struct dirent *entry;errno=0;
  while((entry=readdir(directory))) {
    if(strcmp(entry->d_name,".") && strcmp(entry->d_name,"..")){code=ENOTEMPTY;break;}
    errno=0;
  }
  if(!code && errno)code=errno;
  closedir(directory);if(code){errno=code;return -1;}return 0;
}
int crumb_cache_inspect(const char *root) {
  char library[PATH_MAX],caches[PATH_MAX];struct stat st;
  if(snprintf(library,sizeof(library),"%s/Library",root)>=PATH_MAX ||
     snprintf(caches,sizeof(caches),"%s/Library/Caches",root)>=PATH_MAX){errno=ENAMETOOLONG;return -1;}
  if(lstat(root,&st))return -1;
  if(!S_ISDIR(st.st_mode)){errno=ENOTDIR;return -1;}
  if(lstat(library,&st)){return errno==ENOENT?0:-1;}
  if(!S_ISDIR(st.st_mode)){errno=ELOOP;return -1;}
  if(lstat(caches,&st)){return errno==ENOENT?0:-1;}
  if(!S_ISDIR(st.st_mode)){errno=ELOOP;return -1;}
  return empty_directory(caches);
}
static int take_lock(CrumbCacheLease *l) {
  if(l->lock_fd>=0)return 0;
  char path[PATH_MAX];if(snprintf(path,sizeof(path),"%s/.CrumbWorseQuery.lock",l->root)>=PATH_MAX){errno=ENAMETOOLONG;return -1;}
  int fd=open(path,O_RDWR|O_CREAT|O_NOFOLLOW|O_CLOEXEC,0600);if(fd<0)return -1;
  struct stat st;if(fstat(fd,&st)||!S_ISREG(st.st_mode)){close(fd);errno=EINVAL;return -1;}
  if(flock(fd,LOCK_EX|LOCK_NB)){int code=errno;close(fd);errno=code;return -1;}
  l->lock_fd=fd;return 0;
}
static void release_lock(CrumbCacheLease *l){if(l->lock_fd>=0){close(l->lock_fd);l->lock_fd=-1;}}
static bool same_inode(const struct stat *st,uint64_t device,uint64_t inode) {
  return (uint64_t)st->st_dev==device && (uint64_t)st->st_ino==inode;
}
static bool owned_link(const CrumbCacheLease *l,const char *path) {
  struct stat st;if(lstat(path,&st)||!S_ISLNK(st.st_mode)||!same_inode(&st,l->link_device,l->link_inode))return false;
  char value[PATH_MAX];ssize_t n=readlink(path,value,sizeof(value)-1);if(n<0)return false;value[n]=0;
  return !strcmp(value,l->target);
}
static int remove_own_link(const CrumbCacheLease *l,const char *path) {
  struct stat st;if(lstat(path,&st))return errno==ENOENT?0:-1;
  if(!owned_link(l,path)){errno=ESTALE;return -1;}return unlink(path);
}
int crumb_cache_begin(CrumbCacheLease *l,const char *root,const char *target,const char *staging,const char *identifier) {
  memset(l,0,sizeof(*l));l->lock_fd=-1;
  if(strlen(root)>=sizeof(l->root)||strlen(target)>=sizeof(l->target)||strlen(staging)>=sizeof(l->staging)||strlen(identifier)>=sizeof(l->identifier)){errno=ENAMETOOLONG;return -1;}
  strcpy(l->root,root);strcpy(l->target,target);strcpy(l->staging,staging);strcpy(l->identifier,identifier);
  char caches[PATH_MAX],backup[PATH_MAX],quarantine[PATH_MAX],library[PATH_MAX];
  if(lease_paths(l,caches,backup,quarantine)||take_lock(l))return -1;
  int code=0;struct stat st;
  if(crumb_cache_inspect(root)){code=errno;goto failed;}
  if(snprintf(library,sizeof(library),"%s/Library",root)>=PATH_MAX){code=ENAMETOOLONG;goto failed;}
  if(mkdir(library,0700) && errno!=EEXIST){code=errno;goto failed;}
  if(lstat(caches,&st)==0){l->had_caches=true;l->original_device=st.st_dev;l->original_inode=st.st_ino;}
  else if(errno!=ENOENT){code=errno;goto failed;}
  if(symlink(target,staging)){code=errno;goto failed;}
  if(lstat(staging,&st)){code=errno;unlink(staging);goto failed;}
  l->link_device=st.st_dev;l->link_inode=st.st_ino;return 0;
failed:release_lock(l);errno=code;return -1;
}
int crumb_cache_activate(CrumbCacheLease *l) {
  char caches[PATH_MAX],backup[PATH_MAX],quarantine[PATH_MAX];struct stat st;
  if(lease_paths(l,caches,backup,quarantine)||l->lock_fd<0){errno=EINVAL;return -1;}
  if(l->had_caches){
    if(exclusive_move(caches,backup))return -1;
    if(lstat(backup,&st)||!same_inode(&st,l->original_device,l->original_inode)){errno=ESTALE;return -1;}
    // A file may have arrived since selection. Preserve it and decline the query.
    if(empty_directory(backup))return -1;
  }
  return exclusive_move(l->staging,caches);
}
int crumb_cache_restore(CrumbCacheLease *l) {
  char caches[PATH_MAX],backup[PATH_MAX],quarantine[PATH_MAX];struct stat st;int code=0;
  if(lease_paths(l,caches,backup,quarantine)||take_lock(l))return -1;
  if(owned_link(l,caches)) {
    if(exclusive_move(caches,quarantine)){code=errno;goto done;}
    if(remove_own_link(l,quarantine)){code=errno;goto done;}
  }
  if(remove_own_link(l,l->staging)){code=errno;goto done;}
  if(remove_own_link(l,quarantine)){code=errno;goto done;}
  if(lstat(backup,&st)==0) {
    if(!l->had_caches||!same_inode(&st,l->original_device,l->original_inode)){code=ESTALE;goto done;}
    if(exclusive_move(backup,caches)){code=errno;goto done;}
  }else if(errno!=ENOENT){code=errno;goto done;}
done:release_lock(l);if(code){errno=code;return -1;}return 0;
}
