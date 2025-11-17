package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Daudsaid/daud-expense-tracker-api/internal/handlers"
	"github.com/Daudsaid/daud-expense-tracker-api/internal/storage"
)

func main() {
	// Default data file, override with ENV if needed.
	dataFile := "expenses.json"
	if v := os.Getenv("EXPENSES_DATA_FILE"); v != "" {
		dataFile = v
	}

	store, err := storage.NewStore(dataFile)
	if err != nil {
		log.Fatalf("failed to create store: %v", err)
	}

	mux := http.NewServeMux()

	expenseHandler := handlers.NewExpenseHandler(store)
	expenseHandler.RegisterRoutes(mux)

	addr := ":8080"
	fmt.Printf("✅ Expense Tracker API running on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}
