package temperature

import "testing"

func TestCelsiusToFahrenheit(t *testing.T) {
	tests := []struct {
		name string
		c    float64
		want float64
	}{
		{"zero", 0, 32},
		{"positivo", 28.5, 83.3},
		{"negativo", -10, 14},
		{"ruido de ponto flutuante", 36.6, 97.88},
		{"cem", 100, 212},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CelsiusToFahrenheit(tt.c); got != tt.want {
				t.Errorf("CelsiusToFahrenheit(%v) = %v, want %v", tt.c, got, tt.want)
			}
		})
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	tests := []struct {
		name string
		c    float64
		want float64
	}{
		{"zero", 0, 273},
		{"positivo", 28.5, 301.5},
		{"negativo", -10, 263},
		{"fracionario", 0.1, 273.1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CelsiusToKelvin(tt.c); got != tt.want {
				t.Errorf("CelsiusToKelvin(%v) = %v, want %v", tt.c, got, tt.want)
			}
		})
	}
}
