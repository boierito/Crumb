# WaffleStore 2.3.0-dev.1 — SAP experimental

Este fork conserva el proyecto Xcode, la UI, favoritos, historial, búsqueda y
código original de WaffleStore. **Todavía no es una recuperación funcional del
login o las descargas.** La entrega implementa la prueba SAP prioritaria, un
intérprete sin JIT y una IPA de diagnóstico; falta validarla en un dispositivo
jailed antes de avanzar al login, siguiendo el orden solicitado.

## Estado verificable

| Hito | Resultado |
|---|---|
| 1. Build original | Debug y Release con Xcode 26, IPA original generada |
| 2. SAP dentro de iOS jailed | Núcleo integrado como biblioteca estática arm64; ejecución en dispositivo pendiente |
| 3. ActionSignature | Intercambio SAP real con Apple y firma no vacía en Linux/TCI; la IPA contiene el mismo guest y un diagnóstico para probarlo en iOS |
| 4–6. Login, 2FA, DSID/token/storefront | Pendientes; backend original sigue presente, no conectado al nuevo signer |
| 7–8. Búsqueda/versiones | Código original conservado; recuperación autenticada pendiente |
| 9–10. Descargar última/antigua | Pendiente del signer jailed, login y kbsync |
| 11. Exportar IPA descargada | Share Sheet original conservado; flujo nuevo de download/export no validado |
| 12. IPA normal en iOS 27 | IPA resignable preparada; ejecución real pendiente |

El log `docs/evidence/tci-sap-smoke-linux.log` registra una firma de 501 bytes.
Es una firma de un cuerpo de prueba sin Apple ID; no prueba aceptación de una
petición de login. Los tests con FakeGuest prueban protocolo y estados, nunca
la validez criptográfica frente a Apple.

## 1. Cómo funcionaba WaffleStore y qué se rompió

La UI llama a IPATool/StoreClient en `WaffleStore/Functions/IPATool.swift`.
StoreClient obtiene únicamente authenticateAccount del Bag, construye un GUID
a partir del Apple ID y envía JSON sin firma. Almacena DSID/token/storefront y
cookies para llamar al endpoint volumeStoreDownloadProduct hardcodeado.
Los IDs de versiones vienen de metadata; la versión visible depende también
de un servicio externo. La IPA se descarga, descomprime, recibe metadata/SINF,
se vuelve a empaquetar y se sirve por localhost para instalar mediante Safari.

La implementación moderna de ipatool exige inicialización SAP y firma de los
bytes exactos del cuerpo plist. El original omite ambos. Además, authenticate
devuelve false antes de terminar su Task, hay casts que pueden causar crashes,
no existe purchase moderno y download omite el nuevo kbsync. El mapa detallado
está en [docs/BACKEND_AUDIT.md](docs/BACKEND_AUDIT.md).

El original compila sin modificaciones funcionales. Falló inicialmente con
Xcode 16.4 porque PartyUI exige Swift 6.2; Xcode 26 lo resuelve. El tag
`upstream-wafflestore-2.2.2` conserva exactamente el HEAD original d508e53.

## 2. Referencia y diferencia respecto de ipatool

Referencia fija: `majd/ipatool` commit
`3411d57f451f5111ae115641c22f7ed17bbd5fbe`, 1 de octubre de 2026.
Incluye los cambios posteriores al primer SAP de agosto: 2FA, redirects
firmados, timeouts, identidad estable, versiones y kbsync. Ver
[docs/IPATOOL_AUDIT.md](docs/IPATOOL_AUDIT.md).

El ipatool iOS original es una CLI para jailbreak, con no-container y Unicorn
JIT. Aquí **no se ejecuta ni enlaza su CLI**: se compilan únicamente el loader
Mach-O, assets, shims y machine SAP como una biblioteca Go con ABI C. El
transporte Bag/certificado/setup y control de sesión están adaptados a Swift.
Los assets se descargan de Apple con sus hashes/tamaños originales y se guardan
en un Caches explícito dentro del sandbox. No se incluyen assets Apple en la
IPA ni se accede a frameworks privados del sistema iOS.

## 3. Cómo se reemplazó el JIT

La incompatibilidad concreta está en el traductor de Unicorn: escribe código
host y necesita memoria ejecutable. El parche iOS de ipatool cambia RW a RX
mediante mprotect, abortando si falla. En una prueba Linux con PROT_EXEC
prohibido, Unicorn termina al reservar su buffer. La prueba de dispositivo
solicita esos permisos sin ejecutar memoria unsigned ni abortar.

`scripts/prepare-unicorn-tci.py` restaura los tres archivos TCI de QEMU 5.0,
verificados por SHA-256, sobre Unicorn 2.1.4 también fijado/verificado. Adapta
los campos del TCGContext a Unicorn, añade el intérprete al build, utiliza un
PC de helper en TLS, reserva buffers RW y desactiva MAP_JIT/protecciones APRR
del runtime para este modo. La CPU host ejecuta el intérprete compilado y
firmado; los bloques generados son **bytecode almacenado como datos**.

Pruebas realizadas:

1. Unicorn normal ejecuta MOV EAX,42 y falla al denegar PROT_EXEC.
2. TCI ejecuta call/return, stack y hooks con PROT_EXEC denegado: RAX=42,
   cuatro hooks, status=0.
3. Los tests de máquina de ipatool verifican argumentos de stack, import no
   soportado, allocator y limpieza de memoria con el runtime TCI.
4. El guest real inicializa SAP, obtiene certificado/configuración del Bag,
   realiza el intercambio y firma un cuerpo de prueba en Linux.

Esto convierte la alternativa en código compilable y comprobado en host;
**no demuestra todavía ejecución arm64/iOS ni elimina posibles bugs del
intérprete**. Hay que probar instrucciones/SIMD, hooks anidados, timeout,
faults, consumo de RAM, reinstalaciones y kbsync en la arquitectura real.
TCI se restauró de una versión antigua compatible con la base QEMU de Unicorn;
su rendimiento y cobertura completa no se presentan como producción.

## 4. Arquitectura implementada

```text
Settings → SAPDiagnosticView
  → MemoryCapability + TCIProbe
  → KeychainMachineIdentity
  → SAPProtocol.bag (endpoints/version dinámicos)
  → SAPSession actor
      → NativeSAPGuest → ABI C → guest de ipatool → Unicorn/TCI
      → SAP certificate → Exchange(state 1)
      → SAP setup POST → Exchange(state 0)
      → sign(bytes exactos) → Base64 → X-Apple-ActionSignature
```

La identidad consta de seis bytes aleatorios, unicast y localmente
administrados, persistidos en un item generic-password de Keychain
AfterFirstUnlockThisDeviceOnly. GUID y hardwareID se derivan de los mismos
bytes. No depende del MAC físico ni del Apple ID. El logout original no
borra esta nueva identidad. Una nueva firma/team/access group puede cambiar
el acceso al item: no se promete identidad idéntica entre equipos de firma.
La aceptación completa de esta identidad por login/download queda pendiente.

Los secretos SAP viven en memoria y se reconstruyen al iniciar. El bridge
libera/limpia los buffers y hace teardown; no se guardan secretos SAP en
UserDefaults. No se ha migrado aún el almacenamiento original de cuentas:
su archivo authinfo y fallback .authkey están auditados, pero siguen siendo
legacy. No debe confundirse KeychainIdentity nuevo con una migración terminada
de passwordToken/cookies.

El Bag bootstrap usa init.itunes.apple.com/bag.xml. Setup/certificado y endpoint
de auth salen del Bag; no hay fallback unsigned. Se soporta el envelope
Document/Protocol y certificados del CDN mzstatic. Para la prueba SAP se acepta
que un Bag aún anuncie el auth endpoint legacy; eso no activa login legacy en
el módulo nuevo. Se rechazan endpoints HTTP, hosts ajenos y SAP distinto de 200.

## 5. Login, download e instalación pendientes

No se conectó la UI de login al signer porque la ejecución jailed es el gate
prioritario. Después de validar el diagnóstico en dispositivo:

1. Reemplazar Bool/blocking authenticate por async y estados typed; firmar el
   plist exacto; redirects validados/bounded que preserven POST/body/attempt.
2. 2FA independiente del password persistido; normalización seis dígitos y
   errores Apple útiles; distinguir errores credential de respuestas HTML.
3. Credenciales/cookies en Keychain; logout que cierre sesión/cookies sin
   regenerar hardware; pedir password al expirar en vez de guardar el código.
4. Mantener búsqueda y datos locales; verificar versiones/externalVersionId
   mediante metadata de respuesta y Info.plist real de la IPA.
5. Purchase únicamente gratuito y manejo de interacción/licencia.
6. Portar StoreAgent/kbsync del mismo commit sobre TCI, cache por DSID/GUID y
   regeneración una vez si se rechaza. Bag ent/download y fallbacks actuales;
   no reenviar X-Token a CDN/redirects.
7. URLSession download con progreso real, validación HTTP/ZIP, timeouts/retries
   limitados, sandbox Documents/Downloads, Files y Share Sheet independientes
   de la instalación.

Descargar un paquete App Store con SINF no lo desencripta. SideStore/AltStore
pueden instalar esta **app WaffleStore** al resignarla, pero no se promete que
puedan instalar cualquier IPA cifrada obtenida del App Store. El mecanismo
original itms-services/Safari permanece sin certificar en iOS 26/27. No se
introdujo jailbreak, AppSync, TrollStore ni entitlements privados para ello.

## 6. Compilar y generar IPA

macOS con Xcode 26+ (Swift 6.2 por PartyUI), Command Line Tools, CMake, Python 3
y Go 1.25+. Abrir `WaffleStore.xcodeproj`, seleccionar WaffleStore y un destino
**iOS físico arm64**. El build phase prepara las bibliotecas estáticas y caché
de build basada en huella de sources/SDK; requiere red la primera vez. Para
un build firmado, seleccionar el team habitual. No hace falta no-container,
allow-jit ni otros entitlements privados.

```sh
bash scripts/ensure-native.sh
bash ipabuild.sh             # Release
bash ipabuild.sh --debug     # Debug
swift test --package-path MapleSyrup/SAPKit
```

El script original genera IPA unsigned; el CI conserva logs y empaqueta sin
necesitar certificados. El módulo Swift admite tests host en macOS; el target
completo con Go/TCI actualmente no está configurado para iOS Simulator.
El fingerprint detecta cambios del bridge, scripts y SDK; se puede eliminar
`build/Native` para reconstruir las bibliotecas.

GitHub Actions: push/PR → Debug y Release en macos-15/Xcode 26 + tests del
protocolo en macOS y restricciones/interpreter en Linux → artifacts.
workflow_dispatch permite compilar el tag original con el mismo workflow.
Un tag `v2.3.0-dev.1` genera una prerelease **draft** con ambas IPAs. No publicar
una release hasta resolver las licencias y validar los hitos; ver notices.
CFBundleShortVersionString debe ser numérico: 2.3.0, CFBundleVersion 23001;
el sufijo dev.1 vive en el tag/changelog, no en el plist.

## 7. Instalar y probar iOS 27

Resignar la IPA Release con SideStore, AltStore, Sideloadly o certificado de
desarrollador habitual. Instalar y abrir **sin debugger, JIT externo ni
TrollStore**. Entrar en Settings → SAP diagnostic; ejecutar primero sin la
opción de red y exportar el resultado. Después activar Initialize SAP and sign
a test body y repetir; se descargarán assets Apple al Caches del sandbox.
Los resultados esperados y la matriz están en [TESTING.md](TESTING.md).

Para avanzar hacen falta evidencias de ejecución real, no sólo screenshots de
build: versión iOS/dispositivo/método de firma; errno; TCI status/RAX/hooks;
identidad estable después de matar/reabrir; SAP setup completo y firma no vacía.
No enviar password, token, cookies, códigos 2FA ni contenido de firma. El
diagnóstico exporta sólo estados/tamaños. Login y download se prueban **cuando
se implemente su fase**, no con el formulario legacy de esta entrega.

## 8. Límites y licencias

La memoria guest/scratch/heap/stack y las imágenes Apple tienen un coste alto;
los devices con poca RAM pueden recibir jetsam. La descarga inicial por rangos
del update macOS necesita acceso a la CDN Apple; cambios del paquete/hash deben
fallar explícitamente. Se conservan los bounds del guest y timeouts, pero no
hay botón de cancelación durante una llamada nativa; esto sigue experimental.
Los permisos RW→RX por sí solos no prueban SAP ni autorización del backend.
No hubo iPhone/iPad conectado a este entorno para verificar los hitos jailed.

MIT de ipatool conservada. Unicorn/QEMU TCI incluyen GPL; cada artifact entrega
sus fuentes modificadas correspondientes, y el source de WaffleStore está en
el fork. Upstream WaffleStore no contiene licencia de redistribución: es un
problema pendiente que requiere aclaración de sus titulares, no un permiso que
pueda inferirse del nombre open source. Los assets Apple son propietarios y
no se redistribuyen. Ver [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
