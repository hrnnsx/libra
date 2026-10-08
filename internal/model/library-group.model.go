package model

import "time"

type GroupLibraryAnime struct {
	GroupID        int64     `gorm:"column:group_id;primaryKey" json:"group_id"`
	LibraryAnimeID int64     `gorm:"column:library_anime_id;primaryKey" json:"library_anime_id"`
	CreatedAt      time.Time `json:"created_at"`

	LibraryAnime LibraryAnime `gorm:"foreignKey:LibraryAnimeID" json:"library_anime"`
}

func (GroupLibraryAnime) TableName() string {
	return "group_library_anime"
}
