package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/souvikhaldar/Authentication-Wizard/pkg/db"
)

type httpServer struct {
	DB     *db.DB
	Router *mux.Router
}

func (h *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Router.ServeHTTP(w, r)
}

func NewServer() *httpServer {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
	srv := &httpServer{}
	srv.DB = &db.DB{Redis: client}
	srv.Router = mux.NewRouter()
	return srv
}

func (s *httpServer) RegisterNewUser() http.HandlerFunc {
	type SignupBody struct {
		EmailID  string `json:"email_id"`
		Password string `json:"password"`
	}
	fmt.Println("It should run only once: ")
	return func(w http.ResponseWriter, r *http.Request) {

		user := new(SignupBody)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			e := fmt.Errorf("Error in parsing payload: %s", err)
			log.Println(e)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		token, err := s.DB.AddUser(user.EmailID, user.Password)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		fmt.Println("Signup token is: ", token)

	}
}

func (s *httpServer) Verify() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()
		if _, ok := vars["e"]; !ok {
			err := fmt.Errorf("Error in parsing query parameter")
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		if _, ok := vars["t"]; !ok {
			err := fmt.Errorf("Error in parsing query parameter")
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		email := vars["e"][0]
		token := vars["t"][0]
		log.Println(email, token)
		t, err := s.DB.FetchToken(email)
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}
		if token == t {
			log.Println("Verified")
			w.Write([]byte("Successfully Verified."))
			return
		}
		w.Write([]byte("Verification Unsuccessfull"))
		return
	}
}

func main() {
	s := NewServer()
	s.Router.HandleFunc("/signup", s.RegisterNewUser()).Methods("POST")
	s.Router.HandleFunc("/verify", s.Verify()).Methods("GET")
	//s.router.HandleFunc("/login", login.Login)
	log.Fatal(http.ListenAndServe(":8192", s))

}
