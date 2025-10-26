CONTROL 2 – EJECUCIÓN ESPECULATIVA VS SECUENCIAL EN GO
=====================================================

ANÁLISIS DE RENDIMIENTO
-----------------------------------------------------

CONFIGURACIÓN DE LAS PRUEBAS:
Se realizaron 30 corridas para cada modo (especulativo y secuencial) con los siguientes parámetros:
n = 300
umbral = 0
dificultad = 5
max_primo = 500000

Estos valores garantizan que la rama A (Proof-of-Work) sea la ganadora en todas las ejecuciones, permitiendo una comparación directa entre estrategias.

RESULTADOS PROMEDIO:
-----------------------------------------------------

Estrategia   | Promedio Total (ms) | Promedio Decisión (ms) | Promedio Rama (ms) | Speedup
Secuencial   | 3171.37             | 226.90                  | 2930.43            | —
Especulativa | 3088.50             | 102.33                  | 3046.57            | 1.03×

Cálculo:
Speedup = 3171.37 / 3088.50 = 1.0268 ≈ 1.03

RESULTADOS:
La ejecución especulativa fue aproximadamente 2.7 % más rápida que la secuencial.
Esto ocurre porque la evaluación costosa de la traza de matrices se solapa parcialmente con la ejecución de las ramas, reduciendo el tiempo total efectivo.
El beneficio es limitado porque la decisión (~100 ms) representa una fracción pequeña del tiempo total (~3 s).
Con valores mayores de n (más costo de decisión) o menor dificultad del PoW, el speedup aumentaría significativamente.

-----------------------------------------------------
DESCRIPCIÓN GENERAL
-----------------------------------------------------

Implementación del patrón de Ejecución Especulativa en Go, donde dos tareas pesadas (ramas A y B) se ejecutan en paralelo mientras el hilo principal calcula una condición costosa.
Una vez determinada la condición, el sistema valida la rama ganadora, cancela la perdedora y registra métricas de tiempos y rendimiento en un archivo CSV.
El código incluye además una versión secuencial para comparar resultados y calcular el Speedup global.

-----------------------------------------------------
RAMAS IMPLEMENTADAS
-----------------------------------------------------

Rama A: Proof-of-Work (PoW)
Descripción: Simula minería de un bloque, buscando un hash con prefijo de ceros. Carga altamente exponencial.
Función: SimularPruebaDeTrabajo()

Rama B: Búsqueda de Primos
Descripción: Encuentra todos los números primos hasta un límite max_primo. Complejidad polinomial alta.
Función: EncontrarNumerosPrimos()

Decisión: Cálculo de Traza
Descripción: Multiplica dos matrices n×n y calcula la traza (suma de diagonal).
Función: CalcularTrazaProductoMatrices()

-----------------------------------------------------
ESTRUCTURA DEL PROYECTO
-----------------------------------------------------

c2-especulativa/
├── go.mod
├── main.go
└── README.txt

main.go: contiene toda la lógica de ejecución (modos secuencial y especulativo).
go.mod: define el módulo y versión de Go.
README.txt: documento con análisis, uso y resultados.

-----------------------------------------------------
DETALLES TÉCNICOS DE IMPLEMENTACIÓN
-----------------------------------------------------

CONCURRENCIA Y CANALES
- Cada rama (A y B) se ejecuta en una goroutine separada.
- El hilo principal usa canales con buffer (1) para recibir resultados.
- La cancelación cooperativa se realiza con context.WithCancel(), mecanismo idiomático en Go para señalizar cancelación entre goroutines.

SINCRONIZACIÓN
- Se utiliza sync.WaitGroup y sync.Mutex para coordinar escritura concurrente en el archivo CSV.
- Las métricas se agregan de forma segura y consistente entre goroutines.

MÉTRICAS REGISTRADAS
modo, ejecucion, inicio_ms, fin_ms, total_ms, decision_ms, rama, rama_ms,
otra_cancelada, detalle, traza, umbral, dificultad, maximo_primo

-----------------------------------------------------
REQUISITOS
-----------------------------------------------------

Go versión 1.22 o superior
Git (para control de versiones)
CPU multinúcleo (para apreciar ejecución paralela)
Sistemas operativos compatibles: Windows / Linux / macOS

Verificar instalación:
go version

-----------------------------------------------------
COMPILACIÓN Y LIMPIEZA
-----------------------------------------------------

BORRADO DE CACHÉ Y DEPENDENCIAS
go clean -cache -modcache -i -r
go mod tidy

COMPILACIÓN DEL EJECUTABLE
go build -o c2spec.exe

-----------------------------------------------------
USO Y EJECUCIÓN
-----------------------------------------------------

La decisión se basa en el resultado de CalcularTrazaProductoMatrices(n) y el parámetro -umbral:
Si traza >= umbral → gana rama A (Proof-of-Work)
Si traza < umbral → gana rama B (Primos)

PARÁMETROS DISPONIBLES
./c2spec.exe -mode both -runs 30 -n 300 -umbral 0 -dificultad 5 -max_primo 500000 -block_data "Bloque" -nombre_archivo metricas.csv

-----------------------------------------------------
CONSEJOS DE USO
-----------------------------------------------------

Para forzar la rama A (PoW): use -umbral 0
Para forzar la rama B (Primos): use un umbral muy alto (por ejemplo 999999999)
Para repetir los mismos resultados aleatorios, use -seed con un valor fijo.

-----------------------------------------------------
COMANDOS DE EJECUCIÓN DEFINIDOS
-----------------------------------------------------

Ejecución Especulativa (30 corridas):
./c2spec.exe -mode spec -runs 30 -n 300 -umbral 0 -dificultad 5 -max_primo 500000 -nombre_archivo metricas_spec.csv

Ejecución Secuencial (30 corridas):
./c2spec.exe -mode seq -runs 30 -n 300 -umbral 0 -dificultad 5 -max_primo 500000 -nombre_archivo metricas_seq.csv

Ejecución Combinada (para análisis automático):
./c2spec.exe -mode both -runs 30 -n 300 -umbral 0 -dificultad 5 -max_primo 500000 -nombre_archivo metricas.csv

-----------------------------------------------------
ANÁLISIS EXTENDIDO (REPORTE COMPLETO)
-----------------------------------------------------

Para el Reporte Completo se recomienda cargar metricas_spec.csv y metricas_seq.csv en una herramienta de gráficos (Excel, Google Sheets o Python/Matplotlib) y visualizar:
- Gráficas de barras: tiempo promedio por estrategia
- Gráfica de línea: evolución del tiempo total por corrida
- Gráfica de Speedup: T_seq / T_spec en función de n

Ejemplo de análisis en Python:

import pandas as pd
import matplotlib.pyplot as plt
spec = pd.read_csv("metricas_spec.csv")
seq = pd.read_csv("metricas_seq.csv")
plt.bar(["Secuencial", "Especulativo"], [seq["total_ms"].mean(), spec["total_ms"].mean()])
plt.title("Comparación de tiempos promedio")
plt.ylabel("Tiempo (ms)")
plt.show()

-----------------------------------------------------
CONCLUSIONES
-----------------------------------------------------

- La ejecución especulativa logra reducir el tiempo total al superponer la evaluación de la condición con la ejecución de las ramas.
- El beneficio depende del costo relativo de la decisión frente a las ramas.
- Con n más alto (decisión más costosa), el speedup crece significativamente.
- En este caso, con n=300, el speedup fue de 1.03x (aproximadamente 3% de mejora).

-----------------------------------------------------
AUTORES
-----------------------------------------------------

Francisca Meyer
Lukas Scheel

-----------------------------------------------------
[ENLACE AL REPOSITORIO](https://github.com/lukasscheelb/LDP-Tarea2)
-----------------------------------------------------

-----------------------------------------------------
NOTAS FINALES
-----------------------------------------------------

Lenguaje: Go
Versión recomendada: 1.22 o superior
