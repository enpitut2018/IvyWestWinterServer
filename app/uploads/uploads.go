package uploads

import (
	"encoding/json"
	"net/http"

	"github.com/enpitut2018/IvyWestWinterServer/app/awsutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/faceidentification"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	l "github.com/sirupsen/logrus"
)

func GetUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var uploads models.Uploads
	var user models.User
	token := r.Header.Get("Authorization")
	if ok := user.GetUserFromToken(db, w, token); ok {
		if ok := uploads.GetPhotosByUserID(db, w, user.UserID); ok {
			httputils.RespondJson(w, http.StatusOK, uploads.Uploads)
			l.Infof("Success")
		}
	}
}

type SourceRequest struct {
	Source string
}

type UploadResponse struct {
	UserID          string   `json:"userid"`
	URL             string   `json:"url"`
	DownloadUserIDs []string `json:"downloadUserIds"`
}

const s3FolderPath = "/upload-photos"

func CreateUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := r.Header.Get("Authorization")
	var source SourceRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&source); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		l.Errorf(err.Error())
	} else {
		url := awsutils.UploadPhoto(w, source.Source, s3FolderPath)

		var user models.User
		if ok := user.GetUserFromToken(db, w, token); ok {
			upload := models.Upload{UserID: user.UserID, URL: url}
			if ok := upload.CreateRecord(db, w); ok {
				// 顔認識を使用してDownloadテーブルにレコードを追加する。
				if downloadUserIDs, ok = faceidentification.FaceIdentification(db, w, upload.URL); ok {
					res := UploadResponse{UserID: upload.UserID, URL: upload.URL, DownloadUserIDs: downloadUserIDs}
					httputils.RespondJson(w, http.StatusOK, res)
					l.Infof("Success")
				}
			}
		}
	}
}

func DeleteUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
	l.Infof("Success")
}
