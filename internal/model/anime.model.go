package model

import "time"

type Anime struct {
	ID         int64     `gorm:"primaryKey" json:"id"`
	ExternalID string    `gorm:"column:external_id;uniqueIndex;not null" json:"external_id"`
	Title      string    `gorm:"type:varchar(255);not null" json:"title"`
	Synopsis   *string   `gorm:"type:text" json:"synopsis,omitempty"`
	ImageURL   *string   `gorm:"type:varchar(500)" json:"image_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type LibraryAnime struct {
	ID             int64      `gorm:"primaryKey" json:"id"`
	UserID         int64      `gorm:"column:user_id;not null" json:"user_id"`
	AnimeID        int64      `gorm:"column:anime_id;not null" json:"anime_id"`
	Status         string     `gorm:"type:varchar(30);not null" json:"status"`
	CurrentEpisode int        `gorm:"column:current_episode;not null" json:"current_episode"`
	Rating         *float64   `gorm:"type:decimal(3,1)" json:"rating,omitempty"`
	Notes          *string    `gorm:"type:text" json:"notes,omitempty"`
	StartedAt      *time.Time `gorm:"column:started_at" json:"started_at,omitempty"`
	CompletedAt    *time.Time `gorm:"column:completed_at" json:"completed_at,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

func (Anime) TableName() string {
	return "anime"
}

func (LibraryAnime) TableName() string {
	return "library_anime"
}
