// Package temperature converte temperaturas a partir de graus Celsius.
package temperature

import "math"

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

// CelsiusToFahrenheit aplica F = C × 1.8 + 32, arredondado a 2 casas.
func CelsiusToFahrenheit(c float64) float64 {
	return round2(c*1.8 + 32)
}

// CelsiusToKelvin aplica K = C + 273, arredondado a 2 casas.
func CelsiusToKelvin(c float64) float64 {
	return round2(c + 273)
}
