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

type Downloader struct {
	ID        string `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

type PhotoInfo struct {
	ID       uint   `json:"id"`
	URL      string `json:"url"`
	Uploader struct {
		ID        string `json:"id"`
		AvatarURL string `json:"avatar_url"`
	} `json:"uploader"`
	Downloaders []Downloader `json:"downloaders"`
}

func GetPhotoInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// 自分が送信した写真の情報を返す。
	token := r.Header.Get("Authorization")
	var uploader models.User
	var photoInfos []PhotoInfo
	if err := uploader.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
		l.Errorf("Not valid token.")
		return
	}

	var uploader_photos models.Uploads
	uploader_photos.GetPhotosByUserID(db, uploader.UserID)

	for _, photo := range uploader_photos.Uploads {
		var photoInfo PhotoInfo
		photoInfo.ID = photo.ID
		photoInfo.URL = photo.URL
		photoInfo.Uploader.ID = uploader.UserID
		photoInfo.Uploader.AvatarURL = uploader.AvatarURL

		var downloads models.Downloads
		if err := downloads.GetDownloadsByPhotoID(db, photo.ID); err != nil {
			httputils.RespondError(w, http.StatusBadRequest, "Failed Get PhotoInfo.")
			l.Errorf("Failed Get PhotoInfo. id=%+v", photo.ID)
			return
		}

		for _, photo := range downloads.Downloads {
			var downloader models.User
			var d Downloader
			if err := downloader.SelectByUserID(db, photo.UserID); err != nil {
				httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
				l.Errorf("Can't Get Downloader Info. id=%+v", photo.UserID)
				return
			}
			d.ID = downloader.UserID
			d.AvatarURL = downloader.AvatarURL
			photoInfo.Downloaders = append(photoInfo.Downloaders, d)
		}
		photoInfos = append(photoInfos, photoInfo)
	}
	httputils.RespondJson(w, http.StatusOK, photoInfos)
}

func GetUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var uploads models.Uploads
	var user models.User
	token := r.Header.Get("Authorization")
	if err := user.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
		l.Errorf("Not valid token.")
		return
	}
	if err := uploads.GetPhotosByUserID(db, user.UserID); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Can't get phots by users.")
		l.Errorf("Can't get phots by users.")
		return
	}
	httputils.RespondJson(w, http.StatusOK, uploads.Uploads)
	l.Infof("Success")
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
		return
	}
	url := awsutils.UploadPhoto(w, source.Source, s3FolderPath)

	var user models.User
	if err := user.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
		l.Errorf("Not valid token.")
		return
	}

	upload := models.Upload{UserID: user.UserID, URL: url}
	if err := upload.CreateRecord(db); err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, err.Error())
		l.Errorf(err.Error())
		return
	}
	// 顔認識を使用してDownloadテーブルにレコードを追加する。
	downloadUserIDs, err := faceidentification.FaceIdentification(db, w, upload.URL)
	if err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, err.Error())
		l.Errorf(err.Error())
		return
	}
	res := UploadResponse{UserID: upload.UserID, URL: upload.URL, DownloadUserIDs: downloadUserIDs}
	httputils.RespondJson(w, http.StatusOK, res)
	l.Infof("Success")
	return
}

func DeleteUploads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
	l.Infof("Success")
}
