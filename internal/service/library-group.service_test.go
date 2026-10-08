package service

import (
	"context"
	"errors"
	"testing"

	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
)

type fakeGroupLibraryGroupRepository struct {
	findByIDAndUserFunc func(
		id int64,
		userID int64,
	) (*model.Group, error)
}

func (f *fakeGroupLibraryGroupRepository) Create(
	group *model.Group,
) error {
	return nil
}

func (f *fakeGroupLibraryGroupRepository) FindByUser(
	userID int64,
) ([]model.Group, error) {
	return nil, nil
}

func (f *fakeGroupLibraryGroupRepository) FindByIDAndUser(
	id int64,
	userID int64,
) (*model.Group, error) {
	if f.findByIDAndUserFunc != nil {
		return f.findByIDAndUserFunc(id, userID)
	}

	return nil, nil
}

func (f *fakeGroupLibraryGroupRepository) Update(
	group *model.Group,
) error {
	return nil
}

func (f *fakeGroupLibraryGroupRepository) Delete(
	id int64,
	userID int64,
) error {
	return nil
}

type fakeGroupLibraryAnimeLibraryRepository struct {
	findByIDAndUserFunc func(
		id int64,
		userID int64,
	) (*model.LibraryAnime, error)
}

func (f *fakeGroupLibraryAnimeLibraryRepository) FindByUserAndAnime(
	userID int64,
	animeID int64,
) (*model.LibraryAnime, error) {
	return nil, nil
}

func (f *fakeGroupLibraryAnimeLibraryRepository) FindByUser(
	userID int64,
) ([]model.LibraryAnime, error) {
	return nil, nil
}

func (f *fakeGroupLibraryAnimeLibraryRepository) FindByIDAndUser(
	id int64,
	userID int64,
) (*model.LibraryAnime, error) {
	if f.findByIDAndUserFunc != nil {
		return f.findByIDAndUserFunc(id, userID)
	}

	return nil, nil
}

func (f *fakeGroupLibraryAnimeLibraryRepository) Create(
	libraryAnime *model.LibraryAnime,
) error {
	return nil
}

func (f *fakeGroupLibraryAnimeLibraryRepository) Update(
	libraryAnime *model.LibraryAnime,
) error {
	return nil
}

func (f *fakeGroupLibraryAnimeLibraryRepository) Delete(
	id int64,
	userID int64,
) error {
	return nil
}

type fakeGroupLibraryAnimeRepository struct {
	findByGroupAndLibraryAnimeFunc func(
		groupID int64,
		libraryAnimeID int64,
	) (*model.GroupLibraryAnime, error)

	findByGroupFunc func(
		groupID int64,
	) ([]model.GroupLibraryAnime, error)

	createFunc func(
		relation *model.GroupLibraryAnime,
	) error

	deleteFunc func(
		groupID int64,
		libraryAnimeID int64,
	) error
}

func (f *fakeGroupLibraryAnimeRepository) FindByGroupAndLibraryAnime(
	groupID int64,
	libraryAnimeID int64,
) (*model.GroupLibraryAnime, error) {
	if f.findByGroupAndLibraryAnimeFunc != nil {
		return f.findByGroupAndLibraryAnimeFunc(
			groupID,
			libraryAnimeID,
		)
	}

	return nil, nil
}

func (f *fakeGroupLibraryAnimeRepository) FindByGroup(
	groupID int64,
) ([]model.GroupLibraryAnime, error) {
	if f.findByGroupFunc != nil {
		return f.findByGroupFunc(groupID)
	}

	return nil, nil
}

func (f *fakeGroupLibraryAnimeRepository) Create(
	relation *model.GroupLibraryAnime,
) error {
	if f.createFunc != nil {
		return f.createFunc(relation)
	}

	return nil
}

func (f *fakeGroupLibraryAnimeRepository) Delete(
	groupID int64,
	libraryAnimeID int64,
) error {
	if f.deleteFunc != nil {
		return f.deleteFunc(groupID, libraryAnimeID)
	}

	return nil
}

func TestGroupLibraryAnimeService_AddAnime(t *testing.T) {
	var createdRelation *model.GroupLibraryAnime

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
				Name:   "Favorites",
			}, nil
		},
	}

	libraryRepo := &fakeGroupLibraryAnimeLibraryRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.LibraryAnime, error) {
			return &model.LibraryAnime{
				ID:      id,
				UserID:  userID,
				AnimeID: 100,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		findByGroupAndLibraryAnimeFunc: func(
			groupID int64,
			libraryAnimeID int64,
		) (*model.GroupLibraryAnime, error) {
			if createdRelation != nil {
				return createdRelation, nil
			}

			return nil, repository.ErrGroupLibraryAnimeNotFound
		},
		createFunc: func(
			relation *model.GroupLibraryAnime,
		) error {
			createdRelation = relation
			return nil
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		libraryRepo,
		groupLibraryRepo,
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		1,
		100,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result == nil {
		t.Fatal("expected relation, got nil")
	}

	if result.GroupID != 1 {
		t.Errorf(
			"expected group ID 1, got %d",
			result.GroupID,
		)
	}

	if result.LibraryAnimeID != 100 {
		t.Errorf(
			"expected library anime ID 100, got %d",
			result.LibraryAnimeID,
		)
	}
}

func TestGroupLibraryAnimeService_AddAnime_GroupNotFound(t *testing.T) {
	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return nil, repository.ErrGroupNotFound
		},
	}

	libraryRepo := &fakeGroupLibraryAnimeLibraryRepository{}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		libraryRepo,
		groupLibraryRepo,
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		999,
		100,
	)

	if !errors.Is(err, ErrGroupNotFoundForUser) {
		t.Fatalf(
			"expected ErrGroupNotFoundForUser, got %v",
			err,
		)
	}

	if result != nil {
		t.Fatalf(
			"expected nil result, got %+v",
			result,
		)
	}
}

func TestGroupLibraryAnimeService_AddAnime_GroupRepositoryError(
	t *testing.T,
) {
	expectedErr := errors.New("database error")

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return nil, expectedErr
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		&fakeGroupLibraryAnimeRepository{},
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		1,
		100,
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

func TestGroupLibraryAnimeService_AddAnime_LibraryAnimeNotFound(
	t *testing.T,
) {
	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	libraryRepo := &fakeGroupLibraryAnimeLibraryRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.LibraryAnime, error) {
			return nil, repository.ErrLibraryAnimeNotFound
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		libraryRepo,
		&fakeGroupLibraryAnimeRepository{},
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		1,
		999,
	)

	if !errors.Is(err, ErrLibraryAnimeNotFound) {
		t.Fatalf(
			"expected ErrLibraryAnimeNotFound, got %v",
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupLibraryAnimeService_AddAnime_LibraryRepositoryError(
	t *testing.T,
) {
	expectedErr := errors.New("database error")

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	libraryRepo := &fakeGroupLibraryAnimeLibraryRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.LibraryAnime, error) {
			return nil, expectedErr
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		libraryRepo,
		&fakeGroupLibraryAnimeRepository{},
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		1,
		100,
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

func TestGroupLibraryAnimeService_AddAnime_AlreadyExists(
	t *testing.T,
) {
	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	libraryRepo := &fakeGroupLibraryAnimeLibraryRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.LibraryAnime, error) {
			return &model.LibraryAnime{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		findByGroupAndLibraryAnimeFunc: func(
			groupID int64,
			libraryAnimeID int64,
		) (*model.GroupLibraryAnime, error) {
			return &model.GroupLibraryAnime{
				GroupID:        groupID,
				LibraryAnimeID: libraryAnimeID,
			}, nil
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		libraryRepo,
		groupLibraryRepo,
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		1,
		100,
	)

	if !errors.Is(err, ErrGroupLibraryAnimeExists) {
		t.Fatalf(
			"expected ErrGroupLibraryAnimeExists, got %v",
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupLibraryAnimeService_AddAnime_FindRelationError(
	t *testing.T,
) {
	expectedErr := errors.New("relation lookup failed")

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	libraryRepo := &fakeGroupLibraryAnimeLibraryRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.LibraryAnime, error) {
			return &model.LibraryAnime{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		findByGroupAndLibraryAnimeFunc: func(
			groupID int64,
			libraryAnimeID int64,
		) (*model.GroupLibraryAnime, error) {
			return nil, expectedErr
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		libraryRepo,
		groupLibraryRepo,
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		1,
		100,
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

func TestGroupLibraryAnimeService_AddAnime_CreateError(
	t *testing.T,
) {
	expectedErr := errors.New("create failed")

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	libraryRepo := &fakeGroupLibraryAnimeLibraryRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.LibraryAnime, error) {
			return &model.LibraryAnime{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		findByGroupAndLibraryAnimeFunc: func(
			groupID int64,
			libraryAnimeID int64,
		) (*model.GroupLibraryAnime, error) {
			return nil, repository.ErrGroupLibraryAnimeNotFound
		},
		createFunc: func(
			relation *model.GroupLibraryAnime,
		) error {
			return expectedErr
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		libraryRepo,
		groupLibraryRepo,
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		1,
		100,
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

func TestGroupLibraryAnimeService_GetAnimes(t *testing.T) {
	expectedRelations := []model.GroupLibraryAnime{
		{
			GroupID:        1,
			LibraryAnimeID: 100,
		},
		{
			GroupID:        1,
			LibraryAnimeID: 101,
		},
	}

	var receivedGroupID int64

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		findByGroupFunc: func(
			groupID int64,
		) ([]model.GroupLibraryAnime, error) {
			receivedGroupID = groupID
			return expectedRelations, nil
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		groupLibraryRepo,
	)

	result, err := service.GetAnimes(
		context.Background(),
		10,
		1,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if receivedGroupID != 1 {
		t.Errorf(
			"expected group ID 1, got %d",
			receivedGroupID,
		)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 relations, got %d",
			len(result),
		)
	}
}

func TestGroupLibraryAnimeService_GetAnimes_GroupNotFound(
	t *testing.T,
) {
	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return nil, repository.ErrGroupNotFound
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		&fakeGroupLibraryAnimeRepository{},
	)

	result, err := service.GetAnimes(
		context.Background(),
		10,
		999,
	)

	if !errors.Is(err, ErrGroupNotFoundForUser) {
		t.Fatalf(
			"expected ErrGroupNotFoundForUser, got %v",
			err,
		)
	}

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupLibraryAnimeService_GetAnimes_GroupRepositoryError(
	t *testing.T,
) {
	expectedErr := errors.New("database error")

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return nil, expectedErr
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		&fakeGroupLibraryAnimeRepository{},
	)

	result, err := service.GetAnimes(
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

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupLibraryAnimeService_GetAnimes_RepositoryError(
	t *testing.T,
) {
	expectedErr := errors.New("database error")

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		findByGroupFunc: func(
			groupID int64,
		) ([]model.GroupLibraryAnime, error) {
			return nil, expectedErr
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		groupLibraryRepo,
	)

	result, err := service.GetAnimes(
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

	if result != nil {
		t.Fatalf("expected nil result, got %+v", result)
	}
}

func TestGroupLibraryAnimeService_DeleteAnime(t *testing.T) {
	var receivedGroupID int64
	var receivedLibraryAnimeID int64

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		deleteFunc: func(
			groupID int64,
			libraryAnimeID int64,
		) error {
			receivedGroupID = groupID
			receivedLibraryAnimeID = libraryAnimeID
			return nil
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		groupLibraryRepo,
	)

	err := service.DeleteAnime(
		context.Background(),
		10,
		1,
		100,
	)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if receivedGroupID != 1 {
		t.Errorf(
			"expected group ID 1, got %d",
			receivedGroupID,
		)
	}

	if receivedLibraryAnimeID != 100 {
		t.Errorf(
			"expected library anime ID 100, got %d",
			receivedLibraryAnimeID,
		)
	}
}

func TestGroupLibraryAnimeService_DeleteAnime_GroupNotFound(
	t *testing.T,
) {
	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return nil, repository.ErrGroupNotFound
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		&fakeGroupLibraryAnimeRepository{},
	)

	err := service.DeleteAnime(
		context.Background(),
		10,
		999,
		100,
	)

	if !errors.Is(err, ErrGroupNotFoundForUser) {
		t.Fatalf(
			"expected ErrGroupNotFoundForUser, got %v",
			err,
		)
	}
}

func TestGroupLibraryAnimeService_DeleteAnime_RepositoryError(
	t *testing.T,
) {
	expectedErr := errors.New("delete failed")

	groupRepo := &fakeGroupLibraryGroupRepository{
		findByIDAndUserFunc: func(
			id int64,
			userID int64,
		) (*model.Group, error) {
			return &model.Group{
				ID:     id,
				UserID: userID,
			}, nil
		},
	}

	groupLibraryRepo := &fakeGroupLibraryAnimeRepository{
		deleteFunc: func(
			groupID int64,
			libraryAnimeID int64,
		) error {
			return expectedErr
		},
	}

	service := NewGroupLibraryAnimeService(
		groupRepo,
		&fakeGroupLibraryAnimeLibraryRepository{},
		groupLibraryRepo,
	)

	err := service.DeleteAnime(
		context.Background(),
		10,
		1,
		100,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}
