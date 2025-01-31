package validator_test

import (
	"github.com/meli-fresh-products-api-backend-t1/utils/validator"
	"testing"
)

func TestString(t *testing.T) {
	tests := []struct {
		str  string
		minV int
		maxV int
		want bool
	}{
		{"teste", 3, 10, true},
		{"teste", 6, 10, false},
		{"teste", 3, 4, false},
		{"", 0, 10, false},
		{" ", 0, 10, false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := validator.String(tt.str, tt.minV, tt.maxV); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEmptyString(t *testing.T) {
	tests := []struct {
		str  string
		want bool
	}{
		{"", true},
		{"teste", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := validator.EmptyString(tt.str); got != tt.want {
				t.Errorf("EmptyString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlankString(t *testing.T) {
	tests := []struct {
		str  string
		want bool
	}{
		{"", true},
		{" ", true},
		{"teste", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := validator.BlankString(tt.str); got != tt.want {
				t.Errorf("BlankString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsCep(t *testing.T) {
	tests := []struct {
		cep  string
		want bool
	}{
		{"12345-678", true},
		{"12345678", false},
		{"1239-678", false},
		{"12345-6789", false},
		{"abcde-fgh", false},
	}
	for _, tt := range tests {
		t.Run("TestIsCep", func(t *testing.T) {
			if got := validator.IsCep(tt.cep); got != tt.want {
				t.Errorf("IsCep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEmail(t *testing.T) {
	tests := []struct {
		e    string
		want bool
	}{
		{"teste@example.com", true},
		{"teste.teste@example.com", true},
		{"teste+teste@example.com", false},
		{"teste@example", false},
		{"teste@.com", false},
		{"@example.com", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := validator.IsEmail(tt.e); got != tt.want {
				t.Errorf("IsEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsTelephone(t *testing.T) {
	tests := []struct {
		t    string
		want bool
	}{
		{"12345-6789", true},
		{"(11) 1234-5678", true},
		{"11 12345678", false},
		{"1234-56789", false},
		{"abcde-fghi", false},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := validator.IsTelephone(tt.t); got != tt.want {
				t.Errorf("IsTelephone() = %v, want %v", got, tt.want)
			}
		})
	}
}
