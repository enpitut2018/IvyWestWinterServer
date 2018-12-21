package main

import (
	"log"
	"net/http"
	"os"

	l "github.com/sirupsen/logrus"

	"github.com/enpitut2018/IvyWestWinterServer/app/downloads"
	"github.com/enpitut2018/IvyWestWinterServer/app/models"
	"github.com/enpitut2018/IvyWestWinterServer/app/uploads"
	"github.com/enpitut2018/IvyWestWinterServer/app/userauth"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/labstack/echo"
)

type App struct {
	Router *mux.Router
	//e      *echo.Echo
	DB *gorm.DB
}

func (app *App) Initialize() {
	var err error
	app.DB, err = gorm.Open("postgres", os.Getenv("DATABASE_URL")) // osパッケージが必要
	if err != nil {
		panic("Failed to connect to database")
	}
	app.DB.DB().SetMaxIdleConns(0)
	app.DB.AutoMigrate(&models.User{}, &models.Upload{}, &models.Download{})

	// app.e := echo.New()
	// app.e.Use(middleware.Logger())
	// app.e.Use(middleware.Recover())
}

func (app *App) Run() {
	defer app.DB.Close()
	// app.e.GET("/downloads", downloads.CreateDownloads(app.DB))
	// app.e.POST("/uploads", uploads.CreateUploads(app.DB))
	// app.e.POST("/signup", userauth.Signup(app.DB))
	// app.e.POST("/signin", userauth.Signin(app.DB))
	// app.e.GET("/user", userauth.GetUserInfo(app.DB))
	// l.Fatal(app.e.Start(":" + os.Getenv("PORT")))

	app.Router = mux.NewRouter().StrictSlash(true)
	app.Router.HandleFunc("/downloads", handlerWithDB(downloads.CreateDownloads, app.DB)).Methods("POST")
	app.Router.HandleFunc("/downloads", handlerWithDB(downloads.GetDownloads, app.DB)).Methods("GET")
	app.Router.HandleFunc("/uploads", handlerWithDB(uploads.CreateUploads, app.DB)).Methods("POST")
	app.Router.HandleFunc("/uploads", handlerWithDB(uploads.GetUploads, app.DB)).Methods("GET")
	app.Router.HandleFunc("/uploadPhotoInfos", handlerWithDB(uploads.GetUploadPhotoInfo, app.DB)).Methods("GET")
	app.Router.HandleFunc("/downloadPhotoInfos", handlerWithDB(downloads.GetDownloadPhotoInfo, app.DB)).Methods("GET")
	app.Router.HandleFunc("/uploadUserFace", handlerWithDB(userauth.UploadUserFace, app.DB)).Methods("POST")
	app.Router.HandleFunc("/signup", handlerWithDB(userauth.Signup, app.DB)).Methods("POST")
	app.Router.HandleFunc("/signin", handlerWithDB(userauth.Signin, app.DB)).Methods("POST")
	app.Router.HandleFunc("/user", handlerWithDB(userauth.GetUserInfo, app.DB)).Methods("GET")
	app.Router.HandleFunc("/users", handlerWithDB(userauth.GetUsersInfo, app.DB)).Methods("GET")
	app.Router.Use(loggingMiddleware)
	l.Info(http.ListenAndServe(":"+os.Getenv("PORT"), app.Router))
}

func handlerWithDB(fn func(w http.ResponseWriter, r *http.Request, DB *gorm.DB), DB *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fn(w, r, DB)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%v %v", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func main() {
	app := App{}
	app.Initialize()
	l.SetReportCaller(true)
	l.Infof("connect localhost:8080/")
	app.Run()
}
