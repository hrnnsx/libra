package model

import "time"

type Group struct {
	ID          int64     `gorm:"primaryKey" json:"id"`
	UserID      int64     `gorm:"column:user_id;not null" json:"user_id"`
	Name        string    `gorm:"type:varchar(100);not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
