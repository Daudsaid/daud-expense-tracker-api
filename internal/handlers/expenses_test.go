package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Daudsaid/daud-expense-tracker-api/internal/models"
)

type mockStore struct {
	expenses map[int64]models.Expense
}

func newMockStore() *mockStore {
	return &mockStore{
		expenses: map[int64]models.Expense{
			1: {ID: 1, Amount: 10, Category: "Food", Note: "Test", Date: "2025-11-17"},
		},
	}
}

func (m *mockStore) ListExpenses() ([]models.Expense, error) {
	out := make([]models.Expense, 0, len(m.expenses))
	for _, e := range m.expenses {
		out = append(out, e)
	}
	return out, nil
}

func (m *mockStore) CreateExpense(e models.Expense) (models.Expense, error) {
	if e.ID == 0 {
		e.ID = int64(len(m.expenses) + 1)
	}
	m.expenses[e.ID] = e
	return e, nil
}

func (m *mockStore) GetExpense(id int64) (models.Expense, error) {
	e, ok := m.expenses[id]
	if !ok {
		return models.Expense{}, fmt.Errorf("not found")
	}
	return e, nil
}

func (m *mockStore) DeleteExpense(id int64) error {
	if _, ok := m.expenses[id]; !ok {
		return fmt.Errorf("not found")
	}
	delete(m.expenses, id)
	return nil
}

func TestGetExpenseByID_OK(t *testing.T) {
	store := newMockStore()
	h := NewExpenseHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/expenses/1", nil)
	rr := httptest.NewRecorder()

	h.expenseByIDHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rr.Code)
	}

	var e models.Expense
	if err := json.Unmarshal(rr.Body.Bytes(), &e); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if e.ID != 1 {
		t.Fatalf("expected ID 1, got %d", e.ID)
	}
}

func TestGetExpenseByID_NotFound(t *testing.T) {
	store := newMockStore()
	h := NewExpenseHandler(store)

	req := httptest.NewRequest(http.MethodGet, "/expenses/999", nil)
	rr := httptest.NewRecorder()

	h.expenseByIDHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}

func TestDeleteExpenseByID_OK(t *testing.T) {
	store := newMockStore()
	h := NewExpenseHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/expenses/1", nil)
	rr := httptest.NewRecorder()

	h.expenseByIDHandler(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", rr.Code)
	}

	// Ensure it's gone
	_, err := store.GetExpense(1)
	if err == nil {
		t.Fatalf("expected error after delete, got nil")
	}
}

func TestDeleteExpenseByID_NotFound(t *testing.T) {
	store := newMockStore()
	h := NewExpenseHandler(store)

	req := httptest.NewRequest(http.MethodDelete, "/expenses/999", nil)
	rr := httptest.NewRecorder()

	h.expenseByIDHandler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", rr.Code)
	}
}
