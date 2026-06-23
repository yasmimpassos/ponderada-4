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
		SELECT e.id, e.group_id, e.paid_by, u.name,
		       e.amount, COALESCE(e.description, ''), COALESCE(e.receipt_url, ''),
		       DATE_FORMAT(e.expense_date, '%Y-%m-%d'), e.created_at
		FROM expenses e
		JOIN users u ON u.id = e.paid_by
		WHERE e.group_id = ?
		ORDER BY e.expense_date DESC, e.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	expenses := make([]*domain.Expense, 0)
	for rows.Next() {
		expense := &domain.Expense{}
		err := rows.Scan(
			&expense.ID, &expense.GroupID, &expense.PaidBy, &expense.PaidByName,
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

func (r *ExpenseRepository) GetPersonalBalances(ctx context.Context, userID string) ([]*domain.PersonBalance, error) {
	query := `
		SELECT other_user_id, other_name, SUM(net_amount) AS balance
		FROM (
			SELECT es.user_id AS other_user_id, u.name AS other_name, es.amount_owed AS net_amount
			FROM expenses e
			JOIN expense_splits es ON es.expense_id = e.id
			JOIN users u ON u.id = es.user_id
			WHERE e.paid_by = ? AND es.user_id != ?

			UNION ALL

			SELECT e.paid_by AS other_user_id, u.name AS other_name, -es.amount_owed AS net_amount
			FROM expense_splits es
			JOIN expenses e ON e.id = es.expense_id
			JOIN users u ON u.id = e.paid_by
			WHERE es.user_id = ? AND e.paid_by != ?

			UNION ALL

			SELECT s.payee_id AS other_user_id, u.name AS other_name, s.amount AS net_amount
			FROM settlements s
			JOIN users u ON u.id = s.payee_id
			WHERE s.payer_id = ?

			UNION ALL

			SELECT s.payer_id AS other_user_id, u.name AS other_name, -s.amount AS net_amount
			FROM settlements s
			JOIN users u ON u.id = s.payer_id
			WHERE s.payee_id = ?
		) combined
		GROUP BY other_user_id, other_name
		HAVING ABS(SUM(net_amount)) > 0.001
		ORDER BY SUM(net_amount) ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]*domain.PersonBalance, 0)
	for rows.Next() {
		b := &domain.PersonBalance{}
		if err := rows.Scan(&b.UserID, &b.UserName, &b.Balance); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *ExpenseRepository) CreateSettlement(ctx context.Context, s *domain.Settlement) error {
	s.ID = uuid.New().String()
	s.CreatedAt = time.Now()
	query := `INSERT INTO settlements (id, payer_id, payee_id, amount, created_at) VALUES (?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, s.ID, s.PayerID, s.PayeeID, s.Amount, s.CreatedAt)
	return err
}

func (r *ExpenseRepository) GetUserGlobalBalances(ctx context.Context, userID string) ([]*domain.GroupBalance, error) {
	query := `
		SELECT
			g.id,
			g.name,
			COALESCE(SUM(CASE WHEN e.paid_by = ? THEN e.amount ELSE 0 END), 0)
			- COALESCE(SUM(CASE WHEN es.user_id = ? THEN es.amount_owed ELSE 0 END), 0) AS balance
		FROM user_groups g
		JOIN group_members gm ON gm.group_id = g.id AND gm.user_id = ?
		LEFT JOIN expenses e ON e.group_id = g.id
		LEFT JOIN expense_splits es ON es.expense_id = e.id AND es.user_id = ?
		GROUP BY g.id, g.name
		HAVING balance != 0
		ORDER BY balance ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]*domain.GroupBalance, 0)
	for rows.Next() {
		b := &domain.GroupBalance{}
		if err := rows.Scan(&b.GroupID, &b.GroupName, &b.Balance); err != nil {
			return nil, err
		}
		balances = append(balances, b)
	}
	return balances, nil
}

func (r *ExpenseRepository) GetGroupBalances(ctx context.Context, groupID string) ([]*domain.Balance, error) {
	query := `
		SELECT
			u.id,
			u.name,
			COALESCE(SUM(CASE WHEN e.paid_by = u.id THEN e.amount ELSE 0 END), 0)
			- COALESCE(SUM(es.amount_owed), 0)
			+ COALESCE((
				SELECT SUM(s.amount) FROM settlements s
				JOIN group_members gm2 ON gm2.user_id = s.payee_id AND gm2.group_id = ?
				WHERE s.payer_id = u.id
			), 0)
			- COALESCE((
				SELECT SUM(s.amount) FROM settlements s
				JOIN group_members gm3 ON gm3.user_id = s.payer_id AND gm3.group_id = ?
				WHERE s.payee_id = u.id
			), 0) AS balance
		FROM users u
		JOIN group_members gm ON gm.user_id = u.id
		LEFT JOIN expenses e ON e.group_id = gm.group_id AND e.group_id = ?
		LEFT JOIN expense_splits es ON es.expense_id = e.id AND es.user_id = u.id
		WHERE gm.group_id = ?
		GROUP BY u.id, u.name
	`
	rows, err := r.db.QueryContext(ctx, query, groupID, groupID, groupID, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make([]*domain.Balance, 0)
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
