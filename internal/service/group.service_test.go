package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
)

type fakeGroupRepository struct {
	createFunc       func(group *model.Group) error
	findByUserFunc   func(userID int64) ([]model.Group, error)
	findByIDUserFunc func(id int64, userID int64) (*model.Group, error)
	updateFunc       func(group *model.Group) error
	deleteFunc       func(id int64, userID int64) error
}

func (f *fakeGroupRepository) Create(group *model.Group) error {
	if f.createFunc != nil {
		return f.createFunc(group)
	}
	return nil
}

func (f *fakeGroupRepository) FindByUser(
	userID int64,
) ([]model.Group, error) {
	if f.findByUserFunc != nil {
		return f.findByUserFunc(userID)
	}
	return nil, nil
}

func (f *fakeGroupRepository) FindByIDAndUser(
	id int64,
	userID int64,
) (*model.Group, error) {
	if f.findByIDUserFunc != nil {
		return f.findByIDUserFunc(id, userID)
	}
	return nil, nil
}

func (f *fakeGroupRepository) Update(
	group *model.Group,
) error {
	if f.updateFunc != nil {
		return f.updateFunc(group)
	}
	return nil
}

func (f *fakeGroupRepository) Delete(
	id int64,
	userID int64,
) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(id, userID)
	}
	return nil
}

func TestGroupService_CreateGroup(t *testing.T) {
	description := "My favorite anime"

	var createdGroup *model.Group

	repo := &fakeGroupRepository{
		createFunc: func(group *model.Group) error {
			createdGroup = group
			group.ID = 1
			return nil
		},
	}

	service := NewGroupService(repo)

	result, err := service.CreateGroup(
		context.Background(),
		10,
		"Favorites",
		&description,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected group, got nil")
	}

	if result.ID != 1 {
		t.Errorf("expected ID 1, got %d", result.ID)
	}

	if result.UserID != 10 {
		t.Errorf("expected UserID 10, got %d", result.UserID)
	}

	if result.Name != "Favorites" {
		t.Errorf(
			"expected name Favorites, got %s",
			result.Name,
		)
	}

	if result.Description == nil {
		t.Fatal("expected description, got nil")
	}

	if *result.Description != description {
		t.Errorf(
			"expected description %q, got %q",
			description,
			*result.Description,
		)
	}

	if createdGroup == nil {
		t.Fatal("expected repository Create to be called")
	}
}

func TestGroupService_CreateGroup_Error(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakeGroupRepository{
		createFunc: func(group *model.Group) error {
			return expectedErr
		},
	}

	service := NewGroupService(repo)

	result, err := service.CreateGroup(
		context.Background(),
		10,
		"Favorites",
		nil,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupService_GetGroups(t *testing.T) {
	groups := []model.Group{
		{
			ID:     1,
			UserID: 10,
			Name:   "Favorites",
		},
		{
			ID:     2,
			UserID: 10,
			Name:   "Watching",
		},
	}

	var receivedUserID int64

	repo := &fakeGroupRepository{
		findByUserFunc: func(userID int64) ([]model.Group, error) {
			receivedUserID = userID
			return groups, nil
		},
	}

	service := NewGroupService(repo)

	result, err := service.GetGroups(
		context.Background(),
		10,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 groups, got %d",
			len(result),
		)
	}

	if result[0].Name != "Favorites" {
		t.Errorf(
			"expected Favorites, got %s",
			result[0].Name,
		)
	}

	if result[1].Name != "Watching" {
		t.Errorf(
			"expected Watching, got %s",
			result[1].Name,
		)
	}
}

func TestGroupService_GetGroups_Error(t *testing.T) {
	expectedErr := errors.New("database error")

	repo := &fakeGroupRepository{
		findByUserFunc: func(userID int64) ([]model.Group, error) {
			return nil, expectedErr
		},
	}

	service := NewGroupService(repo)

	result, err := service.GetGroups(
		context.Background(),
		10,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupService_GetGroup(t *testing.T) {
	expectedGroup := &model.Group{
		ID:     1,
		UserID: 10,
		Name:   "Favorites",
	}

	var receivedID int64
	var receivedUserID int64

	repo := &fakeGroupRepository{
		findByIDUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			receivedID = id
			receivedUserID = userID
			return expectedGroup, nil
		},
	}

	service := NewGroupService(repo)

	result, err := service.GetGroup(
		context.Background(),
		10,
		1,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result != expectedGroup {
		t.Fatalf(
			"expected group %p, got %p",
			expectedGroup,
			result,
		)
	}

	if receivedID != 1 {
		t.Errorf(
			"expected ID 1, got %d",
			receivedID,
		)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}
}

func TestGroupService_GetGroup_NotFound(t *testing.T) {
	repo := &fakeGroupRepository{
		findByIDUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return nil, repository.ErrGroupNotFound
		},
	}

	service := NewGroupService(repo)

	result, err := service.GetGroup(
		context.Background(),
		10,
		999,
	)

	if !errors.Is(err, repository.ErrGroupNotFound) {
		t.Fatalf(
			"expected ErrGroupNotFound, got %v",
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupService_UpdateGroup(t *testing.T) {
	description := "Updated description"

	group := &model.Group{
		ID:          1,
		UserID:      10,
		Name:        "Old Name",
		Description: nil,
	}

	var updatedGroup *model.Group

	repo := &fakeGroupRepository{
		findByIDUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return group, nil
		},
		updateFunc: func(group *model.Group) error {
			updatedGroup = group
			return nil
		},
	}

	service := NewGroupService(repo)

	newName := "New Name"

	result, err := service.UpdateGroup(
		context.Background(),
		10,
		1,
		&newName,
		&description,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedGroup == nil {
		t.Fatal("expected repository Update to be called")
	}

	if updatedGroup.Name != "New Name" {
		t.Errorf(
			"expected name New Name, got %s",
			updatedGroup.Name,
		)
	}

	if updatedGroup.Description == nil {
		t.Fatal("expected description, got nil")
	}

	if *updatedGroup.Description != description {
		t.Errorf(
			"expected description %q, got %q",
			description,
			*updatedGroup.Description,
		)
	}

	if result != group {
		t.Fatalf(
			"expected returned group %p, got %p",
			group,
			result,
		)
	}
}

func TestGroupService_UpdateGroup_PartialUpdate(t *testing.T) {
	description := "Original description"

	group := &model.Group{
		ID:          1,
		UserID:      10,
		Name:        "Old Name",
		Description: &description,
	}

	var updatedGroup *model.Group

	repo := &fakeGroupRepository{
		findByIDUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return group, nil
		},
		updateFunc: func(group *model.Group) error {
			updatedGroup = group
			return nil
		},
	}

	service := NewGroupService(repo)

	newName := "New Name"

	result, err := service.UpdateGroup(
		context.Background(),
		10,
		1,
		&newName,
		nil,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if updatedGroup.Name != "New Name" {
		t.Errorf(
			"expected name New Name, got %s",
			updatedGroup.Name,
		)
	}

	if updatedGroup.Description == nil {
		t.Fatal("expected description to remain, got nil")
	}

	if *updatedGroup.Description != description {
		t.Errorf(
			"expected description %q, got %q",
			description,
			*updatedGroup.Description,
		)
	}

	if result != group {
		t.Fatalf(
			"expected returned group %p, got %p",
			group,
			result,
		)
	}
}

func TestGroupService_UpdateGroup_NotFound(t *testing.T) {
	repo := &fakeGroupRepository{
		findByIDUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return nil, repository.ErrGroupNotFound
		},
	}

	service := NewGroupService(repo)

	result, err := service.UpdateGroup(
		context.Background(),
		10,
		999,
		nil,
		nil,
	)

	if !errors.Is(err, repository.ErrGroupNotFound) {
		t.Fatalf(
			"expected ErrGroupNotFound, got %v",
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupService_UpdateGroup_UpdateError(t *testing.T) {
	expectedErr := errors.New("update failed")

	group := &model.Group{
		ID:     1,
		UserID: 10,
		Name:   "Favorites",
	}

	repo := &fakeGroupRepository{
		findByIDUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return group, nil
		},
		updateFunc: func(group *model.Group) error {
			return expectedErr
		},
	}

	service := NewGroupService(repo)

	name := "Updated"

	result, err := service.UpdateGroup(
		context.Background(),
		10,
		1,
		&name,
		nil,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupService_DeleteGroup(t *testing.T) {
	var receivedID int64
	var receivedUserID int64

	repo := &fakeGroupRepository{
		deleteFunc: func(
			id int64,
			userID int64,
		) error {
			receivedID = id
			receivedUserID = userID
			return nil
		},
	}

	service := NewGroupService(repo)

	err := service.DeleteGroup(
		context.Background(),
		10,
		1,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if receivedID != 1 {
		t.Errorf(
			"expected ID 1, got %d",
			receivedID,
		)
	}

	if receivedUserID != 10 {
		t.Errorf(
			"expected userID 10, got %d",
			receivedUserID,
		)
	}
}

func TestGroupService_DeleteGroup_NotFound(t *testing.T) {
	repo := &fakeGroupRepository{
		deleteFunc: func(
			id int64,
			userID int64,
		) error {
			return repository.ErrGroupNotFound
		},
	}

	service := NewGroupService(repo)

	err := service.DeleteGroup(
		context.Background(),
		10,
		999,
	)

	if !errors.Is(err, repository.ErrGroupNotFound) {
		t.Fatalf(
			"expected ErrGroupNotFound, got %v",
			err,
		)
	}
}

func TestGroupService_DeleteGroup_Error(t *testing.T) {
	expectedErr := errors.New("delete failed")

	repo := &fakeGroupRepository{
		deleteFunc: func(
			id int64,
			userID int64,
		) error {
			return expectedErr
		},
	}

	service := NewGroupService(repo)

	err := service.DeleteGroup(
		context.Background(),
		10,
		1,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}
