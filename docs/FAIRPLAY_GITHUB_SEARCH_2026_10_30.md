# GitHub search for imported-IPA decryption

Search date: 2026-10-30. The question was whether a normal iPhone app can import
an encrypted IPA and use FairPlay without per-target reinstallation, jailbreak,
TrollStore, JIT, external enrollment secrets or an external decryption server.
This search did not establish a complete implementation satisfying those terms.
It is not proof that none can exist.

## Relevant source findings

| Source inspected | Finding relevant to Crumb |
| --- | --- |
| [FairPlayed, bad7df1](https://github.com/morzelowski/FairPlayed/tree/bad7df138839b5e0ac8b9f20cc4689140b8436fb) | Still embeds the target in the installed helper. No newer upstream commit or PR was returned. |
| Four listed FairPlayed forks | pegakmop, DaniilDavidovich, c12n and eriksmets all point to the same `bad7df1` head. |
| [fairplayed-build, eb017f6](https://github.com/sscascacsa/fairplayed-build/tree/eb017f641998f6755981788c07fc83b51f636dd4) | A differently laid-out build repository. Its main.m, fpcore.c/h and zipio.c/h are identical Git blobs to upstream; it does not contain a different native access route. |
| [lbr77/unfair, a754437](https://github.com/lbr77/unfair/tree/a754437649b1428e07918c0315a678ed7fc57dec) | Implements IPA staging on iOS, but AppBundleStager writes under `/var/containers/Bundle/Application` and changes ownership. UnfairSupport.c explicitly requires root and `libjailbreak` primitives. Importing this staging code would violate the required runtime conditions. |
| [DerekSelander/yacd, 261219d](https://github.com/DerekSelander/yacd/tree/261219da08eafde2d61076d54edbab65e54499b2) | The no-jailbreak description depends on Psychic Paper and targets iOS 13.4.1 or lower. It is not evidence of access on the tested iOS 27 phone. |
| [londek/ipadecrypt, 849a129](https://github.com/londek/ipadecrypt/tree/849a1290afda0d74e608f80859119391dbab9e86) | The reusable account credential is minted using root SSH on a jailbroken device. Subsequent portability does not solve initial credential production inside Crumb. |
| [EthanArbuckle/fairplay-decrypt, bb09a61](https://github.com/EthanArbuckle/fairplay-decrypt/tree/bb09a61046f5c99eb1795589a307646df04b26d0) | A useful mobile-crypto reference, but its credential retrieval remains a privileged-device dependency; not a demonstrated sandbox-local source for the required credential. |
| [UnfairPlayKit, f7ebf85](https://github.com/maclub-net/UnfairPlayKit/tree/f7ebf8506c64a11e35d9b720b916dfe3873244dc) | Swift/C wrappers and an enrollment pipeline; the documented capture stage needs root on Mac and includes its system SC Info and hardware identity. A Swift interface is not an iPhone enrollment producer. |
| [ipatool-rs, 6741709](https://github.com/everettjf/ipatool-rs/tree/6741709ce3e268479dec50910982e827aa3f21e2) | Its README explicitly defers iOS runtime releases and says private enrollment/decryption flows still need end-to-end validation. It uses Unicorn for selected guest execution paths. |
| [Sequoia, 01a542f](https://github.com/sequoiahq/sequoia/blob/01a542ff4953a105153b454a4e8a689ef1bd331c/src/engines/fairplay/program.rs) | The inspected reader maps with `PROT_READ | PROT_EXEC` and retains the input cryptid. It supplies no verified iPhone data-container license handoff. |
| [MobHunt, 65da3c8](https://github.com/ivRodriguezCA/MobHunt/blob/65da3c82b917fb101cf5ed1eb78e9b2d0b2754e7/docs/ios-decryption.md) | Its direct path is a Mac-hosted acquisition/decryption workflow. Its alternative device strategies use Frida or TrollDecrypt. Platform-specific Mac diagnoses do not identify this phone's remap failure. |
| [Cyanide app-decrypt description](https://github.com/kolbicz/Cyanide/blob/13ad360b3e7d2ec4bfe0c4a703cc85a6cea46e71/Cyanide/tweaks/app_decrypt.h) | Uses kernel read/write to recover pages from a target process. This does not demonstrate an ordinary sandbox mechanism; it is not incorporated. |

Repository/code searches additionally covered mremap/cryptid, fairplayd data
containers, sandbox-extension consumption, and account-credential/keybag/SINF
relationships. Search indexing is incomplete and some results are unrelated;
empty search results are not treated as impossibility evidence.

## Substantive new UnFairPlay implementation

The previous local reference was `41f064bfe0ef947e748b785b8ab62697a64bdbe6`.
The newly inspected upstream head is
[`cc2272b3cdd81321080da553ac68ec6bb849a145`](https://github.com/DrWhax/unfairplay/tree/cc2272b3cdd81321080da553ac68ec6bb849a145),
377 commits ahead in GitHub's comparison. The number of commits is not a
decryption result.

Production recovery now uses typed Go rather than captured executable
execution. Relevant operations include deterministic current-container source
derivation, a footer check, decoded SIDB tree parsing, exact selector routing,
and complete validation of the selected class material against the SINF.
This could remove an interpreter dependency if it becomes compatible with
Crumb's producer, but compatibility must be established with the real license.

Seven selected UnFairPlay forks were compared against this head: Tourountzis
and linfeng810 were identical; ilyasaftr, hyhou, 2dubu, KishanBagaria and
fancha0/unfairplay-ios- were behind without ahead commits. This is a bounded
fork review, not a claim to have reviewed every fork.

### Actual local private control

A local, uncommitted harness called upstream `recovery.Recover` without changing
its algorithms or acceptance gates. It reused the retained control's original
SIDB and six-byte **virtual** identity. The SINF was checked byte-for-byte
against its original IPA member; the original associated SUPF was extracted
from that same IPA. No account credential or foreign enrollment was provided.

Observed progress was:

```text
load-inputs
decode-sidb
search-sinf: registry=1, record=1, window=1, completed=0/1
search-sinf: registry=1, record=1, window=1, completed=1/1
error: authenticated SINF accepted no class precursor
```

The signed-keyfile validation, SIDB footer check, tree parse and exact route
checks precede this event in the unmodified upstream control flow. Reaching it
therefore establishes that these upstream checks passed for this control;
they were not independently reimplemented or audited in this experiment.

No enrollment model was recovered, the SINF was not validated, and the IPA was
not decrypted. The run was finite and did not end in an emulator budget error.
One rejected, deterministically selected class tuple is not evidence that all
other algorithms or producer versions are impossible. It is also not grounds
for trying other account records or random windows.

The private harness would continue with explicit SINF resolution and complete
IPA processing only after recovery succeeded; neither later stage ran. No
private state, identifiers, credential, model or IPA is in this document or
in the committed result summary.

The anonymous result is recorded in
[`research/unfairplay-cc2272-private-control-summary.json`](research/unfairplay-cc2272-private-control-summary.json).
Upstream `go test ./...`, `go vet ./...` and `go test -race ./...` passed.
The first public-suite run hit two permission-fixture failures because the
research process uses umask 0077, turning requested 0644 test files into 0600.
The public tests were rerun with process-local umask 0022; no permission gate or
test assertion was changed. Private-input permissions remained 0600/0700.

## Consequence for integration

The updated typed engine is a concrete new research reference for the
StoreAgent-produced state, not a demonstrated correction of the system-reader
access problem. The next causal boundary for that engine is the selected
class-material transformation and its agreement with Crumb's StoreAgent/CoreFP
producer. This is distinct from changing storage locations for mremap.

No runtime engine is replaced, no SINF check is weakened, and no new device
build is justified by this search alone. Existing public and private controls
must establish the actual correction before it is offered as a working route.
