package domain

import "time"

type Expense struct {
	ID          string         `json:"id"`
	GroupID     string         `json:"group_id"`
	PaidBy      string         `json:"paid_by"`
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

type OCRResult struct {
	Amount      float64 `json:"amount"`
	Description string  `json:"description"`
	Date        string  `json:"date"`
	RawText     string  `json:"raw_text"`
}
