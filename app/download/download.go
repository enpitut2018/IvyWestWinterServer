package download

import (
	"github.com/enpitut2018/IvyWestWinterServer/app/dbutils"
	"github.com/enpitut2018/IvyWestWinterServer/app/httputils"
	"net/http"
)

func CreateDownloads(w http.ResponseWriter, r *http.Request) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}

func GetDownloads(w http.ResponseWriter, r *http.Request) {
	db := dbutils.ConnectPostgres()
	defer db.Close()

	var user dbutils.User
	token := r.Header.Get("Authorization")
	if err := db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user).Error; err != nil {
		httputils.RespondError(w, http.StatusUnauthorized, err.Error())
		panic(err.Error())
	}
	var downloads []dbutils.Download
	if err := db.Raw("SELECT * FROM downloads WHERE userid = ?", user.Userid).Scan(&downloads).Error; err != nil {
		httputils.RespondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	httputils.RespondJson(w, http.StatusOK, downloads)
}

func DeleteDownloads(w http.ResponseWriter, r *http.Request) {
	httputils.RespondJson(w, http.StatusOK, map[string]string{"message": "Sorry. Not Implement."})
}
