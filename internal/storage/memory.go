package storage

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Daudsaid/daud-expense-tracker-api/internal/models"
)

// MemoryStore is an in-memory implementation of Store with JSON file persistence.
type MemoryStore struct {
	mu       sync.Mutex
	expenses []models.Expense
	nextID   int64
	path     string
}

// NewMemoryStore initialises a MemoryStore and loads data from JSON if present.
func NewMemoryStore(path string) (*MemoryStore, error) {
	ms := &MemoryStore{
		expenses: make([]models.Expense, 0),
		nextID:   1,
		path:     path,
	}

	if path == "" {
		// No persistence if path is empty (mainly for tests if desired).
		return ms, nil
	}

	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			// No file yet – start empty.
			return ms, nil
		}
		return nil, err
	}
	defer f.Close()

	data, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	if len(data) == 0 {
		// Empty file – treat as no data.
		return ms, nil
	}

	if err := json.Unmarshal(data, &ms.expenses); err != nil {
		return nil, err
	}

	// Set nextID to max(existingID)+1
	var maxID int64
	for _, e := range ms.expenses {
		if e.ID > maxID {
			maxID = e.ID
		}
	}
	ms.nextID = maxID + 1

	return ms, nil
}

// ListExpenses returns a copy of all expenses.
func (m *MemoryStore) ListExpenses() ([]models.Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	out := make([]models.Expense, len(m.expenses))
	copy(out, m.expenses)
	return out, nil
}

// CreateExpense inserts a new expense, assigning ID and default date if needed.
func (m *MemoryStore) CreateExpense(e models.Expense) (models.Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if strings.TrimSpace(e.Date) == "" {
		e.Date = time.Now().Format("2006-01-02")
	}

	e.ID = m.nextID
	m.nextID++

	m.expenses = append(m.expenses, e)

	if err := m.save(); err != nil {
		return models.Expense{}, err
	}

	return e, nil
}

// GetExpense returns a single expense by ID.
func (m *MemoryStore) GetExpense(id int64) (models.Expense, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, e := range m.expenses {
		if e.ID == id {
			return e, nil
		}
	}
	return models.Expense{}, errors.New("expense not found")
}

// DeleteExpense removes a single expense by ID.
func (m *MemoryStore) DeleteExpense(id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, e := range m.expenses {
		if e.ID == id {
			// Remove from slice
			m.expenses = append(m.expenses[:i], m.expenses[i+1:]...)
			return m.save()
		}
	}

	return errors.New("expense not found")
}

// save writes the current expenses slice to the JSON file if a path is set.
func (m *MemoryStore) save() error {
	if m.path == "" {
		return nil
	}

	data, err := json.MarshalIndent(m.expenses, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(m.path, data, 0o644)
}
