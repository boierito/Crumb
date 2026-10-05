# WaffleStore

## Experimental revival — 2.3.0-dev.1

This fork preserves the original app and adds a no-JIT SAP guest diagnostic.
**Modern login and downloads are not restored yet.** The original and modified
projects build in CI; the priority is validating Settings → SAP diagnostic on a
conventionally sideloaded Release IPA on iOS 27 without a debugger. See
[IMPLEMENTATION.md](IMPLEMENTATION.md), [TESTING.md](TESTING.md),
[changed files](CHANGED_FILES.md), and [license notices](THIRD_PARTY_NOTICES.md).

The upstream instructions below are preserved as historical documentation and
are not a claim that their legacy login/install flow works in this branch.
A **jailed** app store app downgrader, based off of [MuffinStoreJailed](https://github.com/mineek/MuffinStoreJailed-Public) and [PancakeStore](https://github.com/jailbreakdotparty/PancakeStore). Supports iOS 16.4+ and does not use any exploits.

[PancakeStore (by jailbreakdotparty)](https://github.com/jailbreakdotparty/PancakeStore/releases/latest) • [MuffinStoreJailed (by mineek)](https://github.com/mineek/MuffinStoreJailed-Public)

>[!IMPORTANT]
>WafleStore has been temporarily discontinued. Apple has changed their backend for authentication, meaning that we can no longer use the method we have been relied on. WaffleStore will not work until someone finds another backend solution. For the time being, the repo will still be up. Don't come complain about that it isn't working, there is nothing we can do.

>[!WARNING]
>Use this tool at your own risk! You may lose app data, and other damage could occur.

## Prerequisites
- An iPhone or iPad running iOS 16.4 or later.
- An Apple ID that has 2FA **enabled.**
- Some method of sideloading.

## How do I downgrade an application?
1. Make sure that you have uninstalled the app you want to downgrade beforehand. To keep your app data, offload the app instead of uninstalling it.
2. Install WaffleStore using your preferred sideloading method.
3. Log in your Apple ID, and click "Send 2FA Code."
4. You will likely receive a 2FA popup. Click "Allow" and write down the code somewhere. Then, type the code into the text field. If you do not get this pop up, just type six random numbers instead. Finally, click "Log In."
5. Get the app store link of the application you want to downgrade.
6. Type or paste the link to the app into the field.
7. Click "Downgrade App," select the version you'd like to downgrade to, and wait patiently!

## Troubleshooting
- **App crashes when trying to log in:** You likely input the wrong 2FA code. Kill the app from the app switcher and try again.
- **App doesn't progress when clicking "log in":** You likely input the wrong Apple ID and/or Password. Kill the app from the app switcher and try again.
- **I don't receive a 2FA code when clicking "Authenticate":** Type six random numbers into the field and log in. This likely works because you had logged into MuffinStoreJailed or WaffleStore in the past, even on a different device.
- **The app crashes when I try to downgrade:** You likely did not purchase (download) the app in the past. Download the app beforehand and remove it, or purchase it from another device.
- **"Safari can't open the page" error:** Go to Settings > Apps > Safari, and under "Privacy & Security," disable "Not Secure Connection Warning." Also ensure that you don't have any VPN or DNS settings that might be affecting this.
