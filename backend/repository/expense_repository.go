package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"racha-historico/domain"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(ctx context.Context, expense *domain.Expense) error {
	expense.ID = uuid.New().String()
	expense.CreatedAt = time.Now()

	query := `
		INSERT INTO expenses (id, group_id, paid_by, amount, description, receipt_url, expense_date, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		expense.ID, expense.GroupID, expense.PaidBy, expense.Amount,
		expense.Description, expense.ReceiptURL, expense.ExpenseDate, expense.CreatedAt,
	)
	return err
}

func (r *ExpenseRepository) CreateSplit(ctx context.Context, split *domain.ExpenseSplit) error {
	query := `INSERT INTO expense_splits (expense_id, user_id, amount_owed) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, split.ExpenseID, split.UserID, split.AmountOwed)
	return err
}

func (r *ExpenseRepository) FindByGroupID(ctx context.Context, groupID string) ([]*domain.Expense, error) {
	query := `
		SELECT id, group_id, paid_by, amount, COALESCE(description, ''), COALESCE(receipt_url, ''),
		       DATE_FORMAT(expense_date, '%Y-%m-%d'), created_at
		FROM expenses
		WHERE group_id = ?
		ORDER BY expense_date DESC, created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var expenses []*domain.Expense
	for rows.Next() {
		expense := &domain.Expense{}
		err := rows.Scan(
			&expense.ID, &expense.GroupID, &expense.PaidBy,
			&expense.Amount, &expense.Description, &expense.ReceiptURL,
			&expense.ExpenseDate, &expense.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		expenses = append(expenses, expense)
	}
	return expenses, nil
}

func (r *ExpenseRepository) FindByID(ctx context.Context, id string) (*domain.Expense, error) {
	expense := &domain.Expense{}
	query := `
		SELECT id, group_id, paid_by, amount, COALESCE(description, ''), COALESCE(receipt_url, ''),
		       DATE_FORMAT(expense_date, '%Y-%m-%d'), created_at
		FROM expenses
		WHERE id = ?
	`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&expense.ID, &expense.GroupID, &expense.PaidBy,
		&expense.Amount, &expense.Description, &expense.ReceiptURL,
		&expense.ExpenseDate, &expense.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return expense, nil
}

func (r *ExpenseRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM expenses WHERE id = ?`, id)
	return err
}

func (r *ExpenseRepository) GetGroupBalances(ctx context.Context, groupID string) ([]*domain.Balance, error) {
	query := `
		SELECT
			u.id,
			u.name,
			COALESCE(SUM(CASE WHEN e.paid_by = u.id THEN e.amount ELSE 0 END), 0)
			- COALESCE(SUM(es.amount_owed), 0) AS balance
		FROM users u
		JOIN group_members gm ON gm.user_id = u.id
		LEFT JOIN expenses e ON e.group_id = gm.group_id AND e.group_id = ?
		LEFT JOIN expense_splits es ON es.expense_id = e.id AND es.user_id = u.id
		WHERE gm.group_id = ?
		GROUP BY u.id, u.name
	`
	rows, err := r.db.QueryContext(ctx, query, groupID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var balances []*domain.Balance
	for rows.Next() {
		balance := &domain.Balance{}
		err := rows.Scan(&balance.UserID, &balance.Name, &balance.Balance)
		if err != nil {
			return nil, err
		}
		balances = append(balances, balance)
	}
	return balances, nil
}
