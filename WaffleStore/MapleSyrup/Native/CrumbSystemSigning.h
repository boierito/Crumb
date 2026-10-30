#ifndef CRUMB_SYSTEM_SIGNING_H
#define CRUMB_SYSTEM_SIGNING_H
#include <stddef.h>
int CrumbCurrentSigningTeam(char *output, size_t length);
int CrumbValidateSigningIdentity(const char *p12, const char *profile, const char *password,
                                 const char *expectedBundleID, const char *expectedTeamID);
int CrumbSignPreparedBundle(const char *folder, const char *p12,
                           const char *profile, const char *password,
                           const char *expectedBundleID, const char *expectedTeamID);
#endif
