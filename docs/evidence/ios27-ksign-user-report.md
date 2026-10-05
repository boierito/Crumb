# Evidencia de dispositivo — reporte del usuario

- Build: 23001 / 2.3.0-dev.1.
- iOS: 27.0.1; familia iPhone, modelo específico no informado.
- Firma: ksign y certificado; usuario confirma uso sin JIT.
- No se recibió crash report; no se ejecutó una prueba con credenciales.
- Configuración Release, debugger y cierre completo entre pruebas no se
  confirmaron explícitamente. No se infieren de los errno ni del reporte.

Dos diagnósticos con SAP iniciado reportaron:

```text
signed-text-control=42
rw-allocation-errno=0
rw-to-rx-errno=0
rwx-allocation-errno=0
unsigned-code-execution=not-attempted
tci-guest-status=0
tci-guest-rax=42
tci-instruction-hooks=4
tci-test-is-sap=false
keychain-identity-stable=true
store-bag=validated; sap-version=200
sap-initialization=completed
X-Apple-ActionSignature=generated; base64-length=668; contents=withheld
sap-runtime=experimental-tci-static-library
apple-login=not-attempted
```

Esto aporta evidencia de ejecución del guest SAP dentro de esa instalación.
No prueba aceptación de login, validez del token, persistencia entre launches,
compatibilidad con todo método de firma ni funcionamiento de download.
El valor identity-stable compara dos lecturas dentro de la misma ejecución.
