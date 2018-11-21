package download

import (
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
	"github.com/jinzhu/gorm"
)

func CreateDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}

func GetDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {

	var user models.User
	token := r.Header.Get("Authorization")
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	}
	var photos []models.Photo
	if err := db.Raw("SELECT * FROM photos WHERE userid = ?", user.Userid).Scan(&photos).Error; err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	httputils.RespondJson(w, http.StatusOK, photos)
}

func DeleteDownloads(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}
