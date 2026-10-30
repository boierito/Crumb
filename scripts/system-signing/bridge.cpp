#include "bundle.h"
#include "log.h"
#include <mutex>
#include <ctime>
#include <openssl/crypto.h>

static std::mutex signing_lock;

// Creates a local CMS signature. This does not assert that iOS has accepted
// installation, that an Apple identity is unrevoked, or that plaintext is clean.
extern "C" int CrumbSignPreparedBundle(const char *folder, const char *p12,
                                      const char *profile, const char *password,
                                      const char *expectedBundleID, const char *expectedTeamID) {
    if (!folder || !p12 || !profile || !password || !expectedBundleID || !expectedTeamID || !*expectedTeamID) return 1;
    std::lock_guard<std::mutex> lock(signing_lock);
    ZLog::SetLogLever(ZLog::E_NONE);
    std::string secret(password);
    int result=4;
    try {
        ZSignAsset identity;
        if (!identity.Init("",p12,profile,"",secret,false,false,false)) result=2;
        else {
            jvalue info,provision;
            std::string content;
            if (identity.m_strTeamId!=expectedTeamID || !info.read_plist_from_file("%s/Info.plist",folder) ||
                info["CFBundleIdentifier"].as_string()!=expectedBundleID ||
                !ZSignAsset::GetCMSContent(identity.m_strProvData,content) ||
                !provision.read_plist(content) || provision["ExpirationDate"].as_date()<=std::time(nullptr)) result=3;
            else {
                std::string exact=identity.m_strTeamId+"."+expectedBundleID;
                std::string wildcard=identity.m_strTeamId+".*";
                if (identity.m_strApplicationId!=exact && identity.m_strApplicationId!=wildcard) result=3;
                else {
                    ZBundle bundle;
                    const std::vector<std::string> empty;
                    result=bundle.SignFolder(&identity,folder,"","","",empty,empty,true,false,false,false)?0:4;
                }
            }
        }
    } catch (...) { result=4; }
    if (!secret.empty()) OPENSSL_cleanse(&secret[0],secret.size());
    return result;
}
