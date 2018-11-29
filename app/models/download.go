package models

import (
	"net/http"

	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	l "github.com/sirupsen/logrus"
)

type Download struct {
	gorm.Model `json:"-"`
	UserID     string `json:"userid"`
	URL        string `json:"url"`
}

type Downloads struct {
	Downloads []Download
}

func (download *Download) CreateRecord(db *gorm.DB, w http.ResponseWriter) bool {
	if err := db.Create(&download).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, "Can't make record.")
		l.Errorf("Can't make record.")
	}
	return true
}

func (downloads *Downloads) GetDownloadsByUserID(db *gorm.DB, w http.ResponseWriter, userID string) bool {
	if err := db.Find(&downloads.Downloads, "user_id = ?", userID).Error; err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		l.Errorf(err.Error())
	}
	return true
}
