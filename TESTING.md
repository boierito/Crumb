# Testing — 2.3.0-dev.2

No hay un iPhone conectado al entorno. El usuario reportó SAP setup y firma en
un iPhone iOS 27.0.1 con ksign/certificado sin JIT, dos veces, build 23001.
Ver [reporte](docs/evidence/ios27-ksign-user-report.md). La igualdad de identidad
reportada es dentro de una ejecución; persistencia entre launches pendiente.
El reporte no prueba aceptación de login. Dev.2 requiere nueva prueba real.

| iOS version | device | login | 2FA | search | versions | purchase | download | export |
|---|---|---|---|---|---|---|---|---|
| 26.x | iPhone, pendiente | Implementado; pendiente real | Implementado; pendiente real | Conservado; pendiente | Pendiente | No portado | No portado | Original; nuevo pendiente |
| 27.0.1 | iPhone, ksign/certificado; modelo no informado | Implementado dev.2; pendiente real | Implementado dev.2; pendiente real | Conservado; pendiente | Pendiente | No portado | No portado | Original; nuevo pendiente |
| 27.x | iPad, pendiente | Implementado; pendiente real | Implementado; pendiente real | Conservado; pendiente | Pendiente | No portado | No portado | Pendiente |

## Tests automáticos

```sh
swift test --package-path MapleSyrup/SAPKit
python3 -m pip install unicorn==2.1.4
bash scripts/test-no-exec.sh
bash scripts/test-tci-host.sh
# Manual/opcional: Apple assets + SAP setup, sin credenciales/login:
bash scripts/test-sap-host.sh
```

SAPTests cubre Bag, endpoints, identidad/GUID, envelope Apple, handshake, firma
de bytes exactos, cierre y rechazo de runtime JIT/firma vacía/respuestas grandes.
AuthenticationTests usa un signer fixture **sin firmas válidas de Apple** y
respuestas locales. Cubre plist/headers, POST/body/attempt tras redirects,
rechazo de hosts/rutas inseguros y 303, límite de cuatro redirects, -5000 y
attempt 2, challenge 2FA/cookies en memoria y rechazo/normalización de códigos, 204/403/404/429/5xx,
timeout/cancelación, Retry-After/backoff, respuestas incompletas, mensajes
redactados, fallo de persistencia, serialización de cuenta sin password/código,
cookies/expiración/host, GUID mismatch, logout y Keychain real en macOS.
El test macOS de Keychain usa un service UUID con datos fixture y lo elimina.
No certifica Keychain bajo otro certificado/access group ni iOS físico.

Original Debug/Release: https://github.com/boierito/wafflestore/actions/runs/37300555931
Los resultados de dev.2 están en https://github.com/boierito/wafflestore/actions.
Build es una condición necesaria, no prueba de login aceptado.

## Instalación y prueba real del login (prioridad iOS 27)

1. Usar la IPA **Release 23002 / dev.2**. Resignar con el mismo certificado y
   bundle ID si se quiere conservar la identidad/caches de dev.1. Instalar
   mediante ksign/SideStore/AltStore/Sideloadly habitual, sin debugger ni JIT.
   Registrar modelo, iOS, método y configuración Release. No compartir secretos.
2. Al abrir se elimina el almacenamiento legacy (authinfo/.authkey/clave EC).
   Si nunca hubo login moderno, aparecerá el formulario. No importar sesiones
   antiguas ni concatenar un código manualmente al password.
3. Introducir Apple ID/password propios → Send 2FA Code. Esperar etapas Bag,
   SAP, firma y Contacting Apple; un inicio frío puede tardar minutos por assets.
   El botón y campos se bloquean durante el request, la UI sigue respondiendo.
4. Sólo si Apple requiere 2FA aparece el campo. Usar un código de seis dígitos
   vigente del dispositivo confiable → Log In. Puede haber login sin challenge.
5. Esperar estado Signed in y reporte outcome=authenticated,
   DSID/passwordToken/storefront=present, valores withheld. Esto es recepción
   de credenciales reales, a diferencia de una firma de prueba. Pod es opcional.
6. Settings → Export authentication diagnostic. Compartir ese texto y el error
   visible si falla (revisar/redactar datos personales). No enviar logs HTTP,
   plist, password, códigos, cookies, token o firma.
7. Cerrar/reabrir app. Debe cargar cuenta sin password ni código. El estado
   Saved session loaded significa sólo restauración local, no token validado.
   Logout desde menú debe volver al formulario y borrar cuenta/cookies sin
   cerrar la app. Volver a iniciar sesión para confirmar funcionamiento.
8. El botón Downgrade está deshabilitado deliberadamente en esta build.
   La búsqueda/favoritos/historial se conserva; Store/kbsync/download viene después
   de comprobar el login real. No registrar descarga como exitosa por login.

## Casos reales que siguen pendientes

| Caso | Resultado esperado | Estado |
|---|---|---|
| Login correcto sin 2FA | Token/DSID/storefront guardados; password/código vacíos | Fixture; pendiente Apple real |
| 2FA correcto | Challenge → código → autenticado | Fixture; pendiente Apple real |
| Password incorrecta | Dos intentos lógicos máximo; mensaje Apple útil | Fixture; pendiente Apple real |
| Código incorrecto/expirado | Mensaje visible; código borrado; permite código fresco | Fixture; pendiente Apple real |
| Cuenta bloqueada/disabled | Error explícito, no sesión creada | Fixture; pendiente Apple real |
| 429/503/HTML/204 | Tres envíos máx., backoff, Retry-After, error explícito | Fixture; no provocar rate limits reales |
| Cancelar durante red/backoff | No aplica resultado a UI; puede volver a iniciar | Fixture transporte; pendiente UI real |
| Cancelar durante guest C/Go | Espera retorno nativo antes de teardown | Pendiente real; no promesa de cancelación inmediata |
| Sesión tras relaunch | Cuenta/restored sin password; no claim online | Keychain macOS; pendiente iOS |
| Logout | Cuenta/cookies borradas, identidad conservada | Persistencia fixture/Keychain macOS; pendiente UI real |
| Identidad entre launches | Mismo item machine-identity con misma firma/bundle ID | Pendiente iOS real |
| Token inválido/expirado | Pedir nuevo login, sin password guardado | Detección Store pendiente |
| Search/version/purchase/download/export | Última/antigua con externalVersionId, progreso/export independientes | Store moderno pendiente |

## Diagnóstico SAP de dev.1 / regresión dev.2

Settings → SAP diagnostic: primero sin red, luego Initialize SAP and sign a test
body. Esperar tci-status=0, RAX=42, hooks=4; Bag v200/setup completo; firma no vacía.
RW→RX/RWX errno=0 no implica que código unsigned se ejecutó. No compartir buffers.
Stage nativo: 1 argumentos; 2 assets/cache/CDN; 3 runtime; 4 init; 5 exchange;
6 sign; 7 handle. En crash/jetsam guardar reporte sanitizado y memoria/dispositivo.
Repetir cold/warm, offline, app background/foreground en iOS 26/27 e iPad.

No automatizar compras pagas. Una IPA App Store puede seguir cifrada: no prometer
instalación/downgrade porque la app WaffleStore se pueda resignar normalmente.
