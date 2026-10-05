# User report: Store build 23003

Reported iPhone iOS 27.0.1; same installation previously authenticated and restored
its Keychain session. Store probe v3 reached latest iOS externalVersionId, Store
Bag, account-bound kbsync generation and ent/download, then failed with code=10.
The UI requested logout; subsequent login required repeated attempts, and Apple
did not repeat its earlier 2FA challenge. No HTTP/Apple failure code was reported.

code=10 was the local Swift/NSError classification, not Apple's failureType.
Do not infer an expired account or accepted kbsync from it. Dev.4 fixes premature
fallback termination and adds diagnostics. Physical versions/download validation
remains pending. Absence of a new 2FA challenge is not itself a login failure.
