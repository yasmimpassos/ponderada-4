package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"

	"racha-historico/domain"
	"racha-historico/service"
)

type ExpenseHandler struct {
	expenseService *service.ExpenseService
	ocrService     *service.OCRService
}

func NewExpenseHandler(expenseService *service.ExpenseService, ocrService *service.OCRService) *ExpenseHandler {
	return &ExpenseHandler{
		expenseService: expenseService,
		ocrService:     ocrService,
	}
}

type createExpenseRequest struct {
	PaidBy      string               `json:"paid_by"`
	Amount      float64              `json:"amount"`
	Description string               `json:"description"`
	ReceiptURL  string               `json:"receipt_url"`
	ExpenseDate string               `json:"expense_date"`
	Splits      []domain.ExpenseSplit `json:"splits"`
}

func (h *ExpenseHandler) ListGroupExpenses(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := chi.URLParam(r, "id")

	expenses, err := h.expenseService.GetGroupExpenses(r.Context(), groupID, userID)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) CreateExpense(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := chi.URLParam(r, "id")

	var req createExpenseRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil || req.Amount <= 0 || req.ExpenseDate == "" {
		writeError(w, http.StatusBadRequest, "amount e expense_date são obrigatórios")
		return
	}

	expense := &domain.Expense{
		GroupID:     groupID,
		PaidBy:      req.PaidBy,
		Amount:      req.Amount,
		Description: req.Description,
		ReceiptURL:  req.ReceiptURL,
		ExpenseDate: req.ExpenseDate,
		Splits:      req.Splits,
	}

	created, err := h.expenseService.CreateExpense(r.Context(), userID, expense)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	expenseID := chi.URLParam(r, "id")

	err := h.expenseService.DeleteExpense(r.Context(), expenseID, userID)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"message": "gasto removido com sucesso"})
}

func (h *ExpenseHandler) GetGroupBalances(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	groupID := chi.URLParam(r, "id")

	balances, err := h.expenseService.GetGroupBalances(r.Context(), groupID, userID)
	if err != nil {
		writeError(w, http.StatusForbidden, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, balances)
}

func (h *ExpenseHandler) ProcessOCR(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		writeError(w, http.StatusBadRequest, "erro ao ler formulário")
		return
	}

	file, _, err := r.FormFile("image")
	if err != nil {
		writeError(w, http.StatusBadRequest, "campo 'image' não encontrado")
		return
	}
	defer file.Close()

	imageBytes, err := io.ReadAll(file)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao ler imagem")
		return
	}

	result, err := h.ocrService.ExtractFromImage(imageBytes)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao processar imagem")
		return
	}

	writeJSON(w, http.StatusOK, result)
}
