# Login real — evidencia reportada por el usuario

Fecha de incorporación: 2026-10-05. iPhone, iOS 27.0.1, build 23002.
El usuario había informado ksign + certificado sin JIT; no indicó el modelo
exacto, configuración Release o debugger. No se infieren.

Primero aparecieron HTTP404 y HTTP204 con retries. Después Apple pidió 2FA,
se siguió el Store pod y se guardó sesión. Reporte recibido:

```text
WaffleStore authentication probe v2
iOS=27.0.1
app-build=23002
password-persistence=false
signer=tci-no-jit
stage=Resolving Store Bag
stage=Initializing SAP
stage=Signing authentication request
stage=Contacting Apple
stage=Following Store pod
stage=Signing authentication request
stage=Contacting Apple
stage=Saving session in Keychain
outcome=authenticated
DSID=present
passwordToken=present
storefront=present
pod=present
secret-values=withheld
```

El usuario confirmó: «si funciono el reabrir». Eso aporta evidencia de login
aceptado y restauración local en esa instalación, no garantiza vigencia futura
del token. Los errores de búsqueda -999 eran cancelaciones esperadas de consultas
anteriores al escribir. No prueban un fallo de la sesión.

Downgrade no funcionaba porque 23002 lo mantenía deshabilitado. No se recibieron
pruebas de purchase, versions, download, export, logout o instalación de IPA.
