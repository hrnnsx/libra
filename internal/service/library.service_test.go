package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
	"gorm.io/gorm"
)

type fakeAnimeRepository struct {
	findResult *model.Anime
	findErr    error
	createErr  error
	created    *model.Anime
}

func (f *fakeAnimeRepository) FindByExternalID(
	externalID string,
) (*model.Anime, error) {
	return f.findResult, f.findErr
}

func (f *fakeAnimeRepository) Create(
	anime *model.Anime,
) error {
	if f.createErr != nil {
		return f.createErr
	}

	f.created = anime

	return nil
}

type fakeLibraryAnimeRepository struct {
	findByUserAndAnimeResult *model.LibraryAnime
	findByUserAndAnimeErr    error

	findByUserResult []model.LibraryAnime
	findByUserErr    error

	findByIDAndUserResult *model.LibraryAnime
	findByIDAndUserErr    error

	createErr error
	created   *model.LibraryAnime

	updateErr error
	updated   *model.LibraryAnime

	deleteErr   error
	deletedID   int64
	deletedUser int64
}

func (f *fakeLibraryAnimeRepository) FindByUserAndAnime(
	userID int64,
	animeID int64,
) (*model.LibraryAnime, error) {
	return f.findByUserAndAnimeResult, f.findByUserAndAnimeErr
}

func (f *fakeLibraryAnimeRepository) FindByUser(
	userID int64,
) ([]model.LibraryAnime, error) {
	return f.findByUserResult, f.findByUserErr
}

func (f *fakeLibraryAnimeRepository) FindByIDAndUser(
	id int64,
	userID int64,
) (*model.LibraryAnime, error) {
	return f.findByIDAndUserResult, f.findByIDAndUserErr
}

func (f *fakeLibraryAnimeRepository) Create(
	libraryAnime *model.LibraryAnime,
) error {
	if f.createErr != nil {
		return f.createErr
	}

	f.created = libraryAnime

	return nil
}

func (f *fakeLibraryAnimeRepository) Update(
	libraryAnime *model.LibraryAnime,
) error {
	if f.updateErr != nil {
		return f.updateErr
	}

	f.updated = libraryAnime

	return nil
}

func (f *fakeLibraryAnimeRepository) Delete(
	id int64,
	userID int64,
) error {
	f.deletedID = id
	f.deletedUser = userID

	return f.deleteErr
}

func TestLibraryService_AddAnime_ExistingAnime(t *testing.T) {
	animeRepository := &fakeAnimeRepository{
		findResult: &model.Anime{
			ID:         1,
			ExternalID: "20464",
			Title:      "The Irregular at Magic High School",
		},
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByUserAndAnimeErr: repository.ErrLibraryAnimeNotFound,
	}

	service := NewLibraryService(
		animeRepository,
		libraryRepository,
		nil,
	)

	result, err := service.AddAnime(
		context.Background(),
		10,
		"20464",
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if result == nil {
		t.Fatal("expected library anime, got nil")
	}

	if result.UserID != 10 {
		t.Errorf(
			"expected user_id 10, got %d",
			result.UserID,
		)
	}

	if result.AnimeID != 1 {
		t.Errorf(
			"expected anime_id 1, got %d",
			result.AnimeID,
		)
	}

	if result.Status != "PLAN_TO_WATCH" {
		t.Errorf(
			"expected status PLAN_TO_WATCH, got %s",
			result.Status,
		)
	}

	if result.CurrentEpisode != 0 {
		t.Errorf(
			"expected current_episode 0, got %d",
			result.CurrentEpisode,
		)
	}

	if libraryRepository.created == nil {
		t.Fatal("expected repository Create to be called")
	}
}

func TestLibraryService_AddAnime_AlreadyExists(t *testing.T) {
	animeRepository := &fakeAnimeRepository{
		findResult: &model.Anime{
			ID:         1,
			ExternalID: "20464",
			Title:      "The Irregular at Magic High School",
		},
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByUserAndAnimeResult: &model.LibraryAnime{
			ID:      100,
			UserID:  10,
			AnimeID: 1,
		},
	}

	service := NewLibraryService(
		animeRepository,
		libraryRepository,
		nil,
	)

	_, err := service.AddAnime(
		context.Background(),
		10,
		"20464",
	)

	if !errors.Is(err, ErrLibraryAnimeExists) {
		t.Fatalf(
			"expected ErrLibraryAnimeExists, got %v",
			err,
		)
	}
}

func TestLibraryService_AddAnime_InvalidExternalID(t *testing.T) {
	service := NewLibraryService(
		&fakeAnimeRepository{},
		&fakeLibraryAnimeRepository{},
		nil,
	)

	_, err := service.AddAnime(
		context.Background(),
		10,
		"abc",
	)

	if !errors.Is(err, ErrAnimeNotFound) {
		t.Fatalf(
			"expected ErrAnimeNotFound, got %v",
			err,
		)
	}
}

func TestLibraryService_AddAnime_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	animeRepository := &fakeAnimeRepository{
		findErr: expectedErr,
	}

	service := NewLibraryService(
		animeRepository,
		&fakeLibraryAnimeRepository{},
		nil,
	)

	_, err := service.AddAnime(
		context.Background(),
		10,
		"20464",
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}

func TestLibraryService_GetLibrary(t *testing.T) {
	expected := []model.LibraryAnime{
		{
			ID:             1,
			UserID:         10,
			AnimeID:        100,
			Status:         "WATCHING",
			CurrentEpisode: 5,
		},
		{
			ID:             2,
			UserID:         10,
			AnimeID:        200,
			Status:         "COMPLETED",
			CurrentEpisode: 12,
		},
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByUserResult: expected,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	result, err := service.GetLibrary(
		context.Background(),
		10,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if len(result) != 2 {
		t.Fatalf(
			"expected 2 anime, got %d",
			len(result),
		)
	}

	if result[0].Status != "WATCHING" {
		t.Errorf(
			"expected WATCHING, got %s",
			result[0].Status,
		)
	}
}

func TestLibraryService_GetLibrary_Error(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	libraryRepository := &fakeLibraryAnimeRepository{
		findByUserErr: expectedErr,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	_, err := service.GetLibrary(
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
}

func TestLibraryService_GetAnime(t *testing.T) {
	expected := &model.LibraryAnime{
		ID:             1,
		UserID:         10,
		AnimeID:        100,
		Status:         "WATCHING",
		CurrentEpisode: 5,
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserResult: expected,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	result, err := service.GetAnime(
		context.Background(),
		10,
		1,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if result.ID != 1 {
		t.Errorf(
			"expected id 1, got %d",
			result.ID,
		)
	}

	if result.UserID != 10 {
		t.Errorf(
			"expected user_id 10, got %d",
			result.UserID,
		)
	}
}

func TestLibraryService_GetAnime_NotFound(t *testing.T) {
	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserErr: repository.ErrLibraryAnimeNotFound,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	_, err := service.GetAnime(
		context.Background(),
		10,
		999,
	)

	if !errors.Is(err, repository.ErrLibraryAnimeNotFound) {
		t.Fatalf(
			"expected ErrLibraryAnimeNotFound, got %v",
			err,
		)
	}
}

func TestLibraryService_UpdateAnime(t *testing.T) {
	startedAt := "2026-10-01T10:00:00Z"
	completedAt := "2026-10-05T10:00:00Z"

	libraryAnime := &model.LibraryAnime{
		ID:             1,
		UserID:         10,
		AnimeID:        100,
		Status:         "PLAN_TO_WATCH",
		CurrentEpisode: 0,
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserResult: libraryAnime,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	status := "COMPLETED"
	currentEpisode := 12
	rating := 9.5
	notes := "Great anime"

	result, err := service.UpdateAnime(
		context.Background(),
		10,
		1,
		&status,
		&currentEpisode,
		&rating,
		&notes,
		&startedAt,
		&completedAt,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if libraryRepository.updated == nil {
		t.Fatal("expected repository Update to be called")
	}

	updated := libraryRepository.updated

	if updated.Status != "COMPLETED" {
		t.Errorf(
			"expected status COMPLETED, got %s",
			updated.Status,
		)
	}

	if updated.CurrentEpisode != 12 {
		t.Errorf(
			"expected current_episode 12, got %d",
			updated.CurrentEpisode,
		)
	}

	if updated.Rating == nil || *updated.Rating != 9.5 {
		t.Errorf("expected rating 9.5")
	}

	if updated.Notes == nil || *updated.Notes != "Great anime" {
		t.Errorf("expected notes 'Great anime'")
	}

	if updated.StartedAt == nil {
		t.Fatal("expected started_at")
	}

	if updated.CompletedAt == nil {
		t.Fatal("expected completed_at")
	}

	if result != libraryAnime {
		t.Error("expected updated library anime to be returned")
	}
}

func TestLibraryService_UpdateAnime_InvalidStartedAt(t *testing.T) {
	libraryAnime := &model.LibraryAnime{
		ID:     1,
		UserID: 10,
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserResult: libraryAnime,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	invalidTime := "not-a-date"

	_, err := service.UpdateAnime(
		context.Background(),
		10,
		1,
		nil,
		nil,
		nil,
		nil,
		&invalidTime,
		nil,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if libraryRepository.updated != nil {
		t.Error("repository Update should not be called")
	}
}

func TestLibraryService_UpdateAnime_InvalidCompletedAt(t *testing.T) {
	libraryAnime := &model.LibraryAnime{
		ID:     1,
		UserID: 10,
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserResult: libraryAnime,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	invalidTime := "not-a-date"

	_, err := service.UpdateAnime(
		context.Background(),
		10,
		1,
		nil,
		nil,
		nil,
		nil,
		nil,
		&invalidTime,
	)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if libraryRepository.updated != nil {
		t.Error("repository Update should not be called")
	}
}

func TestLibraryService_UpdateAnime_RepositoryError(t *testing.T) {
	expectedErr := errors.New("database unavailable")

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserErr: expectedErr,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	status := "WATCHING"

	_, err := service.UpdateAnime(
		context.Background(),
		10,
		1,
		&status,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}

func TestLibraryService_DeleteAnime(t *testing.T) {
	libraryRepository := &fakeLibraryAnimeRepository{}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	err := service.DeleteAnime(
		context.Background(),
		10,
		1,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if libraryRepository.deletedID != 1 {
		t.Errorf(
			"expected deleted id 1, got %d",
			libraryRepository.deletedID,
		)
	}

	if libraryRepository.deletedUser != 10 {
		t.Errorf(
			"expected deleted user 10, got %d",
			libraryRepository.deletedUser,
		)
	}
}

func TestLibraryService_DeleteAnime_NotFound(t *testing.T) {
	libraryRepository := &fakeLibraryAnimeRepository{
		deleteErr: repository.ErrLibraryAnimeNotFound,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	err := service.DeleteAnime(
		context.Background(),
		10,
		999,
	)

	if !errors.Is(err, repository.ErrLibraryAnimeNotFound) {
		t.Fatalf(
			"expected ErrLibraryAnimeNotFound, got %v",
			err,
		)
	}
}

func TestLibraryService_UpdateAnime_PartialUpdate(t *testing.T) {
	libraryAnime := &model.LibraryAnime{
		ID:             1,
		UserID:         10,
		AnimeID:        100,
		Status:         "WATCHING",
		CurrentEpisode: 5,
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserResult: libraryAnime,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	newEpisode := 6

	result, err := service.UpdateAnime(
		context.Background(),
		10,
		1,
		nil,
		&newEpisode,
		nil,
		nil,
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if result.CurrentEpisode != 6 {
		t.Errorf(
			"expected current_episode 6, got %d",
			result.CurrentEpisode,
		)
	}

	if result.Status != "WATCHING" {
		t.Errorf(
			"expected status to remain WATCHING, got %s",
			result.Status,
		)
	}
}

func TestLibraryService_UpdateAnime_RepositoryUpdateError(t *testing.T) {
	expectedErr := errors.New("update failed")

	libraryAnime := &model.LibraryAnime{
		ID:     1,
		UserID: 10,
		Status: "WATCHING",
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserResult: libraryAnime,
		updateErr:             expectedErr,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	status := "COMPLETED"

	_, err := service.UpdateAnime(
		context.Background(),
		10,
		1,
		&status,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	if !errors.Is(err, expectedErr) {
		t.Fatalf(
			"expected error %v, got %v",
			expectedErr,
			err,
		)
	}
}

func TestLibraryService_DeleteAnime_RepositoryError(t *testing.T) {
	expectedErr := errors.New("delete failed")

	libraryRepository := &fakeLibraryAnimeRepository{
		deleteErr: expectedErr,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	err := service.DeleteAnime(
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

func TestLibraryService_UpdateAnime_PreservesExistingTimestamps(t *testing.T) {
	started := time.Date(
		2026,
		10,
		1,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	completed := time.Date(
		2026,
		10,
		5,
		10,
		0,
		0,
		0,
		time.UTC,
	)

	libraryAnime := &model.LibraryAnime{
		ID:          1,
		UserID:      10,
		StartedAt:   &started,
		CompletedAt: &completed,
		Status:      "WATCHING",
	}

	libraryRepository := &fakeLibraryAnimeRepository{
		findByIDAndUserResult: libraryAnime,
	}

	service := NewLibraryService(
		&fakeAnimeRepository{},
		libraryRepository,
		nil,
	)

	status := "COMPLETED"

	_, err := service.UpdateAnime(
		context.Background(),
		10,
		1,
		&status,
		nil,
		nil,
		nil,
		nil,
		nil,
	)

	if err != nil {
		t.Fatalf(
			"expected no error, got %v",
			err,
		)
	}

	if libraryRepository.updated.StartedAt == nil {
		t.Fatal("expected started_at to remain")
	}

	if libraryRepository.updated.CompletedAt == nil {
		t.Fatal("expected completed_at to remain")
	}

	if !libraryRepository.updated.StartedAt.Equal(started) {
		t.Error("started_at was unexpectedly changed")
	}

	if !libraryRepository.updated.CompletedAt.Equal(completed) {
		t.Error("completed_at was unexpectedly changed")
	}
}

var _ repository.AnimeRepository = (*fakeAnimeRepository)(nil)
var _ repository.LibraryAnimeRepository = (*fakeLibraryAnimeRepository)(nil)

var _ = gorm.ErrRecordNotFound
