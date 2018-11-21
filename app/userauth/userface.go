package userauth

import (
	"encoding/base64"
	"encoding/json"
	"github.com/enpitut2018/IvyWestWinterServer/app/awsutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/rs/xid"
	"net/http"
	"path/filepath"
	"github.com/jinzhu/gorm"
)

type SourceRequest struct {
	Source string `json:"source"`
}

func UploadUserFace(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	token := r.Header.Get("Authorization")
	
	var source SourceRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&source); err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
	}

	data, err := base64.StdEncoding.DecodeString(source.Source)
	if err != nil {
		httputils.RespondError(w, http.StatusBadRequest, "can't decode base64")
	}

	guid := xid.New()
	imageFilePath := filepath.Join("/user-face-photos", guid.String()+".jpg")
	if false == awsutils.UploadS3(data, imageFilePath) {
		httputils.RespondError(w, http.StatusBadRequest, "can't upload the photo")
	}
	urlStr := filepath.Join("https://s3-ap-northeast-1.amazonaws.com/ivy-west-winter/", imageFilePath)

	var user models.User
	if ok := user.GetUserFromToken(db, w, token); ok {
		if user.AbaterUrl != "" {
			httputils.RespondError(w, http.StatusUnauthorized, "Abater photo is already uploaded.")
		} else{
			if ok = user.UpdateAbaterUrl(db, w, urlStr); ok {
				httputils.RespondJson(w, http.StatusOK, user)
			}
		}
	}
}
