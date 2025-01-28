package validator_test

import (
	"github.com/meli-fresh-products-api-backend-t1/utils/validator"
	"testing"
)

func TestIntIsPositive(t *testing.T) {
	tests := []struct {
		input int
		want  bool
	}{
		{1, true},
		{0, false},
		{-1, false},
	}
	for _, tt := range tests {
		t.Run("TestIntIsPositive", func(t *testing.T) {
			if got := validator.IntIsPositive(tt.input); got != tt.want {
				t.Errorf("IntIsPositive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntIsNegative(t *testing.T) {
	tests := []struct {
		input int
		want  bool
	}{
		{-1, true},
		{0, false},
		{1, false},
	}
	for _, tt := range tests {
		t.Run("TestIntIsNegative", func(t *testing.T) {
			if got := validator.IntIsNegative(tt.input); got != tt.want {
				t.Errorf("IntIsNegative() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntIsZero(t *testing.T) {
	tests := []struct {
		input int
		want  bool
	}{
		{0, true},
		{1, false},
		{-1, false},
	}
	for _, tt := range tests {
		t.Run("TestIntIsZero", func(t *testing.T) {
			if got := validator.IntIsZero(tt.input); got != tt.want {
				t.Errorf("IntIsZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntIsGreaterThan(t *testing.T) {
	tests := []struct {
		i    int
		j    int
		want bool
	}{
		{2, 1, true},
		{1, 1, false},
		{1, 2, false},
	}
	for _, tt := range tests {
		t.Run("TestIntIsGreaterThan", func(t *testing.T) {
			if got := validator.IntIsGreaterThan(tt.i, tt.j); got != tt.want {
				t.Errorf("IntIsGreaterThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntIsLessThan(t *testing.T) {
	tests := []struct {
		i    int
		j    int
		want bool
	}{
		{1, 2, true},
		{1, 1, false},
		{2, 1, false},
	}
	for _, tt := range tests {
		t.Run("TestIntIsLessThan", func(t *testing.T) {
			if got := validator.IntIsLessThan(tt.i, tt.j); got != tt.want {
				t.Errorf("IntIsLessThan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIntBetween(t *testing.T) {
	tests := []struct {
		i    int
		minV int
		maxV int
		want bool
	}{
		{2, 1, 3, true},
		{1, 1, 3, true},
		{3, 1, 3, true},
		{0, 1, 3, false},
		{4, 1, 3, false},
	}
	for _, tt := range tests {
		t.Run("TestIntBetween", func(t *testing.T) {
			if got := validator.IntBetween(tt.i, tt.minV, tt.maxV); got != tt.want {
				t.Errorf("IntBetween() = %v, want %v", got, tt.want)
			}
		})
	}
}
