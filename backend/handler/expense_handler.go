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
	PaidBy       string   `json:"paid_by"`
	Amount       float64  `json:"amount"`
	Description  string   `json:"description"`
	ReceiptURL   string   `json:"receipt_url"`
	ExpenseDate  string   `json:"expense_date"`
	SplitUserIDs []string `json:"split_user_ids"`
}

func (h *ExpenseHandler) GetPersonalBalances(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)

	balances, err := h.expenseService.GetPersonalBalances(r.Context(), userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "erro ao buscar balanços")
		return
	}

	writeJSON(w, http.StatusOK, balances)
}

func (h *ExpenseHandler) Settle(w http.ResponseWriter, r *http.Request) {
	payerID := r.Context().Value(UserIDKey).(string)

	var req struct {
		PayeeID string  `json:"payee_id"`
		Amount  float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PayeeID == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "payee_id e amount são obrigatórios")
		return
	}

	if err := h.expenseService.Settle(r.Context(), payerID, req.PayeeID, req.Amount); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "pagamento registrado"})
}

func (h *ExpenseHandler) SettleReceived(w http.ResponseWriter, r *http.Request) {
	payeeID := r.Context().Value(UserIDKey).(string)

	var req struct {
		PayerID string  `json:"payer_id"`
		Amount  float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.PayerID == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "payer_id e amount são obrigatórios")
		return
	}

	if err := h.expenseService.Settle(r.Context(), req.PayerID, payeeID, req.Amount); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{"message": "recebimento registrado"})
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
	}

	created, err := h.expenseService.CreateExpense(r.Context(), userID, expense, req.SplitUserIDs)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, created)
}

func (h *ExpenseHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(UserIDKey).(string)
	expenseID := chi.URLParam(r, "expenseID")

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
