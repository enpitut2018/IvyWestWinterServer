package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"./dbutils"
	"./download"
	"./upload"
	"./userauth"
	"github.com/gorilla/mux"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/downloads", download.CreateDownloads).Methods("POST")
	myRouter.HandleFunc("/downloads", download.GetDownloads).Methods("GET")
	myRouter.HandleFunc("/downloads", download.DeleteDownloads).Methods("DELETE")
	myRouter.HandleFunc("/uploads", upload.CreateUploads).Methods("POST")
	myRouter.HandleFunc("/uploads", upload.GetUploads).Methods("GET")
	myRouter.HandleFunc("/uploads", upload.DeleteUploads).Methods("DELETE")
	myRouter.HandleFunc("/signup", userauth.Signup).Methods("POST")
	myRouter.HandleFunc("/signin", userauth.Signin).Methods("POST")
	// herokuデプロイ用にコメントアウト
	// log.Print(http.ListenAndServe(":8080", myRouter))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), myRouter))
}

func main() {
	dbutils.InitialMigration()
	// initialSQL()
	fmt.Println("------ connect start localhost:8080/ -------")
	handleRequests()
}
