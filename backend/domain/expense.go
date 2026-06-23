package domain

import "time"

type Expense struct {
	ID          string         `json:"id"`
	GroupID     string         `json:"group_id"`
	PaidBy      string         `json:"paid_by"`
	PaidByName  string         `json:"paid_by_name,omitempty"`
	Amount      float64        `json:"amount"`
	Description string         `json:"description,omitempty"`
	ReceiptURL  string         `json:"receipt_url,omitempty"`
	ExpenseDate string         `json:"expense_date"`
	CreatedAt   time.Time      `json:"created_at"`
	Splits      []ExpenseSplit `json:"splits,omitempty"`
}

type ExpenseSplit struct {
	ExpenseID  string  `json:"expense_id,omitempty"`
	UserID     string  `json:"user_id"`
	AmountOwed float64 `json:"amount_owed"`
}

type Balance struct {
	UserID  string  `json:"user_id"`
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type GroupBalance struct {
	GroupID   string  `json:"group_id"`
	GroupName string  `json:"group_name"`
	Balance   float64 `json:"balance"`
}

type PersonBalance struct {
	UserID   string  `json:"user_id"`
	UserName string  `json:"user_name"`
	Balance  float64 `json:"balance"`
}

type Settlement struct {
	ID        string    `json:"id"`
	PayerID   string    `json:"payer_id"`
	PayeeID   string    `json:"payee_id"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}

type OCRResult struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	RawText     string  `json:"raw_text"`
}
