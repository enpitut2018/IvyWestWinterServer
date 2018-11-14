package faceidentification

import "fmt"

func FaceIdentification(photourl) {
	db := dbutils.ConnectPostgres()
	defer db.Close()

	// face identification by Azure API
	// 今回はユーザーが全員写っていると想定する。
	var users []dbutils.User
	if err := db.Raw("SELECT * FROM users").Scan(&users).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, "can't Identification")
		panic("can't Identification")
	}
	for user := range users {
		var download dbutils.Download
		download.UserId = user.UserId
		download.Url = photourl
		if err = db.Create(&download).Error; err != nil {
			httputils.RespondError(w, http.StatusInternalServerError, "Can't make record")
			panic("Can't make record")
		}
	}
}
