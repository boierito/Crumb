# FairPlay en iPhone: implementación 23042 y evaluación de UnfairPlay

El usuario confirmó mantener todo en iPhone, sin jailbreak, sin JIT y sin Mac. No se depende de una captura Mac ni de un modelo generado en Mac por el usuario.

## Cambio implementado

Las pruebas anteriores de 23039/23040 terminaban en bootstrap/raw_unknown_control_flow, text+0xa4be38. Se obtuvo la región real del daemon del runner de desarrollo usando exclusivamente la sección ejecutable del binario thin con SHA256 `1e58d37e9602dd0ca5e24b14762063b6c310248b4d3a94272771d815e4ea1784`.

El generador reutiliza el traductor de la fuente vendorizada `f7242cce58077a4af8db738ed295ee23ad8f067d`, comprobada por hash, y genera Go antes de compilar la aplicación. No ejecuta el daemon ni modifica memoria ejecutable en el iPhone. No usa capturas, cuentas ni secretos. Artefacto de procedencia: https://github.com/boierito/Crumb/actions/runs/40332161032, job fairplay-macos27-profile.

Se incorporó `rawSourceBootstrapBlock100a4be38` al inventario compartido y al paquete compilado. El bloque procesa tablas de bytes, actualiza flags/contador, escribe un byte y obtiene su siguiente destino de una tabla con índice signed. Se preservan las comprobaciones stopPC por instrucción y el comportamiento de errores de memoria. El destino dinámico no se sustituye por una dirección inventada.

Build identificada como 23042. El gate de cobertura ahora encuentra el handler, su registro y pruebas específicas. Esto habilita una prueba de dispositivo; no demuestra descifrado ni cobertura de todos los destinos siguientes.

## Verificación

- Cinco pruebas del generador: sección ejecutable, límites, instrucciones inválidas, hash incorrecto y comparación independiente de Go generado con ARM64 para 22 estados sintéticos.
- Comparación del bloque real traducido con ejecución Unicorn de sus instrucciones ARM64 ensambladas independientemente: 40 estados, todos los registros, NZCV, destino del salto y región completa de salida de 64 KiB. Incluye desbordamientos de 32 bits, límites del contador e índices de dispatch negativos.
- Vectores de referencia públicos en `scripts/fixtures/fairplay-bootstrap-next.json`: diez estados sintéticos con registros, flags y hash de salida. Tests Go adicionales verifican stopPC, rechazo de memoria ausente y signed dispatch.
- El test diferencial compila la función integrada junto con los helpers/state/memoria originales de la fuente vendorizada, evitando una implementación de memoria distinta en el test.
- Como en el runtime existente, la referencia sustituye BRAA por BR. Esto verifica el destino portable; no valida autenticación PAC nativa.
- Ocho pruebas Python existentes de inspección y gate pasaron. La suite Go previa a la incorporación completa pasó; la suite integrada y CI se ejecutan por separado.

## Qué aporta UnfairPlay y qué falta

Se examinó DrWhax/unfairplay `c54e751775cbc2a54510c0fcafb0aef6b1619bc9`. Sus pruebas `internal/fairplay` e `internal/recovery` pasan con Go 1.26.0 y umask 022. El parser `current_sinf.go` es idéntico al de la fuente ef98227 incorporada por Crumb. Eso verifica igualdad de formato, no recuperación real del estado de StoreAgent.

Upstream recupera desde estado coherente de un Mac enrolado. No demuestra que el SIDB legado que obtiene Crumb equivalga a su estado moderno. Esa ruta se descartó como solución de producto por el requisito iPhone-only.

Todavía se necesita una ejecución privada en iPhone de la build nueva para comprobar los sucesores reales del bloque y, después, la compatibilidad de inscripción, el frame y el modelo. No se ha obtenido en este entorno un modelo ni una IPA descifrada. El propio upstream declara que procesamiento completo y plausibilidad ARM64 de primeras páginas no autentican todo el plaintext.

Referencias históricas: https://github.com/boierito/Crumb/pull/7 y https://github.com/boierito/Crumb/pull/8. Implementación actual: https://github.com/boierito/Crumb/pull/9.
