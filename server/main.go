package main

import (
	"log"
	"net/http"
	"pkg/signup"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/signup", signup.RegisterNewUser).Methods("POST")
	router.HandleFunc("/verify", signup.Verify)
	router.HandleFunc("/login", login.Login)
	log.Fatal(http.ListenAndServe(":8192", router))

}
