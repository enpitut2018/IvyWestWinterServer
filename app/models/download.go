package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type Download struct {
	gorm.Model `json:"-"`
	UserID     string `json:"userid"`
	URL        string `json:"url"`
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

func (downloads *Downloads) GetDownloadsByUserID(db *gorm.DB, userID string) error {
	if err := db.Find(&downloads.Downloads, "user_id = ?", userID).Error; err != nil {
		return err
	}
	return nil
}
