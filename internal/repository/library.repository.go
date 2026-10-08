package repository

import (
	"errors"

	"github.com/hrnnsx/libra/internal/model"
	"gorm.io/gorm"
)

var ErrLibraryAnimeNotFound = errors.New("library anime not found")
var ErrLibraryAnimeExists = errors.New("anime already exists in library")

type LibraryAnimeRepository interface {
	FindByUserAndAnime(
		userID int64,
		animeID int64,
	) (*model.LibraryAnime, error)

	Create(libraryAnime *model.LibraryAnime) error
}

type libraryAnimeRepository struct {
	db *gorm.DB
}

func NewLibraryAnimeRepository(db *gorm.DB) LibraryAnimeRepository {
	return &libraryAnimeRepository{
		db: db,
	}
}

func (r *libraryAnimeRepository) FindByUserAndAnime(
	userID int64,
	animeID int64,
) (*model.LibraryAnime, error) {
	var libraryAnime model.LibraryAnime

	err := r.db.
		Where("user_id = ? AND anime_id = ?", userID, animeID).
		First(&libraryAnime).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrLibraryAnimeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &libraryAnime, nil
}

func (r *libraryAnimeRepository) Create(
	libraryAnime *model.LibraryAnime,
) error {
	return r.db.Create(libraryAnime).Error
}
