package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
	"net/http"
	"strings"
)

func downloads(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	fmt.Println("token=%+v", queries["token"][0])

	db := connectPostgres()
	defer db.Close()

	var photos []Photo
	if err := db.Raw("SELECT * FROM Downloads WHERE User.Token = ?", queries["token"][0]).Scan(&photos).Error; err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	respondJson(w, http.StatusOK, photos)
}

func uploads(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	fmt.Println("token=%+v", queries["token"][0])

	db := connectPostgres()
	defer db.Close()

	var photos []Photo
	if err := db.Raw("SELECT * FROM photos WHERE uploader.token = ?", queries["token"][0]).Scan(&photos).Error; err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	respondJson(w, http.StatusOK, photos)
}

func photo(w http.ResponseWriter, r *http.Request) {
	queries := r.URL.Query()
	token := queries["token"][0]
	decoder := json.NewDecoder(r.Body)
	var photo Photo
	if err := decoder.Decode(&photo); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	db := connectPostgres()
	defer db.Close()

	// create record
	var user User
	db.Raw("SELECT * FROM users WHERE token = ?", token).Scan(&user)
	photo.Uploader = user
	fmt.Printf("%+v\n", user)
	if err := db.Create(&photo).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		panic(err.Error())
	}

	respondJson(w, http.StatusOK, photo)
}

func signup(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newuser User
	if err := decoder.Decode(&newuser); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	db := connectPostgres()
	defer db.Close()

	// check already user
	var olduser User
	db.Raw("SELECT * FROM USERS WHERE userid = ?", newuser.Userid).Scan(&olduser)
	if olduser.Userid == newuser.Userid {
		respondError(w, http.StatusInternalServerError, "userid is already used!")
	} else {
		// create new user
		newuser.Token = getToken(newuser.Userid)
		if err := db.Create(&newuser).Error; err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			panic(err)
		}
		respondJson(w, http.StatusOK, map[string]string{"message": "user created!"})
	}
}

func signin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newuser User
	if err := decoder.Decode(&newuser); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		panic(err.Error())
	}

	db := connectPostgres()
	defer db.Close()
	var olduser User
	db.Raw("SELECT * FROM USERS WHERE userid = ?", newuser.Userid).Scan(&olduser)
	if olduser.Userid != newuser.Userid {
		respondError(w, http.StatusBadRequest, "user is not registered!")
	} else {
		if olduser.Password != newuser.Password {
			respondError(w, http.StatusBadRequest, "password is different")
		} else {
			newuser.Token = olduser.Token
			respondJson(w, http.StatusOK, newuser)
		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/downloads", downloads).Methods("GET")
	myRouter.HandleFunc("/uploads", uploads).Methods("GET")
	myRouter.HandleFunc("/signup", signup).Methods("POST")
	myRouter.HandleFunc("/signin", signin).Methods("POST")
	myRouter.HandleFunc("/photo", photo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func getToken(userid string) string {
	h := md5.New()
	h.Write([]byte(strings.ToLower(userid)))
	return hex.EncodeToString(h.Sum(nil))
}

func main() {
	fmt.Println("init migration")
	initialMigration()
	fmt.Println("------ connect start localhost:8080/ -------")
	handleRequests()
}
