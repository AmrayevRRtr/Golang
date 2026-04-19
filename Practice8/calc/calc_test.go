package calc

import (
	"testing"
)

//func TestAdd(t *testing.T) {
//	got := Add(2, 3)
//	want := 6
//	if got != want {
//		t.Errorf("Add(2, 3) = %d; want %d", got, want)
//	}
//}

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

func TestDivide(t *testing.T) {
	result, err := Divide(10, 2)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if result != 5 {
		t.Errorf("Divide(10, 2) = %d; want 5", result)
	}

	_, err = Divide(10, 0)

	if err == nil {
		t.Errorf("expected error for division by zero, got nil")
	}
}

func TestSubtract_TableDriven(t *testing.T) {
	tests := []struct {
		name string
		a, b int
		want int
	}{
		{"both positive", 5, 3, 2},
		{"positive minus zero", 5, 0, 5},
		{"negative minus positive", -5, 3, -8},
		{"both negative", -5, -3, -2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Subtract(tt.a, tt.b)

			if got != tt.want {
				t.Errorf("Subtract(%d, %d) = %d; want %d",
					tt.a, tt.b, got, tt.want)
			}
		})
	}
}
