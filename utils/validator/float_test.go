package validator_test

import (
	"github.com/meli-fresh-products-api-backend-t1/utils/validator"
	"testing"
)

func TestFloatIsPositive(t *testing.T) {
	tests := []struct {
		input float64
		want  bool
	}{
		{1.0, true},
		{0.0, false},
		{-1.0, false},
	}
	for _, tt := range tests {
		t.Run("TestFloatIsPositive", func(t *testing.T) {
			if got := validator.FloatIsPositive(tt.input); got != tt.want {
				t.Errorf("FloatIsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatIsNegative(t *testing.T) {
	tests := []struct {
		input float64
		want  bool
	}{
		{-1.0, true},
		{0.0, false},
		{1.0, false},
	}
	for _, tt := range tests {
		t.Run("TestFloatIsNegative", func(t *testing.T) {
			if got := validator.FloatIsNegative(tt.input); got != tt.want {
				t.Errorf("FloatIsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatIsZero(t *testing.T) {
	tests := []struct {
		input float64
		want  bool
	}{
		{0.0, true},
		{1.0, false},
		{-1.0, false},
	}
	for _, tt := range tests {
		t.Run("TestFloatIsZero", func(t *testing.T) {
			if got := validator.FloatIsZero(tt.input); got != tt.want {
				t.Errorf("FloatIsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatIsGreaterThan(t *testing.T) {
	tests := []struct {
		i    float64
		j    float64
		want bool
	}{
		{2.0, 1.0, true},
		{1.0, 1.0, false},
		{1.0, 2.0, false},
	}
	for _, tt := range tests {
		t.Run("TestFloatIsGreaterThan", func(t *testing.T) {
			if got := validator.FloatIsGreaterThan(tt.i, tt.j); got != tt.want {
				t.Errorf("FloatIsGreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatIsLessThan(t *testing.T) {
	tests := []struct {
		i    float64
		j    float64
		want bool
	}{
		{1.0, 2.0, true},
		{1.0, 1.0, false},
		{2.0, 1.0, false},
	}
	for _, tt := range tests {
		t.Run("TestFloatIsLessThan", func(t *testing.T) {
			if got := validator.FloatIsLessThan(tt.i, tt.j); got != tt.want {
				t.Errorf("FloatIsLessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFloatBetween(t *testing.T) {
	tests := []struct {
		i    float64
		minV float64
		maxV float64
		want bool
	}{
		{2.0, 1.0, 3.0, true},
		{1.0, 1.0, 3.0, true},
		{3.0, 1.0, 3.0, true},
		{0.0, 1.0, 3.0, false},
		{4.0, 1.0, 3.0, false},
	}
	for _, tt := range tests {
		t.Run("TestFloatBetween", func(t *testing.T) {
			if got := validator.FloatBetween(tt.i, tt.minV, tt.maxV); got != tt.want {
				t.Errorf("FloatBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
