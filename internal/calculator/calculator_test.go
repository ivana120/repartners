package calculator

import (
	"testing"
)

func TestCalculator_Calculate(t *testing.T) {
	tests := []struct {
		name           string
		packSizes      []int
		orderQty       int
		wantTotalItems int
		wantTotalPacks int
		wantPacks      map[int]int
		wantErr        error
	}{
		{
			name:           "order 1 item with standard packs",
			packSizes:      []int{250, 500, 1000, 2000, 5000},
			orderQty:       1,
			wantTotalItems: 250,
			wantTotalPacks: 1,
			wantPacks:      map[int]int{250: 1},
		},
		{
			name:           "order exactly 250",
			packSizes:      []int{250, 500, 1000, 2000, 5000},
			orderQty:       250,
			wantTotalItems: 250,
			wantTotalPacks: 1,
			wantPacks:      map[int]int{250: 1},
		},
		{
			name:           "order 251",
			packSizes:      []int{250, 500, 1000, 2000, 5000},
			orderQty:       251,
			wantTotalItems: 500,
			wantTotalPacks: 1,
			wantPacks:      map[int]int{500: 1},
		},
		{
			name:           "order 501",
			packSizes:      []int{250, 500, 1000, 2000, 5000},
			orderQty:       501,
			wantTotalItems: 750,
			wantTotalPacks: 2,
			wantPacks:      map[int]int{500: 1, 250: 1},
		},
		{
			name:           "order 12001",
			packSizes:      []int{250, 500, 1000, 2000, 5000},
			orderQty:       12001,
			wantTotalItems: 12250,
			wantTotalPacks: 4,
			wantPacks:      map[int]int{5000: 2, 2000: 1, 250: 1},
		},
		{
			name:           "exact match with multiple packs",
			packSizes:      []int{250, 500, 1000, 2000, 5000},
			orderQty:       1000,
			wantTotalItems: 1000,
			wantTotalPacks: 1,
			wantPacks:      map[int]int{1000: 1},
		},
		{
			name:           "no pack sizes",
			packSizes:      []int{},
			orderQty:       100,
			wantErr:        ErrNoPackSizes,
		},
		{
			name:           "invalid order quantity",
			packSizes:      []int{250, 500},
			orderQty:       0,
			wantErr:        ErrInvalidOrder,
		},
		{
			name:           "negative order quantity",
			packSizes:      []int{250, 500},
			orderQty:       -5,
			wantErr:        ErrInvalidOrder,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calc := New(tt.packSizes)
			result, err := calc.Calculate(tt.orderQty)

			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("expected error %v, got %v", tt.wantErr, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if result.TotalItems != tt.wantTotalItems {
				t.Errorf("TotalItems: want %d, got %d", tt.wantTotalItems, result.TotalItems)
			}

			if result.TotalPacks != tt.wantTotalPacks {
				t.Errorf("TotalPacks: want %d, got %d", tt.wantTotalPacks, result.TotalPacks)
			}

			if len(result.Packs) != len(tt.wantPacks) {
				t.Errorf("Packs count: want %d, got %d", len(tt.wantPacks), len(result.Packs))
			}

			for size, qty := range tt.wantPacks {
				if result.Packs[size] != qty {
					t.Errorf("Pack[%d]: want %d, got %d", size, qty, result.Packs[size])
				}
			}
		})
	}
}

func TestEdgeCase(t *testing.T) {
	calc := New([]int{23, 31, 53})
	result, err := calc.Calculate(500000)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedPacks := map[int]int{23: 2, 31: 7, 53: 9429}
	expectedTotal := 500000
	expectedPackCount := 2 + 7 + 9429

	if result.TotalItems != expectedTotal {
		t.Errorf("TotalItems: want %d, got %d", expectedTotal, result.TotalItems)
	}

	if result.TotalPacks != expectedPackCount {
		t.Errorf("TotalPacks: want %d, got %d", expectedPackCount, result.TotalPacks)
	}

	for size, qty := range expectedPacks {
		if result.Packs[size] != qty {
			t.Errorf("Pack[%d]: want %d, got %d", size, qty, result.Packs[size])
		}
	}

	calculatedTotal := 0
	for size, qty := range result.Packs {
		calculatedTotal += size * qty
	}
	if calculatedTotal < 500000 {
		t.Errorf("Total items %d is less than order quantity 500000", calculatedTotal)
	}
}

func TestSetPackSizes(t *testing.T) {
	calc := New([]int{100, 200})

	result, _ := calc.Calculate(150)
	if result.TotalItems != 200 {
		t.Errorf("expected 200, got %d", result.TotalItems)
	}

	calc.SetPackSizes([]int{50, 100})

	result, _ = calc.Calculate(150)
	if result.TotalItems != 150 {
		t.Errorf("expected 150, got %d", result.TotalItems)
	}
}

func TestPackSizes(t *testing.T) {
	sizes := []int{100, 200, 300}
	calc := New(sizes)

	got := calc.PackSizes()
	if len(got) != len(sizes) {
		t.Errorf("expected %d pack sizes, got %d", len(sizes), len(got))
	}
}

func BenchmarkCalculate(b *testing.B) {
	calc := New([]int{23, 31, 53})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = calc.Calculate(500000)
	}
}

func BenchmarkCalculateStandardPacks(b *testing.B) {
	calc := New([]int{250, 500, 1000, 2000, 5000})
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = calc.Calculate(12001)
	}
}
