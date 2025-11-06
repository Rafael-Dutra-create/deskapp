package utils

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"strconv"
	"testing"
)

func almostEqual(a, b, epsilon float64) bool {
	return math.Abs(a-b) < epsilon
}

func readCSV(fileName string) []float64 {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Erro ao abrir o arquivo:", err)
		return nil
	}
	defer file.Close()

	// Criar um novo leitor CSV
	reader := csv.NewReader(file)

	// Ler todo o conteúdo do CSV
	records, err := reader.ReadAll()
	r := make([]float64, len(records))
	if err != nil {
		fmt.Println("Erro ao ler o arquivo CSV:", err)
		return nil
	}

	for i, value := range records {
		num, _ := strconv.ParseFloat(value[0], 64)
		r[i] = num
	}
	return r
}

func TestSmooth4253H(t *testing.T) {
	input := readCSV("input.csv")
	target := readCSV("output.csv")
	dataSmooth := Smooth4253H(input)
	for i, v := range target {
		if !almostEqual(v, dataSmooth[i], 0.0001) {
			t.Errorf("\nEsperado: %f - Recebido: %f\n", v, dataSmooth[i])
		}
	}
}

func BenchmarkSmooth4253H(b *testing.B) {
	for i := 0; i < b.N; i++ {
		input := readCSV("input.csv")
		Smooth4253H(input)
	}
}
