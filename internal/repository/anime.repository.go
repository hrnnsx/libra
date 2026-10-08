package repository

import (
	"errors"

	"github.com/hrnnsx/libra/internal/model"
	"gorm.io/gorm"
)

var ErrAnimeNotFound = errors.New("anime not found")

type AnimeRepository interface {
	FindByExternalID(externalID string) (*model.Anime, error)
	Create(anime *model.Anime) error
}

type animeRepository struct {
	db *gorm.DB
}

func NewAnimeRepository(db *gorm.DB) AnimeRepository {
	return &animeRepository{
		db: db,
	}
}

func (r *animeRepository) FindByExternalID(
	externalID string,
) (*model.Anime, error) {
	var anime model.Anime

	err := r.db.
		Where("external_id = ?", externalID).
		First(&anime).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAnimeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &anime, nil
}

func (r *animeRepository) Create(
	anime *model.Anime,
) error {
	return r.db.Create(anime).Error
}
