package repository

import (
	"errors"

	"github.com/hrnnsx/libra/internal/model"
	"gorm.io/gorm"
)

var (
	ErrLibraryAnimeNotFound = errors.New("library anime not found")
	ErrLibraryAnimeExists   = errors.New("anime already exists in library")
)

type LibraryAnimeRepository interface {
	FindByUserAndAnime(
		userID int64,
		animeID int64,
	) (*model.LibraryAnime, error)

	FindByUser(
		userID int64,
	) ([]model.LibraryAnime, error)

	FindByIDAndUser(
		id int64,
		userID int64,
	) (*model.LibraryAnime, error)

	Create(
		libraryAnime *model.LibraryAnime,
	) error

	Update(
		libraryAnime *model.LibraryAnime,
	) error

	Delete(
		id int64,
		userID int64,
	) error
}

type libraryAnimeRepository struct {
	db *gorm.DB
}

func NewLibraryAnimeRepository(
	db *gorm.DB,
) LibraryAnimeRepository {
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
		Where(
			"user_id = ? AND anime_id = ?",
			userID,
			animeID,
		).
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

func (r *libraryAnimeRepository) FindByUser(
	userID int64,
) ([]model.LibraryAnime, error) {
	var libraryAnimes []model.LibraryAnime

	err := r.db.
		Preload("Anime").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&libraryAnimes).
		Error

	if err != nil {
		return nil, err
	}

	return libraryAnimes, nil
}

func (r *libraryAnimeRepository) FindByIDAndUser(
	id int64,
	userID int64,
) (*model.LibraryAnime, error) {
	var libraryAnime model.LibraryAnime

	err := r.db.
		Preload("Anime").
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
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

func (r *libraryAnimeRepository) Update(
	libraryAnime *model.LibraryAnime,
) error {
	return r.db.
		Model(&model.LibraryAnime{}).
		Where(
			"id = ? AND user_id = ?",
			libraryAnime.ID,
			libraryAnime.UserID,
		).
		Updates(map[string]interface{}{
			"status":          libraryAnime.Status,
			"current_episode": libraryAnime.CurrentEpisode,
			"rating":          libraryAnime.Rating,
			"notes":           libraryAnime.Notes,
			"started_at":      libraryAnime.StartedAt,
			"completed_at":    libraryAnime.CompletedAt,
		}).Error
}

func (r *libraryAnimeRepository) Delete(
	id int64,
	userID int64,
) error {
	result := r.db.
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		Delete(&model.LibraryAnime{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrLibraryAnimeNotFound
	}

	return nil
}
