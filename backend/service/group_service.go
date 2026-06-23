package service

import (
	"context"
	"errors"

	"racha-historico/domain"
	"racha-historico/repository"
)

type GroupService struct {
	groupRepo *repository.GroupRepository
}

func NewGroupService(groupRepo *repository.GroupRepository) *GroupService {
	return &GroupService{groupRepo: groupRepo}
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

func (s *GroupService) AddMember(ctx context.Context, groupID string, newUserID string, requestingUserID string) error {
	isMember, err := s.groupRepo.IsMember(ctx, groupID, requestingUserID)
	if err != nil {
		return err
	}
	if !isMember {
		return errors.New("você não faz parte deste grupo")
	}
	return s.groupRepo.AddMember(ctx, groupID, newUserID)
}
