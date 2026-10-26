# Crumb y UnfairPlay: evaluación del 26 de octubre de 2026

## Resultado

El descifrado íntegramente en iPhone sigue sin demostrarse. La última evidencia de dispositivo termina en bootstrap/raw_unknown_control_flow, text+0xa4be38. Los PR 7 y 8 están abiertos y en draft. El PR 8 corrige publicación de salidas inválidas y repeticiones; no implementa el destino faltante.

Se comparó Crumb 4ef48073a45980699bc29cc2fe7953e6632409e4 con DrWhax/unfairplay c54e751775cbc2a54510c0fcafb0aef6b1619bc9.

## Qué aporta UnfairPlay

- Recuperación desde una captura de un Mac enrolado, con SIDB, SINF/SUPF, ejecutable instalado e identidad de hardware coherentes.
- Evaluador semántico portable y soporte de contratos 25G83, 26A428/26A434 y 26B118 declarados por upstream.
- Para capturas exactas 26A428/26A434, páginas de referencia verificadas por hash. El binario se interpreta como datos; no se ejecuta nativamente.
- Modelos FPSENR que Crumb ya sabe interpretar. El archivo internal/fairplay/current_sinf.go es idéntico byte por byte al de la fuente fijada ef98227 que Crumb incorpora (SHA256 b2fde487a79ad7838a68ae4ca52203855997865145adfd3b4808a15b47fa0998). Esto demuestra igualdad del parser, no interoperabilidad real de una captura privada.

## Qué no resuelve

No contiene un handler explícito para 0x100a4be38 ni demuestra la conversión del SIDB legado de StoreAgent a un enrollment moderno. Su normalización describe SIDB Mac v7000b y reconoce tamaños 630/642 y 993/1003 para fuentes de 192/224 bytes; esto no permite reinterpretar los 863/910 bytes observados en Crumb.

Upstream también declara que procesamiento completo no autentica todo el plaintext; la plausibilidad ARM64 de primeras páginas no equivale a autenticación criptográfica de todas las páginas.

## Alternativa descartada por el usuario: captura en Mac

En un Mac compatible, iniciar sesión en App Store, instalar y abrir una app propia, y descargar su IPA con la misma cuenta. Compilar la revisión examinada:

```sh
git clone https://github.com/DrWhax/unfairplay.git
cd unfairplay
git checkout c54e751775cbc2a54510c0fcafb0aef6b1619bc9
go build -o unfairplay ./cmd/unfairplay
sudo ./unfairplay collect --app /Applications/App.app --ipa /ruta/App.ipa --output /ruta/privada/App.capture
./unfairplay recover --capture /ruta/privada/App.capture --output /ruta/privada/App.recovery
./unfairplay verify --input /ruta/App.ipa --model /ruta/privada/App.recovery/fairplay-enrollment.bin
```

Solo si recuperación y verificación pasan, probar la importación explícita de fairplay-enrollment.bin mediante la operación de modelo importado de Crumb. Mantener el modelo privado. No se ejecutaron estos comandos sobre una cuenta ni una captura real; esta alternativa aún necesita prueba de extremo a extremo en Crumb. No es una solución iPhone-only.

## Trabajo necesario si se mantiene iPhone-only

1. Obtener acceso de desarrollo al binario exacto ya identificado por hash y traducir la región desde 0x100a4be38 junto con sus sucesores alcanzables, preservando memoria, registros y flags. El ejecutor actual solo inspecciona estructura; no produce esa traducción.
2. Comparar la traducción con una referencia independiente usando fixtures sintéticos antes de incorporarla al dispatcher sin JIT.
3. Validar por separado el contrato del SIDB legado y la derivación del frame; superar bootstrap por sí solo no demuestra un modelo válido.
4. Probar una transacción privada en dispositivo y comprobar el contenido de salida. No publicar otra IPA como arreglo del descifrado mientras solo pase compilación.

## Verificación realizada

- UnfairPlay: go test ./internal/fairplay ./internal/recovery pasó con Go 1.26.0 y umask 022. La primera ejecución bajo umask privado alteró la fixture de permisos 0644 y falló; se corrigió el entorno y se repitió.
- Crumb: ocho pruebas Python de inspección de destino y gate de cobertura pasaron.
- Sin captura privada, Mac enrolado ni prueba física de iPhone en este entorno.

Referencias: https://github.com/boierito/Crumb/pull/7, https://github.com/boierito/Crumb/pull/8, https://github.com/DrWhax/unfairplay.

## Decisión y trabajo posterior

El usuario confirmó mantener todo en iPhone, sin jailbreak ni JIT. La ruta Mac anterior queda descartada como solución de producto. El uso del binario del runner de desarrollo no implica que el usuario necesite un Mac, ni captura datos de ninguna cuenta.

Se añadió `scripts/generate-fairplay-next-block.py`: genera una función Go candidata para la región real desde 0x100a4be38, utilizando exclusivamente el binario con hash exacto y el traductor del archivo vendorizado verificado. Exige sección ejecutable, instrucciones completas y terminador dentro de un límite; falla ante instrucciones no soportadas. CI guarda únicamente la función candidata, sin binario ni captura de cuenta.

La función candidata no se registra ni se incluye en una IPA: primero debe revisarse y probarse contra una referencia independiente. Cinco tests del generador pasan, incluyendo ejecución de Go generado frente a Unicorn para 22 estados sintéticos (con desbordamiento y destinos variables). Esas pruebas verifican el mecanismo del generador, no las instrucciones del destino real ni el descifrado.
