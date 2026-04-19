package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ─── From tutorial: basic Add test ───────────────────────────────────────────

func TestAdd(t *testing.T) {
	got := Add(2, 3)
	want := 5
	if got != want {
		t.Errorf("Add(2, 3) = %d; want %d", got, want)
	}
}

// Table-driven Add test (from tutorial Step 1)
func TestAddTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 2, 3, 5},
		{"positive + zero", 5, 0, 5},
		{"negative + positive", -1, 4, 3},
		{"both negative", -2, -3, -5},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Add(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

// ─── Task 1: Divide ──────────────────────────────────────────────────────────
// Tests both success cases and division-by-zero error.

func TestDivide(t *testing.T) {
	tests := []struct {
		name    string
		a, b    int
		want    int
		wantErr bool
	}{
		{"10 divided by 2", 10, 2, 5, false},
		{"9 divided by 3", 9, 3, 3, false},
		{"negative divided by positive", -10, 2, -5, false},
		{"zero divided by number", 0, 5, 0, false},
		{"error: division by zero", 5, 0, 0, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Divide(tt.a, tt.b)
			if tt.wantErr {
				require.Error(t, err)
				assert.Equal(t, "division by zero", err.Error())
				assert.Equal(t, 0, got)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

// ─── Task 1: Subtract (table-driven) ─────────────────────────────────────────
// Covers: both positive, positive minus zero, negative minus positive, both negative.

func TestSubtractTableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 5, 3, 2},
		{"positive minus zero", 5, 0, 5},
		{"negative minus positive", -3, 4, -7},
		{"both negative", -2, -3, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)
			assert.Equal(t, tt.want, got, "Subtract(%d, %d)", tt.a, tt.b)
		})
	}
}
