package main

import (
	"context"
	"crypto/sha256"
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"
)

// -------------------------------
// Tipos y modelos de datos
// -------------------------------

type Resultado struct {
	Etiqueta string
	Salida   string
	InicioEn time.Time
	Duracion time.Duration
	Error    error
}

type Metricas struct {
	Modo             string
	Ejecucion        int
	InicioUnixMs     int64
	FinUnixMs        int64
	TotalMs          int64
	DecisionMs       int64
	Rama             string
	DuracionRamaMs   int64
	OtraCancelada    bool
	HashOUltimoPrimo string
	Traza            int
	Umbral           int
	Dificultad       int
	MaximoPrimo      int
}

// -------------------------------
// Carga de trabajo (con cancelación)
// -------------------------------

func SimularProofOfWork(ctx context.Context, datosBloque string, dificultad int) (string, int, error) {
	objetivo := strings.Repeat("0", dificultad)
	for intento := 0; ; intento++ {
		// chequeo de cancelación ligero
		if intento%10000 == 0 {
			if err := ctx.Err(); err != nil {
				return "", -1, err
			}
		}
		hash := fmt.Sprintf("%x", sha256.Sum256([]byte(datosBloque+fmt.Sprint(intento))))
		if strings.HasPrefix(hash, objetivo) {
			return hash, intento, nil
		}
	}
}

func EncontrarPrimos(ctx context.Context, maximo int) ([]int, error) {
	if maximo < 2 {
		return []int{}, nil
	}
	primos := make([]int, 0, maximo/10)
	for i := 2; i < maximo; i++ {
		if i%1024 == 0 {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}
		esPrimo := true
		for j := 2; j*j <= i; j++ {
			if i%j == 0 {
				esPrimo = false
				break
			}
		}
		if esPrimo {
			primos = append(primos, i)
		}
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	return primos, nil
}

func CalcularTrazaProductoMatrices(n int) int {
	matriz1 := make([][]int, n)
	matriz2 := make([][]int, n)
	for i := 0; i < n; i++ {
		matriz1[i], matriz2[i] = make([]int, n), make([]int, n)
		for j := 0; j < n; j++ {
			matriz1[i][j], matriz2[i][j] = rand.Intn(10), rand.Intn(10)
		}
	}
	traza := 0
	for i := 0; i < n; i++ {
		suma := 0
		for k := 0; k < n; k++ {
			suma += matriz1[i][k] * matriz2[k][i]
		}
		traza += suma
	}
	return traza
}

// -------------------------------
// Utilidades
// -------------------------------

func ejecutarTarea(ctx context.Context, etiqueta string, f func(context.Context) (string, error)) Resultado {
	inicio := time.Now()
	salida, err := f(ctx)
	return Resultado{Etiqueta: etiqueta, Salida: salida, InicioEn: inicio, Duracion: time.Since(inicio), Error: err}
}

func asegurarCSV(ruta string) (*os.File, *csv.Writer, error) {
	nuevoArchivo := false
	if _, err := os.Stat(ruta); os.IsNotExist(err) {
		nuevoArchivo = true
	}
	archivo, err := os.OpenFile(ruta, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}
	escritor := csv.NewWriter(archivo)
	if nuevoArchivo {
		_ = escritor.Write([]string{"modo", "ejecucion", "inicio_ms", "fin_ms", "total_ms", "decision_ms", "rama", "rama_ms", "otra_cancelada", "detalle", "traza", "umbral", "dificultad", "maximo_primo"})
		escritor.Flush()
	}
	return archivo, escritor, nil
}

func escribirFila(escritor *csv.Writer, m Metricas, ejecucion int) error {
	return escritor.Write([]string{
		m.Modo, fmt.Sprintf("%d", ejecucion),
		fmt.Sprintf("%d", m.InicioUnixMs),
		fmt.Sprintf("%d", m.FinUnixMs),
		fmt.Sprintf("%d", m.TotalMs),
		fmt.Sprintf("%d", m.DecisionMs),
		m.Rama,
		fmt.Sprintf("%d", m.DuracionRamaMs),
		fmt.Sprintf("%t", m.OtraCancelada),
		m.HashOUltimoPrimo,
		fmt.Sprintf("%d", m.Traza),
		fmt.Sprintf("%d", m.Umbral),
		fmt.Sprintf("%d", m.Dificultad),
		fmt.Sprintf("%d", m.MaximoPrimo),
	})
}

// -------------------------------
// Ejecución Especulativa y Secuencial
// -------------------------------

func ejecutarEspeculativo(n, umbral, dificultad, maximoPrimo int, datosBloque string) (Metricas, error) {
	inicioTotal := time.Now()

	ctxA, cancelarA := context.WithCancel(context.Background())
	ctxB, cancelarB := context.WithCancel(context.Background())
	defer func() { cancelarA(); cancelarB() }()

	resultadoACh := make(chan Resultado, 1)
	resultadoBCh := make(chan Resultado, 1)

	go func() {
		resultadoACh <- ejecutarTarea(ctxA, "A", func(ctx context.Context) (string, error) {
			hash, intento, err := SimularProofOfWork(ctx, datosBloque, dificultad)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("%s (intento=%d)", hash, intento), nil
		})
	}()
	go func() {
		resultadoBCh <- ejecutarTarea(ctxB, "B", func(ctx context.Context) (string, error) {
			primos, err := EncontrarPrimos(ctx, maximoPrimo)
			if err != nil {
				return "", err
			}
			ultimo := -1
			if len(primos) > 0 {
				ultimo = primos[len(primos)-1]
			}
			return fmt.Sprintf("ultimoPrimo=%d (cantidad=%d)", ultimo, len(primos)), nil
		})
	}()

	inicioDecision := time.Now()
	traza := CalcularTrazaProductoMatrices(n)
	duracionDecision := time.Since(inicioDecision)

	ganadora := "B"
	if traza >= umbral {
		ganadora = "A"
	}

	var resultadoA, resultadoB Resultado
	var otraCancelada bool
	if ganadora == "A" {
		resultadoA = <-resultadoACh
		cancelarB()
		otraCancelada = true
		resultadoB = <-resultadoBCh // drenaje
	} else {
		resultadoB = <-resultadoBCh
		cancelarA()
		otraCancelada = true
		resultadoA = <-resultadoACh // drenaje
	}

	total := time.Since(inicioTotal)
	metrica := Metricas{
		Modo: "especulativo", InicioUnixMs: inicioTotal.UnixMilli(), FinUnixMs: inicioTotal.Add(total).UnixMilli(),
		TotalMs: total.Milliseconds(), DecisionMs: duracionDecision.Milliseconds(),
		Rama: ganadora, Traza: traza, Umbral: umbral, Dificultad: dificultad, MaximoPrimo: maximoPrimo,
		OtraCancelada: otraCancelada,
	}
	if ganadora == "A" {
		if resultadoA.Error != nil {
			return metrica, resultadoA.Error
		}
		metrica.DuracionRamaMs, metrica.HashOUltimoPrimo = resultadoA.Duracion.Milliseconds(), resultadoA.Salida
	} else {
		if resultadoB.Error != nil {
			return metrica, resultadoB.Error
		}
		metrica.DuracionRamaMs, metrica.HashOUltimoPrimo = resultadoB.Duracion.Milliseconds(), resultadoB.Salida
	}
	return metrica, nil
}

func ejecutarSecuencial(n, umbral, dificultad, maximoPrimo int, datosBloque string) (Metricas, error) {
	inicio := time.Now()
	inicioDecision := time.Now()
	traza := CalcularTrazaProductoMatrices(n)
	duracionDecision := time.Since(inicioDecision)

	ganadora := "B"
	if traza >= umbral {
		ganadora = "A"
	}
	var duracionRama time.Duration
	var detalle string
	var err error

	if ganadora == "A" {
		inicioRama := time.Now()
		hash, intento, e := SimularProofOfWork(context.Background(), datosBloque, dificultad)
		duracionRama, err = time.Since(inicioRama), e
		detalle = fmt.Sprintf("%s (intento=%d)", hash, intento)
	} else {
		inicioRama := time.Now()
		primos, e := EncontrarPrimos(context.Background(), maximoPrimo)
		duracionRama, err = time.Since(inicioRama), e
		ultimo := -1
		if len(primos) > 0 {
			ultimo = primos[len(primos)-1]
		}
		detalle = fmt.Sprintf("ultimoPrimo=%d (cantidad=%d)", ultimo, len(primos))
	}
	if err != nil {
		return Metricas{}, err
	}

	total := time.Since(inicio)
	return Metricas{
		Modo: "secuencial", InicioUnixMs: inicio.UnixMilli(), FinUnixMs: inicio.Add(total).UnixMilli(),
		TotalMs: total.Milliseconds(), DecisionMs: duracionDecision.Milliseconds(),
		Rama: ganadora, DuracionRamaMs: duracionRama.Milliseconds(),
		HashOUltimoPrimo: detalle, Traza: traza, Umbral: umbral, Dificultad: dificultad, MaximoPrimo: maximoPrimo,
	}, nil
}

// -------------------------------
// main
// -------------------------------

func main() {
	var (
		n            = flag.Int("n", 300, "Dimensión de las matrices (decisión)")
		umbral       = flag.Int("umbral", 0, "Si traza >= umbral -> rama A, si no -> B")
		archivoOut   = flag.String("nombre_archivo", "metricas.csv", "Archivo CSV de métricas")
		modo         = flag.String("mode", "both", "espec | secuencial | ambos")
		repeticiones = flag.Int("runs", 1, "Repeticiones para promediar")
		dificultad   = flag.Int("dificultad", 5, "Dificultad Prueba de Trabajo (A)")
		maximoPrimo  = flag.Int("max_primo", 500000, "Límite para primos (B)")
		datosBloque  = flag.String("block_data", "Bloque de ejemplo", "Datos para PoW")
		semilla      = flag.Int64("seed", time.Now().UnixNano(), "Semilla para reproducibilidad")
	)
	flag.Parse()
	rand.Seed(*semilla)

	// Señales para salida limpia
	senales := make(chan os.Signal, 1)
	signal.Notify(senales, syscall.SIGINT, syscall.SIGTERM)

	// CSV
	archivo, escritor, err := asegurarCSV(*archivoOut)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error CSV:", err)
		os.Exit(1)
	}
	defer archivo.Close()
	defer escritor.Flush()

	titulo := "Ejecución"
	switch *modo {
	case "spec", "espec":
		titulo = "Ejecución Especulativa"
	case "seq", "secuencial":
		titulo = "Ejecución Secuencial"
	case "both", "ambos":
		titulo = "Ejecución Especulativa y Secuencial"
	}

	fmt.Printf("C2 - %s (Go)\nParámetros: n=%d umbral=%d repeticiones=%d dificultad=%d max_primo=%d\n",
		titulo, *n, *umbral, *repeticiones, *dificultad, *maximoPrimo)
	fmt.Println("----------------------------------------------------------------------")

	var corridasEspec, corridasSec []Metricas
	var mutex sync.Mutex
	var wg sync.WaitGroup

	for i := 1; i <= *repeticiones; i++ {
		select {
		case <-senales:
			fmt.Println("\nSeñal recibida, saliendo…")
			return
		default:
		}
		if *modo == "spec" || *modo == "both" {
			wg.Add(1)
			go func(ejecucion int) {
				defer wg.Done()
				m, err := ejecutarEspeculativo(*n, *umbral, *dificultad, *maximoPrimo, *datosBloque)
				if err != nil {
					fmt.Fprintln(os.Stderr, "error especulativo:", err)
					return
				}
				mutex.Lock()
				corridasEspec = append(corridasEspec, m)
				_ = escribirFila(escritor, m, ejecucion)
				escritor.Flush()
				mutex.Unlock()
				fmt.Printf("[espec] ejec=%d total=%dms decision=%dms rama=%s rama_ms=%d cancelada=%t traza=%d\n",
					ejecucion, m.TotalMs, m.DecisionMs, m.Rama, m.DuracionRamaMs, m.OtraCancelada, m.Traza)
			}(i)
		}
		if *modo == "seq" || *modo == "both" {
			wg.Add(1)
			go func(ejecucion int) {
				defer wg.Done()
				m, err := ejecutarSecuencial(*n, *umbral, *dificultad, *maximoPrimo, *datosBloque)
				if err != nil {
					fmt.Fprintln(os.Stderr, "error secuencial:", err)
					return
				}
				mutex.Lock()
				corridasSec = append(corridasSec, m)
				_ = escribirFila(escritor, m, ejecucion)
				escritor.Flush()
				mutex.Unlock()
				fmt.Printf("[secu] ejec=%d total=%dms decision=%dms rama=%s rama_ms=%d traza=%d\n",
					ejecucion, m.TotalMs, m.DecisionMs, m.Rama, m.DuracionRamaMs, m.Traza)
			}(i)
		}
	}
	wg.Wait()

	// Resumen simple
	promedio := func(xs []Metricas, f func(Metricas) int64) float64 {
		if len(xs) == 0 {
			return 0
		}
		var suma int64
		for _, m := range xs {
			suma += f(m)
		}
		return float64(suma) / float64(len(xs))
	}
	if len(corridasEspec) > 0 {
		fmt.Printf("\nResumen especulativo: ejec=%d prom_total=%.2fms prom_decision=%.2fms prom_rama=%.2fms\n",
			len(corridasEspec),
			promedio(corridasEspec, func(m Metricas) int64 { return m.TotalMs }),
			promedio(corridasEspec, func(m Metricas) int64 { return m.DecisionMs }),
			promedio(corridasEspec, func(m Metricas) int64 { return m.DuracionRamaMs }))
	}
	if len(corridasSec) > 0 {
		fmt.Printf("Resumen secuencial : ejec=%d prom_total=%.2fms prom_decision=%.2fms prom_rama=%.2fms\n",
			len(corridasSec),
			promedio(corridasSec, func(m Metricas) int64 { return m.TotalMs }),
			promedio(corridasSec, func(m Metricas) int64 { return m.DecisionMs }),
			promedio(corridasSec, func(m Metricas) int64 { return m.DuracionRamaMs }))
	}
	if len(corridasEspec) > 0 && len(corridasSec) > 0 {
		promEspec := promedio(corridasEspec, func(m Metricas) int64 { return m.TotalMs })
		promSec := promedio(corridasSec, func(m Metricas) int64 { return m.TotalMs })
		if promEspec > 0 {
			fmt.Printf("\nSpeedup = T_sec / T_espec = %.3f\n", promSec/promEspec)
		}
	}
	fmt.Println("\nMétricas escritas en:", *archivoOut)
}
