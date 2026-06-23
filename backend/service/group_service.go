package service

import (
	"context"
	"errors"

	"racha-historico/domain"
	"racha-historico/repository"
)

type GroupService struct {
	groupRepo *repository.GroupRepository
	userRepo  *repository.UserRepository
}

func NewGroupService(groupRepo *repository.GroupRepository, userRepo *repository.UserRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo, userRepo: userRepo}
}

func (s *GroupService) CreateGroup(ctx context.Context, name string, createdByUserID string, memberIDs []string) (*domain.Group, error) {
	group := &domain.Group{
		Name:      name,
		CreatedBy: createdByUserID,
	}

	err := s.groupRepo.Create(ctx, group)
	if err != nil {
		return nil, err
	}

	err = s.groupRepo.AddMember(ctx, group.ID, createdByUserID)
	if err != nil {
		return nil, err
	}

	for _, memberID := range memberIDs {
		if memberID == createdByUserID {
			continue
		}
		err = s.groupRepo.AddMember(ctx, group.ID, memberID)
		if err != nil {
			return nil, err
		}
	}

	return group, nil
}

func (s *GroupService) GetUserGroups(ctx context.Context, userID string) ([]*domain.Group, error) {
	return s.groupRepo.FindByUserID(ctx, userID)
}

func (s *GroupService) GetGroupDetails(ctx context.Context, groupID string, userID string) (*domain.Group, error) {
	isMember, err := s.groupRepo.IsMember(ctx, groupID, userID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, errors.New("você não faz parte deste grupo")
	}
	return s.groupRepo.FindByID(ctx, groupID)
}

func (s *GroupService) AddMemberByEmail(ctx context.Context, groupID string, email string, requestingUserID string) error {
	isMember, err := s.groupRepo.IsMember(ctx, groupID, requestingUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("você não faz parte deste grupo")
	}

	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return errors.New("usuário não encontrado")
	}

	return s.groupRepo.AddMember(ctx, groupID, user.ID)
}

func (s *GroupService) GetMembers(ctx context.Context, groupID string) ([]*domain.User, error) {
	return s.groupRepo.GetMembers(ctx, groupID)
}

func (s *GroupService) JoinGroup(ctx context.Context, groupID string, userID string) error {
	_, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return errors.New("grupo não encontrado")
	}
	return s.groupRepo.AddMember(ctx, groupID, userID)
}
