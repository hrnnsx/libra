package service

import (
	"context"

	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
)

type GroupService interface {
	CreateGroup(
		ctx context.Context,
		userID int64,
		name string,
		description *string,
	) (*model.Group, error)

	GetGroups(
		ctx context.Context,
		userID int64,
	) ([]model.Group, error)

	GetGroup(
		ctx context.Context,
		userID int64,
		id int64,
	) (*model.Group, error)

	UpdateGroup(
		ctx context.Context,
		userID int64,
		id int64,
		name *string,
		description *string,
	) (*model.Group, error)

	DeleteGroup(
		ctx context.Context,
		userID int64,
		id int64,
	) error
}

type groupService struct {
	groupRepository repository.GroupRepository
}

func NewGroupService(
	groupRepository repository.GroupRepository,
) GroupService {
	return &groupService{
		groupRepository: groupRepository,
	}
}

func (s *groupService) CreateGroup(
	ctx context.Context,
	userID int64,
	name string,
	description *string,
) (*model.Group, error) {
	group := &model.Group{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	if err := s.groupRepository.Create(group); err != nil {
		return nil, err
	}

	return group, nil
}

func (s *groupService) GetGroups(
	ctx context.Context,
	userID int64,
) ([]model.Group, error) {
	return s.groupRepository.FindByUser(userID)
}

func (s *groupService) GetGroup(
	ctx context.Context,
	userID int64,
	id int64,
) (*model.Group, error) {
	return s.groupRepository.FindByIDAndUser(
		id,
		userID,
	)
}

func (s *groupService) UpdateGroup(
	ctx context.Context,
	userID int64,
	id int64,
	name *string,
	description *string,
) (*model.Group, error) {
	group, err := s.groupRepository.FindByIDAndUser(
		id,
		userID,
	)

	if err != nil {
		return nil, err
	}

	if name != nil {
		group.Name = *name
	}

	if description != nil {
		group.Description = description
	}

	if err := s.groupRepository.Update(group); err != nil {
		return nil, err
	}

	return s.groupRepository.FindByIDAndUser(
		id,
		userID,
	)
}

func (s *groupService) DeleteGroup(
	ctx context.Context,
	userID int64,
	id int64,
) error {
	return s.groupRepository.Delete(
		id,
		userID,
	)
}
