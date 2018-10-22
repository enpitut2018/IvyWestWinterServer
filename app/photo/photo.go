package photo

import (
	"encoding/base64"
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/awsutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/dbutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
	"path/filepath"
	"strconv"
)

func Downloads(w http.ResponseWriter, r *http.Request) {
	db := dbutils.ConnectPostgres()
	defer db.Close()

	var photos []dbutils.Photo
	token := r.Header.Get("Authorization")
	var user dbutils.User
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	}
	if err := db.Raw("SELECT * FROM Downloads WHERE userid = ?", user.Userid).Scan(&photos).Error; err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	httputils.RespondJson(w, http.StatusOK, photos)
}

func Uploads(w http.ResponseWriter, r *http.Request) {
	db := dbutils.ConnectPostgres()
	defer db.Close()

	var photos []dbutils.Photo
	token := r.Header.Get("Authorization")
	var user dbutils.User
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	}
	if err := db.Raw("SELECT * FROM photos WHERE userid = ?", user.Userid).Scan(&photos).Error; err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	httputils.RespondJson(w, http.StatusOK, photos)
}

func Photo(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	decoder := json.NewDecoder(r.Body)
	var photo dbutils.Photo
	if err := decoder.Decode(&photo); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	db := dbutils.ConnectPostgres()
	defer db.Close()

	// upload to S3
	data, err := base64.StdEncoding.DecodeString(photo.Source)
	if err != nil {
		panic("can't decode base64")
	}
	imageFilePath := filepath.Join("/upload-photos", strconv.FormatUint(uint64(photo.Model.ID), 10)+".jpg")
	err := awsutils.UploadS3(data, imageFilePath)
	if err == nil {
		panic("can't upload the photo")
	}
	urlStr := filepath.Join("https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/", imageFilePath)

	// create record
	var user dbutils.User
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	} else {
		photo.Userid = user.Userid
		photo.Url = urlStr
	}

	if err := db.Create(&photo).Error; err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, err.Error())
		panic(err.Error())
	}

	httputils.RespondJson(w, http.StatusOK, map[string]string{"URL": urlStr})
}
