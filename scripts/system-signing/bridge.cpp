#include "bundle.h"
#include "log.h"
#include <mutex>
#include <ctime>
#include <memory>
#include <fstream>
#include <dirent.h>
#include <sys/stat.h>
#include <openssl/crypto.h>
#include <openssl/err.h>
#include <openssl/provider.h>
#include <openssl/pkcs12.h>
#include <openssl/cms.h>

static std::mutex signing_lock;
static std::once_flag providers;
static void initialize() {
    // Many sideloaders export RC2-encrypted PKCS#12. Legacy is built into the
    // pinned no-module iOS build; never downloads or loads an external module.
    std::call_once(providers, [] { OSSL_PROVIDER_load(NULL,"default"); OSSL_PROVIDER_load(NULL,"legacy"); });
    ZLog::SetLogLever(ZLog::E_NONE);
}
static bool readSmall(const char* path, std::string& out) {
    std::ifstream stream(path,std::ios::binary|std::ios::ate);
    if (!stream || stream.tellg()<=0 || stream.tellg()>8*1024*1024) return false;
    out.resize((size_t)stream.tellg()); stream.seekg(0);
    return (bool)stream.read(&out[0],out.size());
}
static bool permits(const std::string& pattern,const std::string& value) {
    size_t star=pattern.find('*');
    if (star==std::string::npos) return pattern==value;
    // Only a trailing wildcard after a dot, including TEAM.*, is legitimate.
    return star==pattern.size()-1 && star>0 && pattern[star-1]=='.' &&
           value.size()>star && value.compare(0,star,pattern,0,star)==0;
}
static int validate(const char* p12,const char* profile,const char* password,
                    const char* expectedBundleID,const char* expectedTeamID) {
    if (!p12 || !profile || !password || !expectedBundleID || !*expectedBundleID || !expectedTeamID) return 1;
    initialize();
    std::string bytes;
    if (!readSmall(p12,bytes)) return 10;
    const unsigned char* at=(const unsigned char*)bytes.data();
    PKCS12* bag=d2i_PKCS12(NULL,&at,bytes.size());
    EVP_PKEY* key=NULL; X509* certificate=NULL; STACK_OF(X509)* chain=NULL;
    bool decoded=bag && PKCS12_parse(bag,password,&key,&certificate,&chain)==1;
    PKCS12_free(bag); OPENSSL_cleanse(&bytes[0],bytes.size());
    std::unique_ptr<EVP_PKEY,decltype(&EVP_PKEY_free)> releaseKey(key,EVP_PKEY_free);
    std::unique_ptr<X509,decltype(&X509_free)> releaseCert(certificate,X509_free);
    sk_X509_pop_free(chain,X509_free);
    if (!decoded) { ERR_clear_error(); return 11; }
    if (!key || !certificate || X509_check_private_key(certificate,key)!=1) return 12;
    if (X509_cmp_current_time(X509_get0_notBefore(certificate))>=0 ||
        X509_cmp_current_time(X509_get0_notAfter(certificate))<=0) return 13;
    if (!readSmall(profile,bytes)) return 14;
    BIO* input=BIO_new_mem_buf(bytes.data(),(int)bytes.size());
    CMS_ContentInfo* cms=input?d2i_CMS_bio(input,NULL):NULL; BIO_free(input);
    BIO* output=BIO_new(BIO_s_mem());
    // Verify CMS integrity, not Apple trust/revocation or device eligibility.
    bool verified=cms && output && CMS_verify(cms,NULL,NULL,NULL,output,CMS_BINARY|CMS_NO_SIGNER_CERT_VERIFY)==1;
    char* data=NULL; long length=verified?BIO_get_mem_data(output,&data):0;
    jvalue provision;
    bool parsed=length>0 && provision.read_plist(std::string(data,(size_t)length));
    CMS_ContentInfo_free(cms); BIO_free(output);
    if (!parsed) { ERR_clear_error(); return 14; }
    if (provision["ExpirationDate"].as_date()<=std::time(nullptr)) return 15;
    unsigned char* der=NULL; int size=i2d_X509(certificate,&der); bool member=false;
    for (size_t n=0;size>0 && n<provision["DeveloperCertificates"].size();n++) {
        std::string allowed;
        if (provision["DeveloperCertificates"][n].as_data(allowed) && allowed.size()==(size_t)size &&
            CRYPTO_memcmp(allowed.data(),der,size)==0) member=true;
    }
    OPENSSL_free(der);
    if (!member) return 16;
    std::string team=provision["TeamIdentifier"][0].as_string();
    if (team.empty() || (*expectedTeamID && team!=expectedTeamID)) return 17;
    const std::string pattern=provision["Entitlements"]["application-identifier"].as_string();
    // Legacy AppIdentifierPrefix can differ from TeamIdentifier.
    const size_t dot=pattern.find('.');
    if (dot==std::string::npos || !permits(pattern,pattern.substr(0,dot+1)+expectedBundleID)) return 18;
    return 0;
}
extern "C" int CrumbValidateSigningIdentity(const char* p12,const char* profile,const char* password,
                                            const char* bundleID,const char* teamID) {
    std::lock_guard<std::mutex> lock(signing_lock);
    try { return validate(p12,profile,password,bundleID,teamID); } catch (...) { ERR_clear_error(); return 14; }
}
// A single profile may sign multiple bundles only when it actually covers
// each app/extension. Reject before changing any executable, not at installation.
static bool nestedPermitted(const std::string& folder,const std::string& pattern) {
    DIR* directory=opendir(folder.c_str()); if (!directory) return false;
    bool ok=true; struct dirent* entry;
    while (ok && (entry=readdir(directory))) {
        std::string name(entry->d_name);
        if (name=="." || name==".." || name=="CrumbSystemPayload") continue;
        const std::string path=folder+"/"+name; struct stat st;
        if (lstat(path.c_str(),&st)!=0 || S_ISLNK(st.st_mode)) {ok=false;break;}
        if (!S_ISDIR(st.st_mode)) continue;
        if ((name.size()>4 && name.compare(name.size()-4,4,".app")==0) ||
            (name.size()>6 && name.compare(name.size()-6,6,".appex")==0)) {
            jvalue info;
            if (!info.read_plist_from_file("%s/Info.plist",path.c_str())) {ok=false;break;}
            const std::string bundle=info["CFBundleIdentifier"].as_string();
            const size_t dot=pattern.find('.');
            if (bundle.empty() || dot==std::string::npos || !permits(pattern,pattern.substr(0,dot+1)+bundle)) {ok=false;break;}
        }
        if (!nestedPermitted(path,pattern)) ok=false;
    }
    closedir(directory);return ok;
}
// Local CMS signature creation is separate from Apple's installation decision.
extern "C" int CrumbSignPreparedBundle(const char* folder,const char* p12,const char* profile,const char* password,
                                      const char* expectedBundleID,const char* expectedTeamID) {
    if (!folder) return 1;
    std::lock_guard<std::mutex> lock(signing_lock);
    int status=validate(p12,profile,password,expectedBundleID,expectedTeamID);
    if (status) return status;
    std::string secret(password); int result=4;
    try {
        ZSignAsset identity; jvalue info;
        if (!info.read_plist_from_file("%s/Info.plist",folder) || info["CFBundleIdentifier"].as_string()!=expectedBundleID) result=18;
        else if (!identity.Init("",p12,profile,"",secret,false,false,false)) result=2;
        else if (!nestedPermitted(folder,identity.m_strApplicationId)) result=19;
        else { ZBundle bundle; const std::vector<std::string> empty;
            result=bundle.SignFolder(&identity,folder,"","","",empty,empty,true,false,false,false)?0:4; }
    } catch (...) { result=4; }
    if (!secret.empty()) OPENSSL_cleanse(&secret[0],secret.size());
    ERR_clear_error(); return result;
}
