package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"racha-historico/domain"
)

type GroupRepository struct {
	db *sql.DB
}

func NewGroupRepository(db *sql.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *domain.Group) error {
	group.ID = uuid.New().String()
	group.CreatedAt = time.Now()

	query := `INSERT INTO user_groups (id, name, created_by, created_at) VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, group.ID, group.Name, group.CreatedBy, group.CreatedAt)
	return err
}

func (r *GroupRepository) FindByID(ctx context.Context, id string) (*domain.Group, error) {
	group := &domain.Group{}
	query := `SELECT id, name, created_by, created_at FROM user_groups WHERE id = ?`
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return group, nil
}

func (r *GroupRepository) FindByUserID(ctx context.Context, userID string) ([]*domain.Group, error) {
	query := `
		SELECT g.id, g.name, g.created_by, g.created_at
		FROM user_groups g
		JOIN group_members gm ON gm.group_id = g.id
		WHERE gm.user_id = ?
		ORDER BY g.created_at DESC
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user_groups []*domain.Group
	for rows.Next() {
		group := &domain.Group{}
		err := rows.Scan(&group.ID, &group.Name, &group.CreatedBy, &group.CreatedAt)
		if err != nil {
			return nil, err
		}
		user_groups = append(user_groups, group)
	}
	return user_groups, nil
}

func (r *GroupRepository) AddMember(ctx context.Context, groupID string, userID string) error {
	query := `INSERT IGNORE INTO group_members (group_id, user_id, joined_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, groupID, userID, time.Now())
	return err
}

func (r *GroupRepository) IsMember(ctx context.Context, groupID string, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?`
	err := r.db.QueryRowContext(ctx, query, groupID, userID).Scan(&count)
	return count > 0, err
}

func (r *GroupRepository) GetMemberFCMTokens(ctx context.Context, groupID string) ([]string, error) {
	query := `
		SELECT u.fcm_token
		FROM users u
		JOIN group_members gm ON gm.user_id = u.id
		WHERE gm.group_id = ? AND u.fcm_token IS NOT NULL AND u.fcm_token != ''
	`
	rows, err := r.db.QueryContext(ctx, query, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []string
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	return tokens, nil
}
