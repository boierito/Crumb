"""Adapt the pinned MIT zsign library for an in-process, private iOS signer."""
from pathlib import Path
import sys

root = Path(sys.argv[1])

def replace(file, old, new):
    path = root / file
    text = path.read_text()
    if text.count(old) != 1:
        raise SystemExit('unsupported signing source contract: '+file)
    path.write_text(text.replace(old, new))

# No subprocess-based diagnostic helper is included in the app build.
replace('signing.cpp', '#ifndef _WIN32\n\tif (ZFile::IsFileExists("/usr/bin/csreq"))',
        '#if !defined(_WIN32) && !defined(CRUMB_IOS_SIGNER)\n\tif (ZFile::IsFileExists("/usr/bin/csreq"))')
replace('openssl.cpp', 'ERR_print_errors_fp(stdout);', 'ERR_clear_error(); // Private errors must not print identity data.')
path = root / 'openssl.cpp'
text = path.read_text()
start, end = text.index('bool ZSignAsset::GetCMSContent('), text.index('bool ZSignAsset::GetCertSubjectCN(void*')
text = text[:start]+'''bool ZSignAsset::GetCMSContent(const string& input, string& output)
{
    if (input.empty() || input.size() > INT_MAX) return false;
    BIO* bio = BIO_new_mem_buf(input.data(), (int)input.size());
    if (!bio) return false;
    CMS_ContentInfo* cms = d2i_CMS_bio(bio, NULL);
    BIO_free(bio);
    if (!cms) return CMSError();
    ASN1_OCTET_STRING** content = CMS_get0_content(cms);
    bool ok = content && *content && ASN1_STRING_length(*content) > 0;
    if (ok) output.assign((const char*)ASN1_STRING_get0_data(*content), ASN1_STRING_length(*content));
    CMS_ContentInfo_free(cms);
    return ok;
}

''' + text[end:]
path.write_text(text)
replace('openssl.h', '\tZSignAsset();', '\tZSignAsset();\n\t~ZSignAsset();')
replace('openssl.cpp', 'bool ZSignAsset::Init(', '''ZSignAsset::~ZSignAsset()
{
    EVP_PKEY_free((EVP_PKEY*)m_evpPKey);
    X509_free((X509*)m_x509Cert);
    sk_X509_pop_free((STACK_OF(X509)*)m_caCerts, X509_free);
}

bool ZSignAsset::Init(''')
replace('openssl.cpp', '\tm_evpPKey = evpPKey;\n\tm_x509Cert = x509Cert;\n\treturn true;', '''\tm_evpPKey = evpPKey;
    m_x509Cert = x509Cert;
    if (X509_check_private_key(x509Cert, evpPKey) != 1 ||
        X509_cmp_current_time(X509_get0_notBefore(x509Cert)) > 0 ||
        X509_cmp_current_time(X509_get0_notAfter(x509Cert)) <= 0) return false;
    unsigned char* der = NULL;
    int length = i2d_X509(x509Cert, &der);
    bool authorized = false;
    if (length > 0) {
        for (size_t i = 0; i < jvProv["DeveloperCertificates"].size(); i++) {
            string allowed;
            if (jvProv["DeveloperCertificates"][i].as_data(allowed) &&
                allowed.size() == (size_t)length && !memcmp(allowed.data(), der, length)) authorized = true;
        }
    }
    OPENSSL_free(der);
    return authorized;''')

# The purchased target is a RESOURCE. Sign the helper and its own frameworks,
# and include the target in CodeResources, without rewriting any target file.
path = root / 'bundle.cpp'
text = path.read_text()
start, end = text.index('bool ZBundle::GetObjectsToSign('), text.index('bool ZBundle::GenerateCodeResources(')
block = text[start:end]
old = '\t\t\tif (bFolder) {'
assert block.count(old) == 1
block = block.replace(old, '''            const string payload = m_strAppFolder + "/CrumbSystemPayload";
            if (strPath == payload || strPath.compare(0, payload.size()+1, payload+"/") == 0) return false;
''' + old)
old = 'ZFile::EnumFolder(strFolder.c_str(), true, NULL, '
assert block.count(old) == 1
block = block.replace(old, '''ZFile::EnumFolder(strFolder.c_str(), true, [&](bool, const string& candidate) {
        const string payload = m_strAppFolder + "/CrumbSystemPayload";
        return candidate == payload || candidate.compare(0, payload.size()+1, payload+"/") == 0;
    }, ''')
path.write_text(text[:start]+block+text[end:])
