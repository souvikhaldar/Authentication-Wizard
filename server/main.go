package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis/v7"
	"github.com/gorilla/mux"
	"github.com/souvikhaldar/Authentication-Wizard/pkg/db"
	"github.com/souvikhaldar/Authentication-Wizard/pkg/login"
	"github.com/souvikhaldar/Authentication-Wizard/pkg/signup"
	"github.com/souvikhaldar/gomail"
)

type httpServer struct {
	DB     *db.DB
	Router *mux.Router
}
type SignupBody struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
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

	return func(w http.ResponseWriter, r *http.Request) {

		user := new(SignupBody)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			e := fmt.Errorf("Error in parsing payload: %s", err)
			log.Println(e)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		if !signup.IsValidEmail(user.EmailID) {
			err := fmt.Errorf("Invalid input email")
			log.Println(err)
			http.Error(w, err.Error(), 404)
		}

		token, err := s.DB.AddUser(user.EmailID, user.Password)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		msg := fmt.Sprintln("Signup token is: ", token)
		log.Println(msg)
		log.Println(os.Getenv("AW_EMAIL"), os.Getenv("AW_PASSWORD"))
		e, config := gomail.New(os.Getenv("AW_EMAIL"), os.Getenv("AW_PASSWORD"))
		if e != nil {
			fmt.Print(fmt.Errorf("Error in creating config %v", e))
		}
		if e := config.SendMail([]string{user.EmailID}, "Verification", fmt.Sprintf("Click on localhost:8192/verify?e=%s&t=%s", user.EmailID, token)); e != nil {
			fmt.Print(fmt.Errorf("Error in sending mail %v", e))
		}
		w.Write([]byte("Please verify your email address" + msg))
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
			err := fmt.Errorf("Email ID %s does not exist", email)
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		log.Println(t, t)
		if token == t {
			log.Println("Verified")
			if err := s.DB.UpdateValidity(email); err != nil {
				er := fmt.Errorf("Failed to update the validity %s", err.Error())
				log.Println(er)
			}
			w.Write([]byte("Successfully Verified."))
			return
		}
		w.Write([]byte("Verification Unsuccessfull"))
		return
	}
}

func (s *httpServer) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := new(SignupBody)
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			e := fmt.Errorf("Error in parsing payload: %s", err)
			log.Println(e)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		if !login.IsRegistered(s.DB, user.EmailID, user.Password) {
			log.Println("User is not registered")
			http.Error(w, "User is not registerd", 403)
			return
		}
		log.Println("Logged in")
		w.Write([]byte("Successfully logged in"))
	}
}

func main() {
	s := NewServer()
	s.Router.HandleFunc("/signup", s.RegisterNewUser()).Methods("POST")
	s.Router.HandleFunc("/verify", s.Verify()).Methods("GET")
	s.Router.HandleFunc("/login", s.Login()).Methods("POST")
	log.Fatal(http.ListenAndServe(":8192", s))

}
