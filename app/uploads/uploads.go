package uploads

import (
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/awsutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/enpitut2018/IvyWestWinterServer/app/faceidentification"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
	"github.com/jinzhu/gorm"
)

func GetUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var uploads models.Uploads
	var user models.User
	token := r.Header.Get("Authorization")
	user.GetUserFromToken(db, w, token)
	uploads.GetPhotosByUserID(db, w, user.UserID)
	httputils.RespondJson(w, http.StatusOK, uploads.Uploads)
}

type SourceRequest struct {
	Source string
}

func CreateUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := r.Header.Get("Authorization")
	var source SourceRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&source); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	s3FolderPath := "/upload-photos"
	url := awsutils.UploadPhoto(w, source.Source, s3FolderPath)

	var user models.User
	user.GetUserFromToken(db, w, token)
	upload := models.Upload{UserID: user.UserID, URL: url}
	upload.CreateRecord(db, w)

	// face identification
	// 顔認識技術を使用してDownloadテーブルにレコードを追加する。
	faceidentification.FaceIdentification(db, w, upload.URL)

	httputils.RespondJson(w, http.StatusOK, upload)
}

func DeleteUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}
