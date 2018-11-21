package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
)

type Upload struct {
	gorm.Model `json:"-"`
	Userid string `json:"userid"`
	Url    string `json:"url"`
}

type Uploads struct {
	Uploads []Upload
}

func (upload *Upload) CreateRecord(db *gorm.DB, w http.ResponseWriter) bool {
	if err := db.Create(&upload).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, "Can't make record")
		panic("Can't make record")
		return false
	} else {
		return true
	}
}

func (uploads *Uploads) GetPhotosByUserId(db *gorm.DB, w http.ResponseWriter, userid string) bool {
	if err := db.Find(&uploads.Uploads, "userid = ?", userid).Error; err != nil{
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
		return false
	} else {
		return true
	}
}