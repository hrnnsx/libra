package repository

import (
	"errors"

	"github.com/hrnnsx/libra/internal/model"
	"gorm.io/gorm"
)

var (
	ErrGroupLibraryAnimeNotFound = errors.New(
		"group library anime not found",
	)

	ErrGroupLibraryAnimeExists = errors.New(
		"anime already exists in group",
	)
)

type GroupLibraryAnimeRepository interface {
	FindByGroupAndLibraryAnime(
		groupID int64,
		libraryAnimeID int64,
	) (*model.GroupLibraryAnime, error)

	FindByGroup(
		groupID int64,
	) ([]model.GroupLibraryAnime, error)

	Create(
		groupLibraryAnime *model.GroupLibraryAnime,
	) error

	Delete(
		groupID int64,
		libraryAnimeID int64,
	) error
}

type groupLibraryAnimeRepository struct {
	db *gorm.DB
}

func NewGroupLibraryAnimeRepository(
	db *gorm.DB,
) GroupLibraryAnimeRepository {
	return &groupLibraryAnimeRepository{
		db: db,
	}
}

func (r *groupLibraryAnimeRepository) FindByGroupAndLibraryAnime(
	groupID int64,
	libraryAnimeID int64,
) (*model.GroupLibraryAnime, error) {
	var relation model.GroupLibraryAnime

	err := r.db.
		Where(
			"group_id = ? AND library_anime_id = ?",
			groupID,
			libraryAnimeID,
		).
		First(&relation).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGroupLibraryAnimeNotFound
	}

	if err != nil {
		return nil, err
	}

	return &relation, nil
}

func (r *groupLibraryAnimeRepository) FindByGroup(
	groupID int64,
) ([]model.GroupLibraryAnime, error) {
	var relations []model.GroupLibraryAnime

	err := r.db.
		Preload("LibraryAnime").
		Preload("LibraryAnime.Anime").
		Where("group_id = ?", groupID).
		Order("created_at DESC").
		Find(&relations).
		Error

	if err != nil {
		return nil, err
	}

	return relations, nil
}

func (r *groupLibraryAnimeRepository) Create(
	groupLibraryAnime *model.GroupLibraryAnime,
) error {
	return r.db.Create(groupLibraryAnime).Error
}

func (r *groupLibraryAnimeRepository) Delete(
	groupID int64,
	libraryAnimeID int64,
) error {
	result := r.db.
		Where(
			"group_id = ? AND library_anime_id = ?",
			groupID,
			libraryAnimeID,
		).
		Delete(&model.GroupLibraryAnime{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrGroupLibraryAnimeNotFound
	}

	return nil
}
