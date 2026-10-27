# FairPlay: selección, extracción y contrato de inscripción

Base examinada: `a32bb48351133be493eedbe3236a70380ffcc9d2` (build 23045).
Esta rama es un cambio **funcional** de momento de captura, no una
declaración de descifrado ni un parche por dirección ARM64.

## Evidencia y límites

La prueba privada sanitaria del iPhone informa ejecución `completed`,
13.683.433 instrucciones genéricas y `Source`/ `Frame`/ `Pack` de
224/16/4184 bytes. `state.Evaluate` terminó, pero ninguno de los 14
bloques transformados coincidió con `frame.Evaluate`. Por ello no se
escribió modelo ni se verificó la IPA. Un estado nativo 0 del consumidor
StoreAgent no garantiza inscripción del daemon FairPlay moderno.

El runtime vendorizado de `unfairplay` delimita la fase
`enrollment-frame` **después de selección** y separa `enrollment-sinf`.
Su `runRawSourceSelectionCFG` materializa el frame en la segunda llamada
`selectionFrameCall`, mediante `applyRecoveredEnrollmentFrame`. Crumb
anteriormente leía el par después de terminar SINF; eso permite leer
modificaciones posteriores, aun cuando ambos getters devuelven bytes válidos.

El SIDB legado de StoreAgent observado (888 bytes) no corresponde a las
longitudes conocidas de SIDB nativo macOS 26A428 documentadas por
`canonicalSIDB` upstream (630/642 o 993/1003). **El tamaño no basta para
probar incompatibilidad**: es una ausencia de contrato verificado.
Crumb actualmente monta el SIDB legado en dos ubicaciones virtuales de
macOS y registra cinco búsquedas de estado no satisfechas. Sin saber cuáles
eran obligatorias no deben sustituirse por archivos ficticios.

## Cambio

1. Cuenta, bootstrap, writer y selección ejecutan el mismo código previo.
2. Inmediatamente después de selección se copian `Source` y `Frame`;
   cualquier fallo al extraerlos hace fallar la recuperación.
3. SINF se ejecuta sobre **el mismo estado** y genera el `Pack` sin
   alterar el checkpoint ya copiado. Si SINF falla, el checkpoint se borra.
4. Se releen ambos buffers tras SINF y se informa sólo
   `checkpoint_compared_after_sinf`, `source_changed_after_sinf` y
   `frame_changed_after_sinf`, sin direcciones, hashes ni datos privados.
5. El par **sigue sujeto al mismo** `state.Evaluate` +
   `frame.Evaluate` de 14 candidatos, binding SINF, modelo y validación
   completa de la IPA. El cambio no prueba que el SIDB legado sea apto.

La regresión usa el modelo real embebido, memoria sintética, evaluadores
reales e inyección sintética de escrituras posteriores en ambas direcciones.
También verifica que una comparación imposible no se reporte como estable.
No utiliza datos de Apple Account ni dispositivos.

## Interpretación futura, sin inferir éxito de un build

| Checkpoint | Resultado | Hipótesis que respalda |
| --- | --- | --- |
| Comparado; source/frame cambió | Par válido tras selección; falla el leído tardío | La fase SINF sobrescribe parte del output; comprobar después el resto de validaciones |
| Comparado; ambos iguales | Sigue el mismatch de los 14 candidatos | El problema no es la sobrescritura tardía; focalizar SIDB, fidelidad del intérprete y entradas faltantes |
| Comparación no disponible | No se puede juzgar estabilidad | Inspeccionar el fallo de lectura, no interpretar como igualdad |
| Par coherente y falla validación de SINF/IPA | El checkpoint corrigió como mínimo una puerta | Todavía no hay desencriptación validada |
| Modelo y IPA completamente validados | El puente funcionó para esa transacción | Repetir control de otra IPA, sin trasladarlo a cuentas/builds distintos |

## Qué exigir antes de probar otra IPA en iPhone

El test del checkpoint y las suites completas deben pasar en CI. Si el
dispositivo se vuelve a usar, **una sola** recuperación sanitaria comprobará
si los buffers cambian entre selección y SINF y si el par capturado antes
pasa la validación independiente. No exportar SIDB, SINF, claves, keybags,
estado, identificadores ni IPAs personales al repositorio.

No se publica una IPA como «descifrada» sin `target_ipa_validated=true`.
La build 23045 permanece como línea base y no debe renombrarse.
