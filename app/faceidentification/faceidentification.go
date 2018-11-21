package faceidentification

import (
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
)

func FaceIdentification(db *gorm.DB, photourl string) {
	// face identification by Azure API
	// 今回はユーザーが全員写っていると想定する。
	var users []models.User
	if err := db.Raw("SELECT * FROM users").Scan(&users).Error; err != nil {
		panic("can't Identification")
	}
	for _, user := range users {
		var download models.Download
		download.Userid = user.Userid
		download.PhotoUrl = photourl
		if err := db.Create(&download).Error; err != nil {
			panic("Can't make record")
		}
	}
}
