# WaffleStore 2.3.0-dev.2 — login firmado experimental

Este fork conserva el proyecto Xcode, la UI, favoritos, historial y búsqueda de
WaffleStore. Implementa SAP sin JIT y conecta el login/2FA moderno con sesión en
Keychain. **La aceptación del login por Apple todavía requiere prueba real.**
Versions, purchase, kbsync, download y exportación del flujo nuevo siguen
pendientes; el botón de downgrade queda deshabilitado en esta build.

El usuario reportó dos inicializaciones SAP exitosas en un iPhone iOS 27.0.1,
firmado con ksign y certificado, sin JIT: TCI status=0, RAX=42, hooks=4,
SAP v200 setup completo y ActionSignature Base64 de 668 caracteres.
Ver [evidencia de dispositivo](docs/evidence/ios27-ksign-user-report.md).
Es una prueba de ejecución del signer, no de login aceptado por Apple.

## Estado verificable

| Hito | Resultado |
|---|---|
| 1. Build original | Debug y Release con Xcode 26, IPA original generada |
| 2. SAP dentro de iOS jailed | Ejecución TCI/SAP reportada por usuario en iPhone iOS 27.0.1, ksign sin JIT |
| 3. ActionSignature | Setup y generación de firma reportados en iOS 27.0.1; aceptación de login pendiente |
| 4–6. Login, 2FA, DSID/token/storefront | Implementados y conectados; pruebas fixture/Keychain y aceptación física de Apple se distinguen en TESTING.md |
| 7–8. Búsqueda/versiones | Código original conservado; recuperación autenticada pendiente |
| 9–10. Descargar última/antigua | Pendiente de validar login y portar kbsync/download |
| 11. Exportar IPA descargada | Share Sheet original conservado; flujo nuevo de download/export no validado |
| 12. IPA normal en iOS 27 | Signer reportado funcional con firma convencional; flujo Store completo pendiente |

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
bytes. No depende del MAC físico ni del Apple ID. El logout nuevo no
borra esta identidad. Una nueva firma/team/access group puede cambiar
el acceso al item: no se promete identidad idéntica entre equipos de firma.
La aceptación completa de esta identidad por login/download queda pendiente.

Los secretos SAP viven en memoria y se reconstruyen en cada login/validación 2FA.
El bridge libera/limpia buffers y hace teardown. La cuenta se almacena como item
generic-password de Keychain AfterFirstUnlockThisDeviceOnly (sin access group
especial), con DSID, passwordToken, storefront, pod, GUID, email, nombre y cookies.
Nunca se guarda password, código 2FA o buffers SAP. Al actualizar se utiliza
SecItemUpdate; no se borra primero una sesión válida. Al restaurar se valida
estructura, endpoint y GUID, pero **no se afirma que el token siga válido en Apple**.

La migración descarta authinfo, .authkey y la clave EC legacy sin descifrarlos;
se exige un nuevo login firmado. No existe fallback de secretos a archivos ni
UserDefaults. Logout borra la cuenta/cookies, cierra URLSession y limpia el estado
UI sin matar la app ni regenerar la identidad. La persistencia entre reinstalación,
cambio de certificado/team o access group no está garantizada.

El Bag bootstrap usa init.itunes.apple.com/bag.xml. Setup/certificado y endpoint
de auth salen del Bag; no hay fallback unsigned. Se soporta el envelope
Document/Protocol y certificados del CDN mzstatic. Para la prueba SAP se acepta
que un Bag aún anuncie el auth endpoint legacy; eso no activa login legacy en
el módulo nuevo. Se rechazan endpoints HTTP, hosts ajenos y SAP distinto de 200.

## 5. Login y 2FA implementados; download pendiente

```text
Formulario original → AppData.startAppleLogin (async)
  → identidad Keychain → Store Bag actual → endpoint validado
  → SAPSession / NativeSAPGuest / TCI → setup
  → AppleAuthentication → plist XML (appleId, attempt, guid, password, rmp, why)
  → firma de los mismos bytes que URLSession envía
  → POST + X-Apple-ActionSignature
  → respuesta Apple / redirects / 2FA
  → DSID + passwordToken + storefront (+ pod opcional) → Keychain
  → IPATool/StoreClient adapter (Store nuevo todavía no habilitado)
```

Se sigue el protocolo del commit ipatool 3411d57: Content-Type conserva
application/x-www-form-urlencoded aunque el cuerpo es plist XML; se envían `rmp` y `attempt` string como en ipatool. No se manda JSON.
Sólo se usan endpoints de auth del Bag que pasan la misma validación de
host buy.itunes.apple.com o *-buy.itunes.apple.com, path authenticate y HTTPS.
El Bag-only probe admite otros endpoints para diagnosticar SAP; el login
rechaza rutas no soportadas antes de enviar credenciales. No hay fallback.

301/302/307/308 se siguen manualmente, hasta cuatro hops, preservando POST,
body y attempt, y firmando esos bytes en cada envío. 303 y destinos ajenos,
userinfo, fragmentos, puertos diferentes a 443 o rutas codificadas se rechazan.
URLSession tiene un cookie jar efímero por intento de login, timeout de 30s y
redirect automático deshabilitado. Cookies relevantes se guardan con la cuenta
en Keychain; no se utiliza el jar global ni se reenvían secretos a hosts ajenos.

`-5000` en la primera respuesta genera únicamente el segundo intento lógico,
igual que ipatool. BadLogin.Configurator_message activa el campo de 2FA, sin
persistir credenciales. La segunda petición usa password + código ASCII de
seis dígitos sin espacios, sólo en memoria. Códigos inválidos se rechazan antes
de SAP/red. Si Apple pide un código fresco, se conserva la pantalla 2FA y se
limpia el código anterior. Se puede cancelar/reiniciar para otra cuenta.
Los mensajes Apple se limitan y se redactan antes de mostrarlos; el diagnóstico
exportado sólo contiene etapas/categorías/estado y nunca cuerpos o valores.

204/403/404/429/5xx sin respuesta Apple útil y timeouts/transporte transitorio:
hasta tres envíos por petición, backoff 10/20s, Retry-After respetado hasta 30s;
si Apple exige más, se detiene y pide esperar. Un plist con fallo de credenciales
se interpreta directamente y no se trata como HTML transitorio. Cancelación
interrumpe URLSession/backoff, pero una llamada C/Go/TCI en curso debe acabar
antes del teardown; no se promete cancelar instantáneamente el guest.

Queda por validar con Apple real: login, 2FA correcta/incorrecta/expirada,
cuenta bloqueada, token recibido, restauración tras relaunch y logout en iOS.
No hay cuentas/passwords disponibles en el entorno, ni se solicitan aquí.

Siguiente fase Store, después de esa prueba:

1. Verificar versiones/externalVersionId contra metadata y la IPA real.
2. Purchase exclusivamente gratuito/licencia existente, sin automatizar pagos.
3. Portar StoreAgent/kbsync del mismo commit sobre TCI y cache por DSID/GUID.
4. Resolver endpoints/download desde Bag y respetar redirects, retries y CDN.
5. URLSession download con progreso, verificación HTTP/ZIP y exportación
   sandbox/Files/Share Sheet separada de instalación.

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
Un tag `v2.3.0-dev.2` genera una prerelease **draft** con ambas IPAs. No publicar
una release hasta resolver las licencias y validar los hitos; ver notices.
CFBundleShortVersionString debe ser numérico: 2.3.0, CFBundleVersion 23002;
el sufijo dev.1 vive en el tag/changelog, no en el plist.

## 7. Instalar y probar iOS 27

Resignar la IPA Release con SideStore, AltStore, Sideloadly o certificado de
desarrollador habitual. Instalar y abrir **sin debugger, JIT externo ni
TrollStore**. Entrar en Settings → SAP diagnostic; ejecutar primero sin la
opción de red y exportar el resultado. Después activar Initialize SAP and sign
a test body y repetir; se descargarán assets Apple al Caches del sandbox.
Los resultados esperados y la matriz están en [TESTING.md](TESTING.md).

Con dev.2 usar el formulario original para login: Apple ID/password → botón
Send 2FA Code → esperar respuesta → si aparece el campo, introducir un código
actual → Log In. No concatenar manualmente el código al password. Tras éxito,
esperar que el log indique DSID/token/storefront recibidos y guardados (sin valores).
En Settings → Export authentication diagnostic se comparte un reporte sanitizado.
Reabrir confirma sólo la carga de la sesión guardada, no validez online del token.
Probar logout desde el menú y otro login. Seguir los casos de [TESTING.md](TESTING.md).
No compartir password, tokens, cookies, códigos, firma ni capturas con secretos.

## 8. Límites y licencias

La memoria guest/scratch/heap/stack y las imágenes Apple tienen un coste alto;
los devices con poca RAM pueden recibir jetsam. La descarga inicial por rangos
del update macOS necesita acceso a la CDN Apple; cambios del paquete/hash deben
fallar explícitamente. Se conservan los bounds del guest y timeouts, pero no
hay botón de cancelación durante una llamada nativa; esto sigue experimental.
Los permisos RW→RX por sí solos no prueban SAP ni autorización del backend.
No hay iPhone/iPad conectado a este entorno. La evidencia jailed proviene del
reporte del usuario con dev.1; dev.2 y la autenticación requieren nueva prueba.

MIT de ipatool conservada. Unicorn/QEMU TCI incluyen GPL; cada artifact entrega
sus fuentes modificadas correspondientes, y el source de WaffleStore está en
el fork. Upstream WaffleStore no contiene licencia de redistribución: es un
problema pendiente que requiere aclaración de sus titulares, no un permiso que
pueda inferirse del nombre open source. Los assets Apple son propietarios y
no se redistribuyen. Ver [THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).
