package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Daudsaid/daud-expense-tracker-api/internal/models"
)

// ExpenseStore describes what the handler needs from a store.
type ExpenseStore interface {
	ListExpenses() ([]models.Expense, error)
	CreateExpense(e models.Expense) (models.Expense, error)
	GetExpense(id int64) (models.Expense, error)
	DeleteExpense(id int64) error
}

type ExpenseHandler struct {
	store ExpenseStore
}

type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: msg})
}

func NewExpenseHandler(store ExpenseStore) *ExpenseHandler {
	return &ExpenseHandler{store: store}
}

func (h *ExpenseHandler) RegisterRoutes(mux *http.ServeMux) {
	// Collection routes
	mux.HandleFunc("/expenses", h.expensesHandler)
	// Resource routes (with ID)
	mux.HandleFunc("/expenses/", h.expenseByIDHandler)
	mux.HandleFunc("/summary", h.summaryHandler)
}

func (h *ExpenseHandler) expensesHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/expenses" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleListExpenses(w)
	case http.MethodPost:
		h.handleCreateExpense(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
	}
}

// /expenses/{id}
func (h *ExpenseHandler) expenseByIDHandler(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, "/expenses/") {
		http.NotFound(w, r)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/expenses/")
	if idStr == "" {
		http.NotFound(w, r)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGetExpenseByID(w, id)
	case http.MethodDelete:
		h.handleDeleteExpenseByID(w, id)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
	}
}

func (h *ExpenseHandler) summaryHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/summary" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "method not allowed")
		return
	}

	expenses, err := h.store.ListExpenses()
	if err != nil {
		http.Error(w, "failed to list expenses", http.StatusInternalServerError)
		return
	}

	total := models.Total(expenses)
	perCategory := models.TotalsByCategory(expenses)

	resp := map[string]any{
		"count":        len(expenses),
		"total":        total,
		"per_category": perCategory,
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *ExpenseHandler) handleListExpenses(w http.ResponseWriter) {
	expenses, err := h.store.ListExpenses()
	if err != nil {
		http.Error(w, "failed to list expenses", http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) handleCreateExpense(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		http.Error(w, "Content-Type must be application/json", http.StatusBadRequest)
		return
	}

	var input struct {
		Amount   float64 `json:"amount"`
		Category string  `json:"category"`
		Note     string  `json:"note"`
		Date     string  `json:"date"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if input.Amount <= 0 {
		http.Error(w, "amount must be greater than 0", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(input.Category) == "" {
		http.Error(w, "category is required", http.StatusBadRequest)
		return
	}

	expense := models.Expense{
		Amount:   input.Amount,
		Category: input.Category,
		Note:     input.Note,
		Date:     input.Date, // store will fill default date if empty
	}

	created, err := h.store.CreateExpense(expense)
	if err != nil {
		http.Error(w, "failed to create expense", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ExpenseHandler) handleGetExpenseByID(w http.ResponseWriter, id int64) {
	exp, err := h.store.GetExpense(id)
	if err != nil {
		writeError(w, http.StatusNotFound, "expense not found")
		return
	}
	writeJSON(w, http.StatusOK, exp)
}

func (h *ExpenseHandler) handleDeleteExpenseByID(w http.ResponseWriter, id int64) {
	if err := h.store.DeleteExpense(id); err != nil {
		writeError(w, http.StatusNotFound, "expense not found")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// Just log to stderr via fmt.
		fmt.Println("error encoding JSON:", err)
	}
}
