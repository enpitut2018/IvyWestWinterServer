package userface

import (
	"net/http"
)

type Source struct {
	Source string
}

func UploadUserFace(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	var source Source
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&source); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	db := dbutils.ConnectPostgres()
	defer db.Close()

	data, err := base64.StdEncoding.DecodeString(source.Source)
	if err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "can't decode base64")
		panic("can't decode base64")
	}

	guid := xid.New()
	imageFilePath := filepath.Join("/user-face-photos", guid.String()+".jpg")
	if false == awsutils.UploadS3(data, imageFilePath) {
		httputils.RespondError(w, http.StatusBadRequest, "can't upload the photo")
		panic("can't upload the photo")
	}
	urlStr := filepath.Join("https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/", imageFilePath)

	var user dbutils.User
	var user_face_photo dbutils.UserFacePhoto

	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	} else {
		user_face_photo.Userid = user.Userid
		user_face_photo.Url = urlStr
		user_face_photo.XID = guid.String()
		if err = db.Create(&user_face_photo).Error; err != nil {
			httputils.RespondError(w, http.StatusInternalServerError, "Can't make record")
			panic("Can't make record")
		}
	}

	httputils.RespondJson(w, http.StatusOK, user_face_photo)
}
