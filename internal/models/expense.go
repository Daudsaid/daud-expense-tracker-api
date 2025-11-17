package models

// Expense represents a single expense entry.
type Expense struct {
	ID       int64   `json:"id"`
	Amount   float64 `json:"amount"`
	Category string  `json:"category"`
	Note     string  `json:"note"`
	Date     string  `json:"date"` // "YYYY-MM-DD"
}

// Total returns the sum of all expense amounts.
func Total(expenses []Expense) float64 {
	var total float64
	for _, e := range expenses {
		total += e.Amount
	}
	return total
}

// TotalsByCategory returns a map of category -> total amount.
func TotalsByCategory(expenses []Expense) map[string]float64 {
	totals := make(map[string]float64)
	for _, e := range expenses {
		totals[e.Category] += e.Amount
	}
	return totals
}
