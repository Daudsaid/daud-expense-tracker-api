package models

import "testing"

func TestTotal(t *testing.T) {
	expenses := []Expense{
		{Amount: 10},
		{Amount: 5.5},
	}
	got := Total(expenses)
	want := 15.5
	if got != want {
		t.Fatalf("Total() = %v, want %v", got, want)
	}
}

func TestTotalsByCategory(t *testing.T) {
	expenses := []Expense{
		{Amount: 10, Category: "Food"},
		{Amount: 5, Category: "Food"},
		{Amount: 7, Category: "Transport"},
	}

	got := TotalsByCategory(expenses)

	if got["Food"] != 15 {
		t.Fatalf("Food total = %v, want %v", got["Food"], 15)
	}
	if got["Transport"] != 7 {
		t.Fatalf("Transport total = %v, want %v", got["Transport"], 7)
	}
}
