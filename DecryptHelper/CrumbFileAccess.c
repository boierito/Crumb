#include "CrumbFileAccess.h"
#include <errno.h>
#include <fcntl.h>
#include <string.h>
#include <sys/stat.h>
#include <unistd.h>

int crumb_open_installed_image(const char *path, CrumbFileAccess *status) {
  memset(status, 0, sizeof(*status));
  int fd;
  do {
    fd = open(path, O_RDONLY | O_NOFOLLOW | O_CLOEXEC);
  } while (fd < 0 && errno == EINTR);
  int saved = fd < 0 ? errno : 0;
  struct stat st;
  int inspected = fd < 0 ? lstat(path, &st) : fstat(fd, &st);
  if (inspected == 0) {
    status->metadata_available = 1;
    status->mode = (uint32_t)st.st_mode;
    status->owner = (uint32_t)st.st_uid;
    status->group = (uint32_t)st.st_gid;
  }
  status->error_code = saved;
  errno = saved;
  return fd;
}
