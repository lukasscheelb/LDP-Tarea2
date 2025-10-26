import pandas as pd
import matplotlib.pyplot as plt

# ==========================
# CARGA DE DATOS
# ==========================

# Archivos CSV (deben estar en la misma carpeta que este script)
seq = pd.read_csv("metricas_seq.csv")
spec = pd.read_csv("metricas_spec.csv")

# Mostrar las columnas detectadas (verificación)
print("Columnas detectadas:", list(seq.columns))

# ==========================
# CÁLCULOS DE PROMEDIOS Y SPEEDUP
# ==========================

prom_seq = seq["total_ms"].mean()
prom_spec = spec["total_ms"].mean()
speedup = prom_seq / prom_spec

print(f"Promedio secuencial : {prom_seq:.2f} ms")
print(f"Promedio especulativo: {prom_spec:.2f} ms")
print(f"Speedup: {speedup:.3f}x")

# ==========================
# GRÁFICO 1: Barras comparando promedios
# ==========================

plt.figure(figsize=(7, 5))
plt.bar(["Secuencial", "Especulativo"], [prom_seq, prom_spec],
        color=["#FF9999", "#90EE90"], edgecolor="black")
plt.title("Comparación de Tiempo Promedio Total", fontsize=13)
plt.ylabel("Tiempo (ms)")
plt.text(0, prom_seq + 50, f"{prom_seq:.2f}", ha="center")
plt.text(1, prom_spec + 50, f"{prom_spec:.2f}", ha="center")
plt.grid(axis="y", linestyle="--", alpha=0.6)
plt.tight_layout()
plt.savefig("grafico_promedios.png", dpi=300)
plt.close()

# ==========================
# GRÁFICO 2: Distribución de tiempos (Boxplot)
# ==========================

plt.figure(figsize=(7, 5))
plt.boxplot([seq["total_ms"], spec["total_ms"]],
            labels=["Secuencial", "Especulativo"],
            patch_artist=True,
            boxprops=dict(facecolor="#AED6F1", color="black"),
            medianprops=dict(color="red"))
plt.title("Distribución de Tiempos Totales (30 corridas)", fontsize=13)
plt.ylabel("Tiempo (ms)")
plt.grid(axis="y", linestyle="--", alpha=0.6)
plt.tight_layout()
plt.savefig("grafico_boxplot.png", dpi=300)
plt.close()

# ==========================
# GRÁFICO 3: Evolución temporal por corrida
# ==========================

plt.figure(figsize=(8, 5))
plt.plot(seq["ejecucion"], seq["total_ms"], marker="o", label="Secuencial", color="#E74C3C")
plt.plot(spec["ejecucion"], spec["total_ms"], marker="s", label="Especulativo", color="#27AE60")
plt.title("Evolución del Tiempo Total por Corrida", fontsize=13)
plt.xlabel("Número de corrida")
plt.ylabel("Tiempo total (ms)")
plt.legend()
plt.grid(True, linestyle="--", alpha=0.6)
plt.tight_layout()
plt.savefig("grafico_evolucion.png", dpi=300)
plt.close()

# ==========================
# GRÁFICO 4: Visual del Speedup
# ==========================

plt.figure(figsize=(6, 4))
plt.bar(["Speedup"], [speedup], color="#F7DC6F", edgecolor="black")
plt.text(0, speedup / 2, f"{speedup:.3f}x", ha="center", va="center", fontsize=16, fontweight="bold")
plt.ylim(0, max(1.0, speedup + 0.1))
plt.title("Speedup Global (T_sec / T_spec)", fontsize=13)
plt.ylabel("Factor de Mejora")
plt.grid(axis="y", linestyle="--", alpha=0.6)
plt.tight_layout()
plt.savefig("grafico_speedup.png", dpi=300)
plt.close()

# ==========================
# RESUMEN EN CONSOLA
# ==========================

print("\nGráficos generados correctamente:")
print("- grafico_promedios.png   → Comparación de promedios")
print("- grafico_boxplot.png      → Distribución de tiempos")
print("- grafico_evolucion.png    → Evolución por corrida")
print("- grafico_speedup.png      → Visualización del Speedup")
print(f"\nSpeedup calculado: {speedup:.3f}x ({(speedup - 1)*100:.2f}% de mejora)")
