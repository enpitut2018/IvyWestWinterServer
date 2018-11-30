package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Upload struct {
	gorm.Model `json:"-"`
	UserID     string `json:"userid"`
	URL        string `json:"url"`
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
	if err := db.Find(&uploads.Uploads, "userid = ?", userid).Error; err != nil {
		return err
	}
	return nil
}
