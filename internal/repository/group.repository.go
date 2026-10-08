package repository

import (
	"errors"

	"github.com/hrnnsx/libra/internal/model"
	"gorm.io/gorm"
)

var ErrGroupNotFound = errors.New("group not found")

type GroupRepository interface {
	Create(group *model.Group) error

	FindByUser(userID int64) ([]model.Group, error)

	FindByIDAndUser(
		id int64,
		userID int64,
	) (*model.Group, error)

	Update(group *model.Group) error

	Delete(
		id int64,
		userID int64,
	) error
}

type groupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) GroupRepository {
	return &groupRepository{
		db: db,
	}
}

func (r *groupRepository) Create(
	group *model.Group,
) error {
	return r.db.Create(group).Error
}

func (r *groupRepository) FindByUser(
	userID int64,
) ([]model.Group, error) {
	var groups []model.Group

	err := r.db.
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Find(&groups).
		Error

	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (r *groupRepository) FindByIDAndUser(
	id int64,
	userID int64,
) (*model.Group, error) {
	var group model.Group

	err := r.db.
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		First(&group).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrGroupNotFound
	}

	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *groupRepository) Update(
	group *model.Group,
) error {
	return r.db.
		Model(&model.Group{}).
		Where(
			"id = ? AND user_id = ?",
			group.ID,
			group.UserID,
		).
		Updates(map[string]interface{}{
			"name":        group.Name,
			"description": group.Description,
		}).
		Error
}

func (r *groupRepository) Delete(
	id int64,
	userID int64,
) error {
	result := r.db.
		Where(
			"id = ? AND user_id = ?",
			id,
			userID,
		).
		Delete(&model.Group{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrGroupNotFound
	}

	return nil
}
