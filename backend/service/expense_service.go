package service

import (
	"context"
	"errors"
	"fmt"

	"racha-historico/domain"
	"racha-historico/repository"
)

type ExpenseService struct {
	expenseRepo         *repository.ExpenseRepository
	groupRepo           *repository.GroupRepository
	notificationService *NotificationService
}

func NewExpenseService(expenseRepo *repository.ExpenseRepository, groupRepo *repository.GroupRepository, notificationService *NotificationService) *ExpenseService {
	return &ExpenseService{
		expenseRepo:         expenseRepo,
		groupRepo:           groupRepo,
		notificationService: notificationService,
	}
}

func (s *ExpenseService) CreateExpense(ctx context.Context, requestingUserID string, expense *domain.Expense) (*domain.Expense, error) {
	if expense.PaidBy == "" {
		expense.PaidBy = requestingUserID
	}

	isMember, err := s.groupRepo.IsMember(ctx, expense.GroupID, expense.PaidBy)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("usuário não faz parte deste grupo")
	}

	err = s.expenseRepo.Create(ctx, expense)
	if err != nil {
		return nil, err
	}

	for i := range expense.Splits {
		expense.Splits[i].ExpenseID = expense.ID
		err = s.expenseRepo.CreateSplit(ctx, &expense.Splits[i])
		if err != nil {
			return nil, err
		}
	}

	go s.sendNewExpenseNotification(expense)

	return expense, nil
}

func (s *ExpenseService) sendNewExpenseNotification(expense *domain.Expense) {
	ctx := context.Background()
	tokens, err := s.groupRepo.GetMemberFCMTokens(ctx, expense.GroupID)
	if err != nil {
		return
	}
	title := "Novo gasto no grupo"
	body := fmt.Sprintf("R$ %.2f registrado: %s", expense.Amount, expense.Description)
	s.notificationService.NotifyGroupMembers(tokens, title, body)
}

func (s *ExpenseService) GetGroupExpenses(ctx context.Context, groupID string, userID string) ([]*domain.Expense, error) {
	isMember, err := s.groupRepo.IsMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("você não faz parte deste grupo")
	}
	return s.expenseRepo.FindByGroupID(ctx, groupID)
}

func (s *ExpenseService) DeleteExpense(ctx context.Context, expenseID string, userID string) error {
	expense, err := s.expenseRepo.FindByID(ctx, expenseID)
	if err != nil {
		return errors.New("gasto não encontrado")
	}

	isMember, err := s.groupRepo.IsMember(ctx, expense.GroupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("você não faz parte deste grupo")
	}

	return s.expenseRepo.Delete(ctx, expenseID)
}

func (s *ExpenseService) GetGroupBalances(ctx context.Context, groupID string, userID string) ([]*domain.Balance, error) {
	isMember, err := s.groupRepo.IsMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("você não faz parte deste grupo")
	}
	return s.expenseRepo.GetGroupBalances(ctx, groupID)
}
