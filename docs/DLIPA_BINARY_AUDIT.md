# DLiPA v1.4 purchase and region audit

Inspected October 5, 2026, using its public release IPA and read-only ARM64
Objective-C metadata/disassembly. No application code or binary is copied or
redistributed. The repository contains README/screenshots, not app source.

- Release: https://github.com/AhmedBafkir/DLiPA/releases/tag/v1.4
- IPA SHA256: `72012beabf55a12d70c5ad3bbe95384e7fd450c8de1bf6fbc62f8844512eea0f`
- README commit: `8bc61ee9d4642f2706b99e9012425c58077015b7`
- Comparison ipatool: `3411d57f451f5111ae115641c22f7ed17bbd5fbe`

## Follow the caller, not just strings

DLIPA Store manager `downloadDataForModel:specificVer:retriesCount:handler:`
starts at `0x10005a06c`: reads model.appId then calls
`downloadWithAppId:specificVer:completion:` at `0x10005a134`.
Its completion switches on returned error category. License-required category 3
(`0x10005a250–0x10005a340`) presents an acquisition action, which calls
`handlePurchaseForModel:retriesCount:handler:` at `0x10005a570`.
That obtains model.appId and calls `purchaseWithAppId:completion:` at
`0x10005a630`. Purchase-success category 1 resumes download at `0x10005a6cc`.
Categories 4–6 trigger login; other failures are shown, not assumed successful.

The purchase worker starts at `0x10004a548`. It obtains DSID/pod from
APStoreAccount.sharedInstance and derives p<pod>-buy.itunes.apple.com, using
MZFinance buyProduct without a GUID query. It sends account DSID/token and the
account storefront (`0x10004a854–0x10004a888`). The selected search country is
stored separately (`Storefronts_selected`, `0x100023930–0x100023964`). No account
storefront replacement was found on this selection path.

Its request body builds **nine** keys (`0x10004a8e4–0x10004a9a4`):
appExtVrsId, hasAskedToFulfillPreorder, buyWithoutAuthorization, hasDoneAgeCheck,
guid, price, pricingParameters, productType and salableAdamId. Values are the
same zero/true/STDQ/C pattern; appId is an NSString property and flows unchanged
from model through caller/block into salableAdamId. The User-Agent literal at
`0x10004a7f8` is Configurator/2.20 (Macintosh; OS X 26.5.1; 25F80)
AppleWebKit/1624.2.5.11.4. ipatool uses twelve body keys, including needDiv,
origPage and origPageLocation, an integer app ID, and its shared Configurator/2.17
User-Agent. WaffleStore 23014 matched the ipatool form, not DLiPA's exact request.

On an explicit 2059 (`0x10004abd8`), DLiPA changes pricingParameters to GAME and
re-serializes XML (`0x10004abfc–0x10004ac98`). Its loop can revisit that branch;
WaffleStore intentionally permits only one alternate. A final 2059 is treated
as temporarily unavailable (`0x10004adfc–0x10004ae60`). DLiPA also handles HTTP
302 by replacing the request URL from Location, with a five-redirect bound
(`0x10004aa94–0x10004ab14`). This is not evidence that a redirect caused the
user's 2059. Following arbitrary Location with account secrets would be unsafe;
no unvalidated redirect handling is copied into this change.

## What 23015 changes and what remains unknown

For an app resolved in a country different from the account, reproduce DLiPA's
nine-field plist, string app ID and 2.20 purchase User-Agent; retain account
storefront/token/pod/GUID and the bounded STDQ/GAME flow. Other purchases retain
the existing ipatool profile/2.17 agent. This is a controlled, evidence-backed
compatibility attempt, not a proven fix for live regional acquisition.

The static binary does not show that Venmo currently grants a license to the same
AR account, nor provide Apple's server-side reason for 2059. The user reports
23014 still failed. Compare the **same app, account and ownership** in DLiPA and
WaffleStore; distinguish a previous US license from first-time acquisition.
The limited existing log now records profile, plist key count, pricing and
numeric Apple failure. It never prints credentials, full cookies or signed URLs.
No additional purchase retries or invented account-country change are added.
