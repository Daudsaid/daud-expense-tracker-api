package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Daudsaid/daud-expense-tracker-api/internal/models"
)

func TestMemoryStore_GetAndDelete(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "test_expenses.json")

	store, err := NewMemoryStore(path)
	if err != nil {
		t.Fatalf("failed to create store: %v", err)
	}

	created, err := store.CreateExpense(models.Expense{
		Amount:   9.99,
		Category: "Test",
		Note:     "Unit test",
	})
	if err != nil {
		t.Fatalf("failed to create expense: %v", err)
	}

	// Get by ID
	got, err := store.GetExpense(created.ID)
	if err != nil {
		t.Fatalf("GetExpense returned error: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected ID %d, got %d", created.ID, got.ID)
	}

	// Delete
	if err := store.DeleteExpense(created.ID); err != nil {
		t.Fatalf("DeleteExpense returned error: %v", err)
	}

	// Ensure it's not found
	if _, err := store.GetExpense(created.ID); err == nil {
		t.Fatalf("expected error for deleted expense, got nil")
	}

	// Optional: ensure JSON file is written
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected JSON file to exist, got error: %v", err)
	}
}
