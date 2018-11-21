package upload

import (
	"encoding/base64"
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/awsutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/enpitut2018/IvyWestWinterServer/app/faceidentification"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/rs/xid"
	"net/http"
	"path/filepath"
	"github.com/jinzhu/gorm"
)

func GetUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {

	var photos []models.Photo
	token := r.Header.Get("Authorization")
	var user models.User
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

type Source struct {
	Source string
}

func CreateUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := r.Header.Get("Authorization")
	var source Source
	decoder := json.NewDecoder(r.Body)
	var photo models.Photo
	if err := decoder.Decode(&source); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	// upload to S3
	data, err := base64.StdEncoding.DecodeString(source.Source)
	if err != nil {
		panic("can't decode base64")
	}
	guid := xid.New() // xidというユニークなID
	imageFilePath := filepath.Join("/upload-photos", guid.String()+".jpg")
	if false == awsutils.UploadS3(data, imageFilePath) {
		panic("can't upload the photo")
	}
	urlStr := filepath.Join("https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/", imageFilePath)

	// create record
	var user models.User
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	} else {
		photo.Userid = user.Userid
		photo.Url = urlStr
		photo.XID = guid.String()
		if err = db.Create(&photo).Error; err != nil {
			httputils.RespondError(w, http.StatusInternalServerError, "Can't make record")
			panic("Can't make record")
		}
	}

	// face identification
	// 顔認識技術を使用してDownloadテーブルにレコードを追加する。
	faceidentification.FaceIdentification(db, photo.Url)

	httputils.RespondJson(w, http.StatusOK, photo)
}

func DeleteUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}
