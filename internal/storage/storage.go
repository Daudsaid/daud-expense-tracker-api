package storage

import (
	"github.com/Daudsaid/daud-expense-tracker-api/internal/models"
)

// DataFile is the default JSON file used for persistence.
const DataFile = "expenses.json"

// Store defines behaviour for persisting expenses.
type Store interface {
	ListExpenses() ([]models.Expense, error)
	CreateExpense(e models.Expense) (models.Expense, error)
	GetExpense(id int64) (models.Expense, error)
	DeleteExpense(id int64) error
}

// NewStore returns a default store implementation.
func NewStore(path string) (Store, error) {
	return NewMemoryStore(path)
}
