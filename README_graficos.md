README – Generador de Gráficos de Análisis de Rendimiento
==========================================================

NOMBRE DEL SCRIPT:
graficos_metricas.py

DESCRIPCIÓN GENERAL:
----------------------------------------------------------
Sirve para generar automáticamente los gráficos necesarios para el
“Reporte Completo”.

El programa lee los archivos CSV generados por la aplicación principal (main_MS.go):
- metricas_seq.csv  → resultados del modo secuencial
- metricas_spec.csv → resultados del modo especulativo

A partir de ellos:
1. Calcula los tiempos promedio de ejecución y el Speedup global.
2. Genera gráficos comparativos de rendimiento y distribución.
3. Guarda las imágenes en formato PNG.

----------------------------------------------------------
REQUISITOS:
----------------------------------------------------------
- Python 3.8 o superior (recomendado)
- Librerías:
    pandas
    matplotlib

INSTALACIÓN DE DEPENDENCIAS:
----------------------------------------------------------
Ejecute el siguiente comando para instalar las librerías necesarias:

pip install pandas matplotlib

----------------------------------------------------------
ARCHIVOS NECESARIOS:
----------------------------------------------------------
Asegúrese de tener los archivos CSV en la misma carpeta que el script:

metricas_seq.csv
metricas_spec.csv
graficos_metricas.py

----------------------------------------------------------
USO:
----------------------------------------------------------
Ejecute el script desde la terminal o consola con el siguiente comando:

python graficos_metricas.py

Durante la ejecución, el programa mostrará por pantalla los promedios
calculados y el Speedup total, además de confirmar la creación de los gráficos.

----------------------------------------------------------
SALIDA DEL PROGRAMA:
----------------------------------------------------------
Se generarán los siguientes archivos de imagen en la misma carpeta:

1. grafico_promedios.png   → Comparación de tiempos promedio total.
2. grafico_boxplot.png      → Distribución de tiempos (variabilidad de corridas).
3. grafico_evolucion.png    → Evolución del tiempo total por corrida.
4. grafico_speedup.png      → Visualización global del Speedup calculado.

Se eligieron los gráficos de barras, boxplot y evolución temporal porque son los que mejor ilustran las diferencias de tiempo promedio, la consistencia de las corridas y la tendencia general del rendimiento. En conjunto, permiten visualizar claramente la mejora obtenida por el patrón de ejecución especulativa, cumpliendo lo que pide el enunciado.

----------------------------------------------------------
DESCRIPCIÓN DE LOS GRÁFICOS:
----------------------------------------------------------
1. Gráfico de Barras (grafico_promedios.png):
   Muestra la comparación directa del tiempo promedio total entre
   la ejecución secuencial y la especulativa.

2. Gráfico Boxplot (grafico_boxplot.png):
   Presenta la distribución y dispersión de los tiempos de las 30 corridas
   en ambos modos, evidenciando estabilidad y variación.

3. Gráfico de Evolución (grafico_evolucion.png):
   Permite observar cómo cambia el tiempo de cada corrida, mostrando
   la consistencia del patrón especulativo frente al secuencial.

4. Gráfico de Speedup (grafico_speedup.png):
   Resume visualmente el factor de mejora obtenido (Speedup = T_seq / T_spec).

----------------------------------------------------------
INTERPRETACIÓN DE RESULTADOS:
----------------------------------------------------------
- El valor del Speedup indica la mejora de rendimiento:
  Speedup > 1.0 → Mejora
  Speedup = 1.0 → Igual rendimiento
  Speedup < 1.0 → Peor rendimiento (no esperado)

----------------------------------------------------------
NOTAS FINALES:
----------------------------------------------------------
- Asegúrese de ejecutar el script después de haber generado los CSV desde
  la aplicación Go, con las 30 corridas completadas.

----------------------------------------------------------
Script elaborado por:
- Lukas Scheel
- Francisca Meyer
