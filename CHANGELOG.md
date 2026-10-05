# Development changelog

## 2.3.0-dev.4

- Corrige rechazo de ent/download: regenera caché rechazada y prueba pod validado.
- No interpreta HTTP 401 sin error Apple como sesión expirada ni pide logout.
- Conserva errores explícitos de sesión/licencia devueltos por el pod.
- Login con conexiones separadas y cookies efímeras compartidas entre intentos.
- Admite respuestas Document/Protocol con pares plist sin contenedor dict.
- Diagnóstico v4: HTTP, fallo Apple numérico y categoría estable, sin secretos.
- Explica que Apple decide si requiere 2FA; no fuerza un challenge nuevo.
- Tests de recuperación, clasificación, cookies y parsing; prueba iOS pendiente.

## 2.3.0-dev.3

- Registra login/2FA y reapertura confirmados por usuario en iOS 27.0.1 (23002).
- Porta StoreAgent kbsync a TCI, sin decryption/JIT; cache aceptado en Keychain.
- Store Swift async: Bag, storefront/country, ent/download y fallbacks del ipatool actual.
- Versiones por externalVersionId; número visible leído del Info.plist de la IPA.
- Adquisición exclusiva de apps verificadas gratuitas; no compras pagas/subscriptions.
- CDN aislada sin secretos, progreso real, retries y validación ZIP/CRC/MD5/identidad.
- Empaquetado streaming con extras Apple y metadata/SINF, sin force unwraps.
- IPA persistente en Documents/Downloads, Files/Share Sheet; instalación separada.
- Búsqueda: cancelación esperada silenciada, query bien codificada y país de la cuenta.
- Tests Store/ZIP/Range y diagnóstico sanitizado; descarga física dev.3 aún pendiente.

## 2.3.0-dev.2

- Registra prueba SAP en iPhone iOS 27.0.1, ksign/certificado sin JIT.
- Conecta login firmado y 2FA a la UI existente, con estados async/cancelación.
- Adapta redirects/retries/Retry-After del ipatool moderno y errores Apple.
- Sesión/cookies en Keychain; no persiste password ni código; descarta legacy.
- Logout sin regenerar identidad ni terminar la app; diagnóstico sanitizado.
- Login real todavía pendiente de validación. Store/download sigue pendiente.

## 2.3.0-dev.1

- Conserva el upstream original y genera IPAs de referencia en macOS CI.
- Audita el backend original y el ipatool del 1 de octubre, incluido kbsync.
- Añade protocolo SAP v200 Swift, identidad Keychain y bridge del guest SAP.
- Restaura experimentalmente Unicorn/TCI para interpretar sin memoria ejecutable.
- Añade diagnóstico de dispositivo y prueba de firma SAP sin credenciales.
- Login, 2FA, purchase y download modernos pendientes de validar SAP jailed.

No es una release estable ni una afirmación de recuperación de WaffleStore.
