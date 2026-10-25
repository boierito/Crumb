# Pinned FairPlay public sources

These are the previously pinned public DrWhax/ipatool source distributions.
Both original GitHub archive URLs began returning 404 on 2026-10-25. Keeping
these exact distributions here makes the native build and portable tests
reproducible without access to that repository. Each includes its upstream
license and recovery attribution. They contain no Crumb account captures or
runtime recovered enrollment models.

| Revision | Archive SHA-256 |
| --- | --- |
| `ef98227e061838320424744d21a513baba7ade28` | `49f169ef2e14c89e3b55a9528057e7db66f236ca5425ff00e47d32f64dddf44d` |
| `f7242cce58077a4af8db738ed295ee23ad8f067d` | `e89a89c4af74a8b0e6f682efdf89edc36bc093874d706d8f27f43e657bd3ec39` |

The first archive was reproduced from the already available Git object using
`git archive --format=tar --prefix=ipatool-REV/ REV | gzip -n`; its bytes match
the historically verified archive hash. Its source tree is
`d9f1768d69fb5532d2ca16e7712cf770f7210479`. The second archive is the original
verified download retained from the successful 23031/23032 builds.

The scripts verify the archive hashes before extraction. SAP/StoreAgent remains
independently pinned. Changing these archives requires an explicit dependency
update and updated integrity checks; a moving upstream branch is not a fallback.
