package faceidentification

import (
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	"net/http"
)

func FaceIdentification(db *gorm.DB, w http.ResponseWriter, url string) {
	// face identification by Azure API
	// 今回はユーザーが全員写っていると想定する。
	var users models.Users
	users.GetAllUsers(db, w)
	for _, user := range users.Users {
		download := models.Download{Userid: user.Userid, Url: url}
		download.CreateRecord(db, w)
	}
}
