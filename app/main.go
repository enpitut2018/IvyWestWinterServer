package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./downloads"
	"./uploads"
	"./userauth"
	"./models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type App struct {
	Router *mux.Router
	DB *gorm.DB
}

func (app *App) Initialize() {
	var err error
	app.DB, err = gorm.Open("postgres", os.Getenv("DATABASE_URL")) // osパッケージが必要
	if err != nil {
		panic("Failed to connect to database")
	}
	// defer app.DB.Close()
	app.DB.DB().SetMaxIdleConns(0)
	app.DB.AutoMigrate(&models.User{}, &models.Upload{}, &models.Download{})
}

func (app *App) Run() {
	defer app.DB.Close()
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/downloads", handlerWithDB(downloads.CreateDownloads, app.DB)).Methods("POST")
	myRouter.HandleFunc("/downloads", handlerWithDB(downloads.GetDownloads, app.DB)).Methods("GET")
	myRouter.HandleFunc("/downloads", handlerWithDB(downloads.DeleteDownloads, app.DB)).Methods("DELETE")
	myRouter.HandleFunc("/uploads", handlerWithDB(uploads.CreateUploads, app.DB)).Methods("POST")
	myRouter.HandleFunc("/uploads", handlerWithDB(uploads.GetUploads, app.DB)).Methods("GET")
	myRouter.HandleFunc("/uploads", handlerWithDB(uploads.DeleteUploads, app.DB)).Methods("DELETE")
	myRouter.HandleFunc("/uploadUserFace", handlerWithDB(userauth.UploadUserFace, app.DB)).Methods("POST")
	myRouter.HandleFunc("/signup", handlerWithDB(userauth.Signup, app.DB)).Methods("POST")
	myRouter.HandleFunc("/signin", handlerWithDB(userauth.Signin, app.DB)).Methods("POST")
	myRouter.HandleFunc("/user", handlerWithDB(userauth.GetUserInfo, app.DB)).Methods("GET")
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func handlerWithDB(fn func(w http.ResponseWriter, r *http.Request, DB *gorm.DB), DB *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, DB)
	}
}

func main() {
	app := App{}
	app.Initialize()
	fmt.Println("\n------ connect start localhost:8080/ -------\n")
	app.Run()
}
