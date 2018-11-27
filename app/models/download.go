package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"net/http"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
)

type Download struct {
	gorm.Model `json:"-"`
	UserID   string `json:"userid"`
	URL string `json:"url"`
}

type Downloads struct {
	Downloads []Download
}

func (download *Download) CreateRecord(db *gorm.DB, w http.ResponseWriter) bool {
	if err := db.Create(&download).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, "Can't make record.")
		panic("Can't make record.")
	}
	return true
}

func (downloads *Downloads) GetDownloadsByUserID(db *gorm.DB, w http.ResponseWriter, userID string) bool {
	if err := db.Find(&downloads.Downloads, "userID = ?", userID).Error; err != nil{
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}
	return true
}