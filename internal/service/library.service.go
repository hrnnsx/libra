package service

import (
	"context"
	"errors"
	"strconv"

	"github.com/hrnnsx/libra/external/anilist"
	"github.com/hrnnsx/libra/internal/model"
	"github.com/hrnnsx/libra/internal/repository"
)

var (
	ErrLibraryAnimeExists = errors.New("anime already exists in library")
)

type LibraryService interface {
	AddAnime(
		ctx context.Context,
		userID int64,
		externalID string,
	) (*model.LibraryAnime, error)

	GetLibrary(
		ctx context.Context,
		userID int64,
	) ([]model.LibraryAnime, error)
}

type libraryService struct {
	animeRepository        repository.AnimeRepository
	libraryAnimeRepository repository.LibraryAnimeRepository
	anilistClient          *anilist.Client
}

func NewLibraryService(
	animeRepository repository.AnimeRepository,
	libraryAnimeRepository repository.LibraryAnimeRepository,
	anilistClient *anilist.Client,
) LibraryService {
	return &libraryService{
		animeRepository:        animeRepository,
		libraryAnimeRepository: libraryAnimeRepository,
		anilistClient:          anilistClient,
	}
}

func (s *libraryService) AddAnime(
	ctx context.Context,
	userID int64,
	externalID string,
) (*model.LibraryAnime, error) {

	id, err := strconv.Atoi(externalID)
	if err != nil {
		return nil, ErrAnimeNotFound
	}

	anime, err := s.animeRepository.FindByExternalID(
		externalID,
	)

	if err != nil &&
		!errors.Is(err, repository.ErrAnimeNotFound) {
		return nil, err
	}

	if errors.Is(err, repository.ErrAnimeNotFound) {
		externalAnime, err := s.anilistClient.GetAnime(
			ctx,
			id,
		)

		if err != nil {
			return nil, ErrAnimeNotFound
		}

		synopsis := externalAnime.Description
		imageURL := externalAnime.CoverImage.Large

		anime = &model.Anime{
			ExternalID: externalID,
			Title:      getAnimeTitle(externalAnime),
			Synopsis:   &synopsis,
			ImageURL:   &imageURL,
		}

		if err := s.animeRepository.Create(anime); err != nil {
			return nil, err
		}
	}

	_, err = s.libraryAnimeRepository.FindByUserAndAnime(
		userID,
		anime.ID,
	)

	if err == nil {
		return nil, ErrLibraryAnimeExists
	}

	if !errors.Is(
		err,
		repository.ErrLibraryAnimeNotFound,
	) {
		return nil, err
	}

	libraryAnime := &model.LibraryAnime{
		UserID:         userID,
		AnimeID:        anime.ID,
		Status:         "PLAN_TO_WATCH",
		CurrentEpisode: 0,
	}

	if err := s.libraryAnimeRepository.Create(
		libraryAnime,
	); err != nil {
		return nil, err
	}

	return libraryAnime, nil
}

func (s *libraryService) GetLibrary(
	ctx context.Context,
	userID int64,
) ([]model.LibraryAnime, error) {

	return s.libraryAnimeRepository.FindByUser(userID)
}

func getAnimeTitle(anime *anilist.Anime) string {
	if anime.Title.English != "" {
		return anime.Title.English
	}

	if anime.Title.Romaji != "" {
		return anime.Title.Romaji
	}

	return anime.Title.Native
}
