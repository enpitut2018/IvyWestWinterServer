package main

import (
	"./dbutils"
	"./photo"
	"./userauth"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/downloads", photo.Downloads).Methods("GET")
	myRouter.HandleFunc("/uploads", photo.Uploads).Methods("GET")
	myRouter.HandleFunc("/signup", userauth.Signup).Methods("POST")
	myRouter.HandleFunc("/signin", userauth.Signin).Methods("POST")
	myRouter.HandleFunc("/photo", photo.Photo).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	dbutils.InitialMigration()
	// initialSQL()
	fmt.Println("------ connect start localhost:8080/ -------")
	handleRequests()
}
