# Crumb / FairPlay — estado técnico, pruebas y próximos pasos

**Fecha de corte:** 25 de octubre de 2026  
**Proyecto:** `boierito/Crumb`  
**Rama de trabajo:** `fix/fairplay-consumer-io`  
**Última IPA probada en dispositivo:** build **23030**  
**Estado global:** importación nativa del keybag aceptada; **recuperación de inscripción portátil y descifrado NO implementados ni validados**.

> **Alcance de este documento.** Este informe es un registro de ingeniería reproducible, basado en los builds, las pruebas sintéticas, los reportes sanitizados V1–V8 y las implementaciones de esta rama. No contiene keybags, SIDB, SINF, DSID, identificadores de hardware, tokens, capturas privadas, material de claves ni datos de cuenta. Una devolución de código nativo `0` no es prueba de que una IPA se haya descifrado.

## 1. Resumen ejecutivo

El objetivo del laboratorio es investigar la compatibilidad de la ruta de importación FairPlay de una transacción de App Store y, por separado, evaluar si existe una **inscripción portátil validada** que permita descifrar una IPA correctamente.

Se logró:

1. Conservar una única sesión de StoreAgent emulado desde la generación de KBSync hasta el consumo del keybag correspondiente, sin recrear el contexto antes de consumirlo.
2. Reproducir dentro del emulador los accesos necesarios a un sistema de archivos virtual `SC Info`, incluyendo `F_GETPATH`, `lstat`, `FSRef` y las operaciones Carbon sobre metadatos.
3. Detectar el rechazo nativo `-50` causado por `FSSetCatalogInfo` con máscara `0xC00` y corregirlo implementando los campos `0x400` (permisos) y `0x800` (Finder info de 16 bytes).
4. Obtener por primera vez en **V6 / build 23028** `emulator_completed=true`, `native_status=0` e `import_succeeded=true`; las capturas posteriores V7 y V8 reprodujeron ese resultado.
5. Validar los ZIP de captura mediante integridad SHA-256 y tamaños del manifiesto, con resúmenes sin bytes privados.
6. Separar de forma explícita **importación del keybag**, **SIDB legado observado**, **raw state portátil**, **enrollment frame** y **modelo portátil validado**; los tres últimos siguen sin verificarse.
7. Mostrar el estado de estas cinco etapas en la interfaz y evitar análisis/exportaciones idénticos repetidos para la misma transacción (build 23030).

**No se ha logrado:** generar el modelo de inscripción portátil compatible con el recuperador actual, autenticar páginas descifradas, producir una IPA realmente descifrada, ni verificar instalación/apertura de una IPA descifrada.

## 2. Arquitectura y límites de compatibilidad

### Ruta implementada

```text
Apple Account autenticada en Crumb
   │
   ├─ StoreAgent retenido en intérprete / contexto original
   │    └─ genera KBSync (196 bytes en las pruebas observadas)
   │
   ├─ transacción de descarga de App Store
   │    └─ respuesta de entitlement + keybag + SINF
   │
   ├─ validación de KBSync y keybag en el mismo contexto
   │    └─ invocación del consumidor nativo emulado
   │         └─ FS virtual aislado (SC Info, SIDB y lock file)
   │
   ├─ resultado del consumidor (V6–V8: status 0)
   │
   └─ ZIP de captura PRIVADO + análisis SANITIZADO
        ├─ comprobación del manifiesto e integridad
        ├─ estado y tamaños; nunca material secreto en el log
        └─ indicadores de readiness individuales
```

### Ruta todavía NO implementada

```text
SIDB legado observado
   ──[puente NO demostrado / no asumir equivalencia]──> raw state portátil
   ──[orquestación NO implementada]──────────────────────> enrollment frame
   ──[validación para build/cuenta/SINF pendiente]───────> modelo de inscripción
   ──[descifrado y autenticación pendientes]─────────────> IPA en plaintext
```

La ruta portátil fijada corresponde a **DrWhax/ipatool**, revisión `ef98227e061838320424744d21a513baba7ade28`. Los evaluadores `state`, `frame` y `SIDB` son componentes específicos de un build de `fairplayd`, con pruebas sintéticas públicas. El proyecto de origen exige contexto de inscripción privado y comprobación de compatibilidad; la salida de StoreAgent 10.9 **no equivale de forma demostrada** a una captura de `fairplayd` moderno. El emulador SAP/StoreAgent y los evaluadores modernos son rutas distintas.

El `SC Info.sidb` virtual contiene datos opacos cuya longitud varía entre transacciones. **No** cortar arbitrariamente 16, 192 o 224 bytes del SIDB para pasarlos como frame, raw state o clave. Tampoco declarar `ready_to_decrypt=true` por la sola existencia del archivo o por `native_status=0`.

## 3. Evolución por build

| Build | Implementación / resultado principal |
| --- | --- |
| **23023** | Capturas estructuradas, análisis del manifiesto, validación del SINF, conservación del productor StoreAgent y contexto de generación privado. |
| **23024** | Invocación de la rutina real de consumo del keybag en el StoreAgent emulado, con filesystem virtual aislado y `consumer-state.json` privado. |
| **23025** | Soporte para `F_GETPATH`, `lstat`, `FSRef` y metadatos Carbon; el replay offline avanza, pero aún rechaza la importación (ejemplo: `-42052`). |
| **23026** | Diagnóstico por etapas, trazas acotadas de virtual I/O y distinción explícita entre llamada completada, importación aceptada y modelo portátil. |
| **23027** | Corrección del contador Swift de archivos/bytes cuando JSON mezcla base64 y `null`; trazas clasificadas del error `FSSetCatalogInfo`. |
| **23028** | Compatibilidad para `FSSetCatalogInfo(0xC00)`; lectura y persistencia de permisos y Finder info; primera aceptación del keybag en V6. |
| **23029** | `Verify current capture — no import required`, análisis local de integridad y `enrollment_bridge` sanitario dentro del diagnóstico. |
| **23030** | Panel `Enrollment readiness`, deduplicación del análisis/exportación por transacción y regresiones para SIDB de tamaños variables. |

Commits relevantes de la rama:

- 23027: `38b33007` (conteo Swift), `0a906e6d` (traza de máscara).
- 23028: `52b5d6c6`, `877c124c`, `0d0e275a` (compatibilidad Carbon, lectura y pruebas).
- 23029: `138b2d2d` (diagnóstico de puente), `407762df` (análisis local), `2c3c48f2` (regresión).
- 23030: `cebab512` y `419ff9a0` (panel y caché), `f3997750` (tamaños variables), `c0866883` (documentación/build final).

## 4. Cronología de capturas y mediciones observadas

Los ZIP de captura privados **no están incorporados al repositorio**. La siguiente tabla utiliza únicamente metadatos permitidos de las pruebas reportadas.

| Control | Build | KBSync (B) | Keybag (B) | SIDB virtual (B) | Estado nativo | Evidencia / interpretación |
| --- | ---: | ---: | ---: | ---: | ---: | --- |
| **V1–V3** | 23023–23026 | — | — | — | Rechazos en pruebas antiguas | Exploración inicial; el replay offline no preservaba necesariamente el contexto de transacción. Véase `FAIRPLAY_LAB.md`. |
| **V4** | 23026 | 196 | 904 | 865 | `-50` | `FSSetCatalogInfo` negativo. El UI decía 0 archivos por una conversión Swift incorrecta de un JSON con `null`; el ZIP tenía 2 entradas. |
| **V5** | 23027 | 196 | 888 | 878 | `-50` | Traza: `unsupported_catalog_mask`, máscara `3072/0xC00`. |
| **V6** | 23028 | 196 | 904 | 862 | **0** | Primera importación aceptada; `FSSetCatalogInfo(0xC00)` retorna `zero`. |
| **V7** | 23029 | 196 | 904 | 883 | **0** | Importación aceptada; análisis local muestra `missing_portable_enrollment_source`; se observaban logs duplicados al verificar/exportar. |
| **V8** | 23030 | 196 | 920 | 885 | **0** | Importación aceptada, panel de cinco indicadores coherente: 2 verdes y 3 pendientes. |

Para V4–V8 los reportes describen **2 archivos de FS virtual**: un lock file vacío y un SIDB no vacío. En V6–V8 se registraron **21 operaciones de I/O virtual** y éxito de las operaciones Carbon relevantes. Las capturas V6–V8 verificaron **5/5 entradas del manifiesto** por tamaño y hash SHA-256. El SINF observado tiene disposición actual, con **440 bytes en su campo privado** y `sinf_count=1`. No se infieren contenidos de esos 440 bytes.

El análisis de la misma aplicación de prueba ha mostrado **2 imágenes Mach-O cifradas** y **712 páginas cifradas** (708 + 4), junto con `supf_count=2`. Este dato mide la IPA **original inspeccionada**: no implica que el consumidor haya intentado descifrar esas páginas.

### Secuencia causal comprobada

1. Inicialmente faltaba compatibilidad del FS virtual; se observaron fallos `-42023` y `-42052` en contextos de replay, además de rechazos `-50` en dispositivo.
2. La traza V5 aisló el error de `FSSetCatalogInfo(0xC00)`: `0x400` representa permisos y `0x800` la información Finder.
3. El emulador pasó a aplicar y conservar los campos requeridos, incluidos los 16 bytes de Finder info, sin aceptar máscaras desconocidas.
4. V6, V7 y V8 devolvieron **`native_status=0`** en transacciones reales retenidas, por lo que el rechazo de la llamada al consumidor quedó resuelto para esos casos.
5. El cuello de botella posterior es **la inscripción moderna y su validación criptográfica**, no otro código negativo del consumidor.

## 5. Significado de indicadores y códigos

| Campo | Lectura correcta |
| --- | --- |
| `consumer_emulator_completed=true` | El intérprete completó la llamada sin error interno; no garantiza éxito del consumidor. |
| `consumer_native_status=0` | El wrapper observado devolvió código nativo cero. |
| `consumer_import_succeeded=true` | Clasificación de Crumb: intérprete completado + retorno nativo cero. **No acredita descifrado.** |
| `legacy_sidb_observed=true` | Se observó un archivo SIDB virtual con bytes. **No es un frame validado.** |
| `portable_raw_state_verified=false` | Falta una fuente de raw state específica y validada. |
| `portable_enrollment_frame_verified=false` | No se ha obtenido ni validado el frame portátil del perfil objetivo. |
| `validated_portable_model=false` | No existe todavía un modelo compatible demostrado para la cuenta, build y SINF. |
| `ready_to_recover=false`, `ready_to_decrypt=false` | **Bloquear** recuperación y descifrado automático en ausencia de las pruebas requeridas. |

Los vectores sintéticos `state`, `frame` y `sidb` con resultado `passed` verifican la implementación de los evaluadores públicos, **no** una inscripción real. Los cambios de longitud del SIDB (865, 878, 862, 883 y 885 B) impiden asumir un tamaño constante como contrato del puente.

## 6. Ubicaciones relevantes en el código

| Archivo | Responsabilidad |
| --- | --- |
| `MapleSyrup/NativeSAP/gueststate.go` | FS virtual, `FSRef`, `FSGetCatalogInfo`, `FSSetCatalogInfo`, permisos/Finder info y traza acotada. |
| `MapleSyrup/NativeSAP/gueststate_test.go` | Pruebas de ABI Carbon, FS virtual, máscaras permitidas y protección de trazas. |
| `MapleSyrup/NativeSAP/kbsyncretained.go` | Sesión emulada retenida; consumo del keybag en el contexto original. |
| `MapleSyrup/NativeSAP/kbsyncbridge.go` | API de intercambio KBSync/keybag y resultado privado del consumidor. |
| `MapleSyrup/NativeSAP/captureanalysis.go` | Validación del ZIP, integridad, límites, SINF, clasificación de etapas y `enrollment_bridge`. |
| `MapleSyrup/NativeSAP/captureanalysis_test.go` | Regresiones sintéticas de privacidad, estados nativos y tamaños de SIDB. |
| `MapleSyrup/NativeSAP/fairplaylabcore.go` | Operaciones `selftest`, `inspect`, `analyze-capture`, `export` y `decrypt` condicionado a modelo. |
| `WaffleStore/App/FairPlayLabView.swift` | Flujo iOS, logs, botones del laboratorio, `Enrollment readiness` y exportación. |
| `scripts/prepare-fairplay-portable.py` | Adaptación de evaluadores portátiles `state`/`frame`/`sidb` fijados en upstream. |
| `docs/FAIRPLAY_LAB.md` | Guía operativa y notas de implementación pormenorizadas por build. |

## 7. Procedimiento actual de reproducción — build 23030

1. Instalar la IPA **23030**, firmada con un certificado válido. La build publicada en GitHub Actions es **sin firmar**.
2. Iniciar sesión en Crumb con una cuenta propia autorizada.
3. Abrir `Settings → FairPlay Laboratory`.
4. Ejecutar `Run self-tests`: `state/frame/sidb = passed`.
5. Usar `Download and capture transaction` con una aplicación gratuita a la que la cuenta tenga acceso.
6. Ejecutar `Inspect IPA`: anotar imágenes y páginas cifradas.
7. Ejecutar `Verify current capture — no import required`.
8. Comprobar el panel: `Native keybag import` y `Legacy SIDB captured` pueden aparecer correctos; `Portable raw state verified`, `Enrollment frame verified` y `Portable model validated` siguen pendientes.
9. Ejecutar `Prepare exports` para obtener `diagnostic.json` (sanitizado) y el ZIP **privado**. Si la captura no cambió, el análisis no debe repetirse innecesariamente.
10. No usar `Import model and decrypt` sin disponer de un modelo propio validado. `Import an older capture ZIP` es opcional y **no** necesario para la transacción actual.

**CI de referencia (build 23030):**

- Build IPA (Debug y Release): https://github.com/boierito/Crumb/actions/runs/40227115134 — `success`.
- SAP protocol tests: https://github.com/boierito/Crumb/actions/runs/40227115142 — `success`.
- Rama: https://github.com/boierito/Crumb/tree/fix/fairplay-consumer-io

La aprobación de CI prueba compilación y regresiones **sintéticas**; las pruebas V6–V8 aportan evidencia en dispositivo exclusivamente para el consumidor legado y la captura.

## 8. Seguridad y tratamiento de datos

- El ZIP `Private-FairPlay-capture*.zip` contiene material derivado de cuenta/dispositivo. **Nunca subirlo al repositorio ni adjuntarlo a issues o pull requests públicos.**
- `context.json`, `consumer-state.json`, `keybag.bin`, `kbsync.bin`, `sinf-*.bin` y archivos SIDB reales se consideran privados.
- Los logs y la tabla de este documento únicamente incluyen tamaños, estados, identificadores técnicos públicos, conteos y resultados de integridad.
- El laboratorio crea capturas con permisos restringidos y protección de archivo de iOS; el ZIP de exportación se mantiene privado y se elimina con los datos temporales de la sesión al cerrar/borrar el laboratorio.
- El filesystem de StoreAgent emulado está separado del filesystem host. No se debe debilitar esa separación para forzar un código nativo cero.
- No almacenar ni imprimir claves recuperadas, bytes del SIDB, DSID, identificadores de hardware o direcciones de memoria sensibles en logs públicos.

## 9. Trabajo pendiente y criterios de aceptación

### Prioridad 1 — Determinar la compatibilidad del estado legado

- Establecer **con evidencia diferencial**, y no por similitud de nombre, si existe una correspondencia reproducible entre el `SC Info.sidb` del StoreAgent 10.9 y cualquier estructura que espere `fairplayd` moderno.
- Si no existe tal correspondencia, declarar explícitamente la ruta legado → moderno **incompatible** y estudiar una fuente de inscripción soportada, sin convertir bytes arbitrarios.
- Registrar únicamente la metodología y resultados sanitizados; conservar muestras privadas fuera de Git.

**Aceptación:** especificación respaldada por pruebas de al menos dos transacciones/inscripciones autorizadas, o incompatibilidad demostrada y documentada.

### Prioridad 2 — Fuente validada de raw state y enrollment frame

- Reconstruir o integrar, para el perfil exacto de `fairplayd`, la obtención lícita/autorizada del raw state requerido y la composición de la entrada del evaluador.
- Verificar tamaño, estructura, build hash, asociación con cuenta/inscripción y SINF, sin reutilizar modelos ajenos.
- Comparar salidas del evaluador con vectores independientes de la **misma inscripción**.

**Aceptación:** coincidencia reproducible de los estados y frames esperados en más de una muestra autorizada; las tres banderas portátiles solo podrán cambiar a `true` cuando exista la evidencia correspondiente.

### Prioridad 3 — Recuperación y autenticación de la IPA

- Producir un modelo de inscripción portátil **solamente después** de validar las etapas anteriores.
- Descifrar una IPA para la cual se cuente con autorización y verificar resultado mediante comprobaciones independientes de plaintext y hashes por página.
- Inspeccionar todos los ejecutables cifrados y probar la instalación/ejecución tras firmar la salida.
- No presentar un `cryptid=0` aislado ni una transformación estructural como prueba de contenido correcto.

**Aceptación:** validación criptográfica de todas las páginas relevantes, modelo asociado a la inscripción correcta, empaquetado íntegro e instalación/apertura comprobadas. Hasta entonces **`ready_to_decrypt=false`**.

### Prioridad 4 — Calidad de UX y diagnósticos

- Mantener el panel de readiness sin interpretaciones falsas; mostrar claramente qué pasos están hechos y cuáles faltan.
- Evitar logs duplicados y errores en operaciones opcionales sin keybag, con pruebas de regresión.
- En una nueva build, agregar una prueba de dispositivo **solo cuando exista un cambio funcional pertinente a la etapa aún abierta**.

## 10. Decisión de ingeniería a fecha de corte

**Hito cerrado:** consumidor legado de keybag en contexto retenido, en las transacciones probadas (V6, V7, V8), con retorno cero y persistencia virtual del SIDB.

**Hito abierto y bloqueante:** conversión/obtención del estado de inscripción portátil validado para el perfil moderno, seguida de descifrado y autenticación real.

**Recomendación:** dejar de iterar sobre descargas idénticas de la 23030 hasta disponer de una hipótesis y una implementación verificables para esa etapa. No confundir nuevos SIDB de distinta longitud con evidencia de avance en descifrado.

---

**Referencias internas:** [Guía FairPlay](FAIRPLAY_LAB.md); [acciones de compilación](https://github.com/boierito/Crumb/actions/runs/40227115134); [pruebas de protocolo](https://github.com/boierito/Crumb/actions/runs/40227115142); [seguimiento técnico](https://github.com/boierito/Crumb/issues/1).

## 11. Continuación: regresiones portátiles y alcance multiapp

Se añadió `bash scripts/test-fairplay-host.sh` y el job `fairplay-portable`
(Linux) para comprobar el paquete FairPlay fijado, sus evaluadores públicos y
los contratos de captura sin depender de Xcode. La preparación del paquete se
comparte con la compilación iOS mediante `prepare-fairplay-package.py`; no se
mantiene una segunda implementación del descifrador. La fuente descargada se
verifica por SHA-256. Las pruebas se ejecutan con detector de carreras.

Estos cambios no completan el puente de inscripción ni modifican las banderas de
readiness. No constituyen una nueva build probada en dispositivo.

### Límites concretos que impiden prometer «cualquier app»

- El upstream fijado exigía exactamente un SINF en `DecryptIPA`. La continuación
  local elimina esa restricción: `fairplaymulti.go` asocia cada ejecutable con su
  SINF y resuelve su modelo independientemente, antes de crear la salida. No se
  elige el primer SINF ni se reutiliza una clave global.
- Las pruebas públicas incluyen clasificación Mach-O universal ARMv7/ARM64,
  geometría inválida, alcance de mutaciones y selección de modelos por selector.
  Son regresiones sintéticas; no prueban cobertura de aplicaciones reales.
- El análisis de plausibilidad de instrucciones no autentica todas las páginas.
  El resultado `completed` del puente sigue separado de
  `plaintext_authenticated=false`.
- La inscripción portátil sigue sin fuente validada. Para evaluar la prioridad 1
  hacen falta las muestras diferenciales autorizadas indicadas en la sección 9;
  no se incorporaron muestras reales ni se dedujo un frame cortando el SIDB.

### Implementación multi-SINF

- Asociación exacta `SC_Info/<ejecutable>.sinf` o `.v5.sinf`; si falta, búsqueda
  de un único SINF en el directorio `SC_Info` del bundle contenedor. Un framework
  puede usar la licencia de su app; una extensión no hereda la licencia de la app
  atravesando su frontera `.appex`. Las asociaciones ambiguas se rechazan.
- Cada ejecutable recibe su máscara y derivador propios; sus slices universales
  comparten la licencia asociada. Todos los modelos se resuelven antes de crear
  el ZIP temporal de salida; un modelo faltante aborta la operación completa.
- El reporte privado registra `sinf_member` por imagen y evidencia por ejecutable
  cuando hay varios dominios. No se publica un hash de clave global ficticio.
- El laboratorio permite importar un archivo o una carpeta de modelos `.bin`,
  copiados a un directorio protegido temporal. Se rechazan enlaces y modelos
  vacíos o mayores de 1 MiB; máximo 256 modelos. El backend valida sus selectores.
- Las regresiones usan dos SINF V3 y modelos independientes, verifican las claves
  esperadas para app/extensión, múltiples slices, asociación de framework,
  ambigüedad y rechazo de un registro incompleto. Pasan junto con las suites
  portátiles bajo `-race` en Linux.

Pendiente: compilación iOS y prueba de la importación de carpetas en dispositivo;
validación con transacciones reales; fuente de inscripción portátil y
verificación criptográfica del plaintext. El éxito sintético no acredita esos
pasos ni permite afirmar compatibilidad con cualquier app.

## 12. Ruta moderna de recuperación integrada (continuación local)

Se incorporó por separado el motor de DrWhax/ipatool
`f7242cce58077a4af8db738ed295ee23ad8f067d`, posterior a la revisión del decoder.
Su runtime generado construye el source de inscripción, selecciona el frame y
continúa hasta producir un pack de inscripción ligado al selector. Se ejecuta
como paquete Go aislado dentro del proceso; no lanza comandos ni ejecuta código
invitado. La descarga está fijada por revisión y SHA-256; también se verifica el
modelo público empaquetado antes de cada recuperación.

### Entrada e integración

- Operación nativa `recover-enrollment`: `Source` es un directorio de captura
  **nativa moderna**, `IPA` la IPA objetivo y `Destination` un directorio nuevo.
- La interfaz ofrece `Recover model from native enrollment folder`. Copia al
  sandbox privado únicamente el manifiesto y los archivos de ejecutable,
  hardware, licencias y estado necesarios; no necesita copiar la IPA ni el
  daemon incluidos en una captura. Verifica los hashes y tamaños del manifiesto
  antes de exponer las entradas al filesystem virtual del motor.
- Primero utiliza el conjunto mínimo de estado probado por upstream; puede
  reintentar con los archivos adicionales de la captura. No introduce un segundo
  SINF instalado en el mismo intento. Los archivos de estado se mantienen
  separados de los del StoreAgent legado.
- Contrasta el source derivado con los evaluadores state/frame existentes: un
  candidato transformado debe reproducir el frame seleccionado por el motor.
  Luego valida el pack contra el SINF de cada ejecutable de la IPA objetivo y
  procesa todas las imágenes inventariadas.
- No devuelve claves, frames, bytes de estado ni rutas privadas en el resultado
  del bridge. El runtime se serializa y descarta sus referencias y buffers
  privados después de cada ejecución. Si falla la validación de la IPA, elimina
  el directorio de recuperación generado. En la interfaz se eliminan también el
  modelo y reporte temporales al terminar.

### Alcance que sigue sin poder afirmarse

La captura `consumer-state.json`/SIDB del StoreAgent 10.9 **no es** la nueva
entrada. Hace falta una captura moderna de una inscripción nativa propia (por
ejemplo, la creada por `ipatool fairplay-enroll collect` en un Mac con la cuenta
inicializada en App Store), con su ejecutable, hardware y estado correspondientes.
No se ha demostrado aquí una conversión de la inscripción legado a esa entrada.
Los indicadores del ZIP legado siguen describiendo esa captura; el resultado
moderno se muestra por separado para no atribuirle evidencia de otra fuente.

`source_frame_cross_checked` y `model_sinf_validated` no significan autenticación
de todas las páginas de plaintext. El resultado conserva
`plaintext_authenticated=false`: el decoder fijado comprueba los hashes de la
fuente cifrada y plausibilidad de instrucciones, pero eso no sustituye una
referencia independiente del plaintext. Tampoco hay garantía de soporte de todas
las generaciones de FairPlay, ni prueba en dispositivo o una IPA real de esta
continuación local. No se ha cerrado el criterio «cualquier IPA 100 % descifrada».

La continuación funcional se identifica como build **23031**; la última build
probada en dispositivo sigue siendo **23030**. No se atribuyen a la 23031 los
resultados V6–V8 ni se declara logrado el descifrado real por actualizar la versión.

Verificación local de esta continuación: `scripts/test-fairplay-host.sh` pasó
las suites del decoder, state/frame/SIDB, rawsource y contratos del laboratorio
con `-race`. El parser de Swift no detectó errores de sintaxis en la vista
modificada; esto no sustituye compilación/typechecking con Xcode. Los tests
locales no incluyen una inscripción privada real ni una IPA real descifrada.


## 13. Intento directo desde la captura retenida — build 23032

`Attempt recovery from current capture` conecta el ZIP de la sesión con el motor
portátil. Valida su manifiesto y los hashes de las entradas utilizadas; exige
una sola transacción, importación nativa aceptada y la identidad original de seis
bytes de `context.json`. Pasa el **SIDB completo sin modificar** al filesystem
virtual del runtime, junto con las licencias, SUPF y ejecutable de la IPA.
Cada SINF objetivo debe coincidir por SHA-256 con uno de la transacción capturada.

Esto implementa una **hipótesis de compatibilidad**, no demuestra que el formato
legado sea compatible. No se cortan bytes del SIDB ni se fabrica estado o claves.
El motor debe derivar un source/frame coherentes, validarlos con los evaluadores
anteriores, producir modelos que validen los SINF y procesar todas las imágenes.
Solo después pueden actualizarse los indicadores derivados. Una incompatibilidad
produce `legacy_state_not_compatible` y no publica ninguna IPA de salida.

La operación procesa dominios SINF independientes y deduplica packs idénticos
antes de validar la IPA con el registro resultante. Los modelos temporales y
estados se borran al terminar. El plaintext sigue marcado como no autenticado.

Pruebas locales nuevas: identidad original y SIDB de 883 bytes preservados
íntegros, rechazo de identidad ausente/inválida, validación del frame en un slot
posterior y rechazo de frame distinto. Decoder y contratos del laboratorio pasan
con `-race`; el parser Swift acepta la vista. **No se ha ejecutado una captura
privada real con esta ruta; su compatibilidad y el descifrado real siguen sin
validarse.** La build 23031 aprobó Debug/Release y todas las suites de GitHub; la
23032 agrega este intento directo y requiere su propia compilación y prueba.


## 14. Rechazo temprano observado en 23032; corrección 23033

El diagnóstico sanitizado de dispositivo reportó self-tests correctos,
`native_status=0`, SIDB de 859 bytes, dos archivos virtuales (uno no vacío),
una transacción íntegra y fallos de `recover-capture` en 0–3 ms con el código
inespecífico `validation_failed`. No aporta el nombre del lock ni material
privado; por tanto no demuestra el punto exacto del rechazo ni incompatibilidad
del SIDB con el runtime portátil.

La revisión encontró un rechazo artificial previo al motor: el decoder nuevo
exigía el prefijo `SC Info.` a **todos** los archivos. El contrato del filesystem
virtual permite otros basenames para locks. Un ZIP sintético íntegro con hardware
original de seis bytes, SIDB de 859 bytes y lock vacío llamado `lock` o `.lock`
reproduce el rechazo con la condición anterior. La 23033 elimina ese requisito,
preserva los bytes completos y mantiene los límites y el rechazo de rutas,
symlinks del ZIP, `.`/`..`, backslashes y traversals.

Ahora se retornan etapas/códigos explícitos y sanitizados, por ejemplo:
`recover-capture/identity` + `capture_hardware_length_invalid`,
`recover-capture/state-files` + `state_filename_invalid`,
`recover-capture/ipa-inventory` + `ipa_inventory_failed`,
`recover-capture/portable-runtime` + `legacy_portable_runtime_rejected`.
La causa privada del error nunca se serializa. Los tests comprueban que un lock
sin ese prefijo permite avanzar hasta inventariar la IPA, y que los errores no
exponen un canario privado.

También se retiraron los textos obsoletos de inspección que decían que la única
ruta era importar un modelo o una captura moderna: inspección mide geometría;
recuperación es una operación aparte. Decryption informa un rechazo y remite al
código/etapa de Enrollment.

Decoder y contratos de laboratorio pasan con `-race` tras demostrar primero el
fallo de los casos nuevos con la condición antigua. No se ha probado todavía
esta corrección contra el ZIP privado real ni se ha logrado validar el plaintext
de una IPA real. La compatibilidad legado → portátil sigue abierta.

La compilación nativa deja de recompilar dos veces el enorme CFG para ejecutar
los mismos tests de host: las suites completas rawsource/lab siguen en el job
Linux `fairplay-portable`, y el build nativo compila el motor para iOS y ejecuta
las suites state/frame/SIDB y los checks SAP existentes.
