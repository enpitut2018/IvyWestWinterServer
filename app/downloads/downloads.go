package downloads

import (
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/jinzhu/gorm"
	l "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func CreateDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}

func GetDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user models.User
	token := r.Header.Get("Authorization")
	if err := user.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
		l.Errorf("Not valid token.")
		return
	}

	var downloads models.Downloads
	if err := downloads.GetDownloadsByUserID(db, user.UserID); err != nil {
		httputils.RespondError(w, http.StatusInternalServerError, err.Error())
		l.Errorf(err.Error())
		return
	}

	httputils.RespondJson(w, http.StatusOK, downloads.Downloads)
	l.Infof("Success.")
}

type ResUser struct {
	UserID    string `json:"id"`
	AvatarURL string `json:"avatar_url"`
}

type PhotoInfo struct {
	ID          uint      `json:"id"`
	URL         string    `json:"url"`
	Uploader    ResUser   `json:"uploader"`
	Downloaders []ResUser `json:"downloaders"`
}

func GetDownloadPhotoInfo(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	// 自分が写っている写真の情報を返す。
	token := r.Header.Get("Authorization")
	var user models.User // tokenを持っている本人のuser info
	if err := user.GetUserFromToken(db, token); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "Not valid token.")
		l.Errorf("Not valid token.")
		return
	}
	queryString := r.URL.Query().Get("userids")
	l.Infof("Query String: %+v", queryString)

	var user_photos models.Downloads
	user_photos.GetDownloadsByUserID(db, user.UserID)

	var photoInfos []PhotoInfo
	for _, photo := range user_photos.Downloads {

		var uploader models.User
		uploader.SelectByUserID(db, photo.UserID)

		// 写真に写っているuserを全て返すクエリ
		var resUsers []ResUser
		if err := db.Raw(`SELECT users.user_id, users.avatar_url FROM users 
						LEFT OUTER JOIN downloads ON (downloads.user_id = users.user_id) 
						WHERE downloads.photo_id = ?`, photo.PhotoID).Scan(&resUsers).Error; err != nil {
			httputils.RespondError(w, http.StatusBadRequest, "Failed Get Users in the Photo")
			l.Errorf("Failed Get Users in the Photo, %+v", photo)
			return
		}
		if queryString == "" {
			var photoInfo PhotoInfo
			photoInfo.ID = photo.PhotoID
			photoInfo.URL = photo.URL
			photoInfo.Uploader.UserID = uploader.UserID
			photoInfo.Uploader.AvatarURL = uploader.AvatarURL
			photoInfo.Downloaders = resUsers
			photoInfos = append(photoInfos, photoInfo)
		} else {
			for _, queryUserID := range strings.Split(queryString, ",") {
				for _, resUser := range resUsers {
					if resUser.UserID == queryUserID {
						var photoInfo PhotoInfo
						photoInfo.ID = photo.PhotoID
						photoInfo.URL = photo.URL
						photoInfo.Uploader.UserID = uploader.UserID
						photoInfo.Uploader.AvatarURL = uploader.AvatarURL
						photoInfo.Downloaders = resUsers
						photoInfos = append(photoInfos, photoInfo)
					}
				}
			}
		}
	}
	httputils.RespondJson(w, http.StatusOK, photoInfos)
	return
}
