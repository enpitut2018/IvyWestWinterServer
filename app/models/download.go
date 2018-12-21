package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Download struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	PhotoID   uint       `json:"photoid"`
	UserID    string     `json:"userid"`
	URL       string     `json:"url"`
}

type Downloads struct {
	Downloads []Download
}

func (download *Download) CreateRecord(db *gorm.DB) error {
	if err := db.Create(&download).Error; err != nil {
		return err
	}
	return nil
}

func (download *Download) GetDownloadByPhotoID(db *gorm.DB, photoID uint) error {
	db.Find(&download, "photo_id = ?", photoID)
	return nil
}

func (downloads *Downloads) GetDownloadsByUserID(db *gorm.DB, userID string) error {
	if err := db.Find(&downloads.Downloads, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}

func (downloads *Downloads) GetDownloadsByPhotoID(db *gorm.DB, photoID uint) error {
	db.Find(&downloads.Downloads, "photo_id = ?", photoID)
	return nil
}
