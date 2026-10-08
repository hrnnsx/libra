package service

import (
	"context"
	"errors"

	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
)

var (
	ErrGroupLibraryAnimeNotFound = errors.New(
		"group library anime not found",
	)

	ErrGroupLibraryAnimeExists = errors.New(
		"anime already exists in group",
	)

	ErrLibraryAnimeNotFound = errors.New(
		"library anime not found",
	)

	ErrGroupNotFoundForUser = errors.New(
		"group not found",
	)
)

type GroupLibraryAnimeService interface {
	AddAnime(
		ctx context.Context,
		userID int64,
		groupID int64,
		libraryAnimeID int64,
	) (*model.GroupLibraryAnime, error)

	GetAnimes(
		ctx context.Context,
		userID int64,
		groupID int64,
	) ([]model.GroupLibraryAnime, error)

	DeleteAnime(
		ctx context.Context,
		userID int64,
		groupID int64,
		libraryAnimeID int64,
	) error
}

type groupLibraryAnimeService struct {
	groupRepository             repository.GroupRepository
	libraryAnimeRepository      repository.LibraryAnimeRepository
	groupLibraryAnimeRepository repository.GroupLibraryAnimeRepository
}

func NewGroupLibraryAnimeService(
	groupRepository repository.GroupRepository,
	libraryAnimeRepository repository.LibraryAnimeRepository,
	groupLibraryAnimeRepository repository.GroupLibraryAnimeRepository,
) GroupLibraryAnimeService {
	return &groupLibraryAnimeService{
		groupRepository:             groupRepository,
		libraryAnimeRepository:      libraryAnimeRepository,
		groupLibraryAnimeRepository: groupLibraryAnimeRepository,
	}
}

func (s *groupLibraryAnimeService) AddAnime(
	ctx context.Context,
	userID int64,
	groupID int64,
	libraryAnimeID int64,
) (*model.GroupLibraryAnime, error) {
	// Pastikan group milik user.
	_, err := s.groupRepository.FindByIDAndUser(
		groupID,
		userID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrGroupNotFound) {
			return nil, ErrGroupNotFoundForUser
		}

		return nil, err
	}

	// Pastikan library anime milik user.
	_, err = s.libraryAnimeRepository.FindByIDAndUser(
		libraryAnimeID,
		userID,
	)

	if err != nil {
		if errors.Is(
			err,
			repository.ErrLibraryAnimeNotFound,
		) {
			return nil, ErrLibraryAnimeNotFound
		}

		return nil, err
	}

	// Pastikan anime belum ada di group.
	_, err = s.groupLibraryAnimeRepository.FindByGroupAndLibraryAnime(
		groupID,
		libraryAnimeID,
	)

	if err == nil {
		return nil, ErrGroupLibraryAnimeExists
	}

	if !errors.Is(
		err,
		repository.ErrGroupLibraryAnimeNotFound,
	) {
		return nil, err
	}

	relation := &model.GroupLibraryAnime{
		GroupID:        groupID,
		LibraryAnimeID: libraryAnimeID,
	}

	if err := s.groupLibraryAnimeRepository.Create(
		relation,
	); err != nil {
		return nil, err
	}

	return s.groupLibraryAnimeRepository.FindByGroupAndLibraryAnime(
		groupID,
		libraryAnimeID,
	)
}

func (s *groupLibraryAnimeService) GetAnimes(
	ctx context.Context,
	userID int64,
	groupID int64,
) ([]model.GroupLibraryAnime, error) {
	// Pastikan group milik user.
	_, err := s.groupRepository.FindByIDAndUser(
		groupID,
		userID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrGroupNotFound) {
			return nil, ErrGroupNotFoundForUser
		}

		return nil, err
	}

	return s.groupLibraryAnimeRepository.FindByGroup(
		groupID,
	)
}

func (s *groupLibraryAnimeService) DeleteAnime(
	ctx context.Context,
	userID int64,
	groupID int64,
	libraryAnimeID int64,
) error {
	// Pastikan group milik user.
	_, err := s.groupRepository.FindByIDAndUser(
		groupID,
		userID,
	)

	if err != nil {
		if errors.Is(err, repository.ErrGroupNotFound) {
			return ErrGroupNotFoundForUser
		}

		return err
	}

	return s.groupLibraryAnimeRepository.Delete(
		groupID,
		libraryAnimeID,
	)
}
