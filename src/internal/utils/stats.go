package utils

import "sort"

// Map percorre o array e retorna de acordo com a funcao
func Map[S ~[]E, E any, R any](array S, fn func(E) R) []R {
	var result []R
	for _, elem := range array {
		result = append(result, fn(elem))
	}
	return result
}

// Mean calcula a média
func Mean[S []E, E int | int16 | int8 | float32 | float64](results S) float64 {
	if len(results) == 0 {
		return 0
	}
	var sumOne int = 0
	total := len(results)
	for _, v := range results {
		sumOne += int(v)
	}
	return float64(sumOne) / float64(total)
}

// Median calcula a Mediana de um slice
func Median(slice []float64) float64 {
	sort.Float64s(slice)
	n := len(slice)
	if n%2 == 0 {
		return (slice[n/2-1] + slice[n/2]) / 2
	}
	return slice[n/2]
}

// Função h para suavizar a série y
func h(y []float64) []float64 {
	N := len(y)
	z := make([]float64, N+1)

	// Suavização inicial com Mediana
	z[0] = y[0]
	z[1] = Median([]float64{y[0], y[1]})
	z[2] = Median([]float64{y[1], y[2]})

	for i := 3; i < N; i++ {
		z[i] = Median([]float64{y[i-3], y[i-2], y[i-1], y[i]})
	}

	z[N-2] = Median([]float64{y[N-3], y[N-2]})
	z[N-1] = Median([]float64{y[N-2], y[N-1]})
	z[N] = y[N-1]

	// Segunda fase de suavização
	z1 := make([]float64, N)
	for i := 0; i < N; i++ {
		z1[i] = (z[i] + z[i+1]) / 2
	}

	// Terceira fase de suavização
	z2 := make([]float64, N)
	z2[0] = z1[0]
	z2[1] = Median([]float64{z1[0], z1[1], z1[2]})

	for i := 2; i < N-2; i++ {
		z2[i] = Median([]float64{z1[i-2], z1[i-1], z1[i], z1[i+1], z1[i+2]})
	}
	z2[N-2] = Median([]float64{z1[N-3], z1[N-2], z1[N-1]})
	z2[N-1] = z1[N-1]

	// Quarta fase de suavização
	z3 := make([]float64, N)
	z3[0] = z2[0]
	for i := 1; i < N-1; i++ {
		z3[i] = Median([]float64{z2[i-1], z2[i], z2[i+1]})
	}
	z3[N-1] = z2[N-1]

	// Última fase de suavização com média ponderada
	z4 := make([]float64, N)
	z4[0] = z3[0]
	for i := 1; i < N-1; i++ {
		z4[i] = (z3[i-1] + z3[i] + z3[i+1]) / 4
	}
	z4[N-1] = z3[N-1]

	// Ajuste dos extremos
	z4[0] = Median([]float64{z4[0], z4[1], (3*z4[1] - 2*z4[2])})
	z4[N-1] = Median([]float64{z4[N-1], z4[N-2], (3*z4[N-3] - 2*z4[N-2])})

	return z4
}

// Função principal para suavização e remoção de ruído
func Smooth4253H(y []float64) []float64 {
	sm := h(y)

	// Cálculo dos resíduos (y - sm)
	rf := make([]float64, len(y))
	for i := 0; i < len(y); i++ {
		rf[i] = y[i] - sm[i]
	}

	// Suavizar os resíduos
	smRf := h(rf)

	// Suavização final
	smooth := make([]float64, len(y))
	for i := 0; i < len(y); i++ {
		smooth[i] = sm[i] + smRf[i]
	}

	return smooth
}
