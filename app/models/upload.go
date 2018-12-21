package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Upload struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	UserID    string     `json:"userid"`
	URL       string     `json:"url"`
}

type Uploads struct {
	Uploads []Upload
}

func (upload *Upload) CreateRecord(db *gorm.DB) error {
	if err := db.Create(&upload).Error; err != nil {
		return err
	}
	return nil
}

func (uploads *Uploads) GetPhotosByUserID(db *gorm.DB, userid string) error {
	if err := db.Find(&uploads.Uploads, "user_id = ?", userid).Error; err != nil {
		return err
	}
	return nil
}

func (upload *Upload) GetPhotoByPhotoID(db *gorm.DB, photoid string) error {
	if err := db.Raw("SELECT * FROM uploads WHERE id = ?", photoid).Scan(&upload).Error; err != nil {
		return err
	}
	return nil
}
