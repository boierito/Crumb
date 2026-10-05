# Testing — 2.3.0-dev.1

**No hay un dispositivo iOS conectado.** Las filas de dispositivo son pendientes,
no resultados inferidos a partir del build o del signer en Linux.

| iOS version | device | login | 2FA | search | versions | purchase | download | export |
|---|---|---|---|---|---|---|---|---|
| 26.x | iPhone físico, pendiente | Bloqueado por fase SAP | Pendiente | Legacy, no validado | Pendiente | No implementado | Pendiente | Pendiente |
| 26.x | iPad físico, pendiente | Bloqueado por fase SAP | Pendiente | Legacy, no validado | Pendiente | No implementado | Pendiente | Pendiente |
| 27.x | iPhone físico, prioridad | Bloqueado por fase SAP | Pendiente | Legacy, no validado | Pendiente | No implementado | Pendiente | Pendiente |
| 27.x | iPad físico, pendiente | Bloqueado por fase SAP | Pendiente | Legacy, no validado | Pendiente | No implementado | Pendiente | Pendiente |

## Evidencias automáticas

Original Debug/Release: https://github.com/boierito/wafflestore/actions/runs/37300555931
El source original d508e53 produjo IPAs sin cambios funcionales.
Los runs finales están en https://github.com/boierito/wafflestore/actions

```sh
swift test --package-path MapleSyrup/SAPKit
python3 -m pip install unicorn==2.1.4
bash scripts/test-no-exec.sh
bash scripts/test-tci-host.sh
# Opcional/manual: hace llamadas SAP a Apple, descarga assets, no login:
bash scripts/test-sap-host.sh
```

Los tests cubren Bag/endpoints/version, envelope Apple, identidad/GUID,
handshake certificado/setup/state 1→0, firma de bytes exactos, sesión cerrada,
rechazo de JIT y firma vacía, error HTTP 204/403/404/429/5xx y límite de respuesta.
Los fixtures FakeGuest **no generan firmas válidas de Apple**. El test de HTTP
no implementa todavía los retries de login; no presentar eso como cobertura de
retries, 2FA, versiones, kbsync o persistencia de cuenta completa.

El test Linux negativo ejecuta Unicorn normal y luego lo restringe en un
subproceso. El positivo ejecuta el mismo TCI/TCIProbe.c usado por la app con
PROT_EXEC prohibido: status=0, RAX=42, instruction-hooks=4. El test manual SAP
usa endpoints del Bag, guest real y firmas diferentes para cuerpos diferentes.
Es evidencia host; no sustituye ejecución arm64 ni aceptación de login.

## Gate de dispositivo: antes de login

1. Descargar artifact **Release**, resignar mediante SideStore/AltStore/
   Sideloadly/certificado habitual. Registrar iOS, modelo, método, bundle ID y
   build. No adjuntar debugger ni activar JIT; no usar TrollStore/jailbreak.
2. Abrir Settings → SAP diagnostic, ejecutar sin red. Esperar
   signed-text-control=42, rw-allocation-errno=0; registrar RW→RX/RWX errno.
   Un valor 0 no significa que el JIT pueda ejecutarse ni que SAP funciona.
3. Esperar tci-guest-status=0, tci-guest-rax=42 y tci-instruction-hooks=4.
   Si falla/crashea, conservar crash report y diagnóstico sanitizado: no
   comenzar el login para ocultar el fallo.
4. Esperar keychain-identity-stable=true. Matar/reabrir y repetir. Comparar
   identidad en debugger está prohibido para el test jailed; el test de
   igualdad dentro de la app y una prueba controlada de persistencia son
   necesarios. El reporte no exporta el GUID; no considerar la igualdad en
   una sola ejecución prueba completa de persistencia entre launches.
5. Activar Initialize SAP and sign a test body; comprobar Bag/version 200,
   setup completo y X-Apple-ActionSignature=generated con longitud no cero.
   Primer arranque puede tardar minutos descargando/verificando assets. No
   registrar certificate/buffers/firma completos.
6. Repetir en frío/caliente, airplane mode, timeout/red con fallos y con la
   app en background/foreground. Registrar error stage y crash/jetsam cuando
   corresponda. Un test no vacío valida generación local, no el login.
7. Repetir en iOS 26 y 27, al menos iPhone; iPad para presentación del Share
   Sheet. Adjuntar diagnósticos sanitizados y fechas en esta matriz.

Stage nativo: 1 argumentos; 2 assets/cache/CDN; 3 runtime/emulator;
4 init guest; 5 exchange; 6 sign; 7 handle. Errores HTTP/setup son independientes.

## Matriz de pruebas que se habilita después del gate SAP

| Flujo | Caso positivo | Casos de error obligatorios | Estado de esta entrega |
|---|---|---|---|
| Login | Plist firmado, redirects preservan cuerpo, DSID/token/storefront | password incorrecta, disabled/locked, HTML/204, 403/404/429/5xx, timeout/cancel | Backend moderno pendiente |
| 2FA | Código vigente, espacios normalizados | incorrecto, expirado, seis dígitos inválidos, demasiados intentos | Pendiente |
| Sesión | Reabrir sin password/código guardados; renovar/pedir credenciales | token inválido/expirado, Keychain inaccesible, logout sin cookies/tokens | Sólo identidad nueva implementada |
| Search | Nombre/App Store ID/Bundle ID | sin resultados, región, offline, lookup inválido | Código original conservado |
| Versions | Visible version ↔ externalVersionId verificado por IPA | metadata desactualizada, ID inválido, versión no disponible | Pendiente |
| Purchase | App gratuita/licencia existente | paga rechazada, usuario debe aceptar términos/interacción | Pendiente |
| Download | Última y antigua, kbsync y versión/ID correctos | blob cache rechazado, token refresh, redirects, CDN/ZIP corrupto, poco disco, timeout | Pendiente |
| Export | IPA verificada en Documents, Files/Share Sheet | archivo ausente, cancel, escenas/iPad popover | Flujo original conservado, nuevo pendiente |
| Instalación | Método externo soportado | IPA cifrada/no resignable, OS rechaza manifest | No certificado; separado de download |

No probar compras pagas automáticas. No marcar éxito al abrir una página de
instalación: descargar/exportar y la instalación confirmada por iOS son estados
distintos. No declarar downgrade funcional por un build o una firma de prueba.
