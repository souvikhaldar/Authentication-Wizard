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

func (s *httpServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.Router.ServeHTTP(w, r)
}

// newServer sets up a new server and returns pointer to
// the new server instance
func newServer() *httpServer {
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

// RegisterNewUser registers/signs up a new user after verifying
// the AUTHENTICATION using email
func (s *httpServer) RegisterNewUser() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		user := new(SignupBody)
		// parse the JSON payload
		if err := json.NewDecoder(r.Body).Decode(user); err != nil {
			e := fmt.Errorf("Error in parsing payload: %s", err)
			log.Println(e)
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}

		// Check if the provided email is in correct format
		if !signup.IsValidEmail(user.EmailID) {
			err := fmt.Errorf("Invalid input email")
			log.Println(err)
			http.Error(w, err.Error(), 404)
			return
		}
		// store the store data until further verification
		token, err := signup.SignupUser(s.DB, user.EmailID, user.Password)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		log.Println(fmt.Sprintln("Sign-up token is: ", token))
		if os.Getenv("AW_EMAIL") == "" || os.Getenv("AW_PASSWORD") == "" {
			err := fmt.Errorf(`Sender email and password is not exported. Please export them. Eg. export AW_PASSWORD=<password> export AW_EMAIL=<email>`)
			log.Println(err)
			http.Error(w, err.Error(), 422)
			return
		}

		// setup the email client
		e, config := gomail.New(os.Getenv("AW_EMAIL"), os.Getenv("AW_PASSWORD"))
		if e != nil {
			err = fmt.Errorf("Error in creating email config %v", e)
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}

		emailBody := fmt.Sprintf("Copy paste this URL on browser to verify: localhost:8192/verify?e=%s&t=%s", user.EmailID, token)

		if e := config.SendMail([]string{user.EmailID}, "Verification", emailBody); e != nil {
			err = fmt.Errorf("Error in sending mail %v", e)
			log.Println(err)
			// since there was error in sending mail
			// hence the user can't sign up
			// so need to delete the user from db
			// so that he/she can try again
			// delete the user details from the db
			if err := s.DB.DeleteUserDetails(user.EmailID); err != nil {
				log.Println("Error in deleting user details upon failed verification: ", err)
			}
			http.Error(w, err.Error(), 500)
			return
		}
		w.Write([]byte("Mail sent. Please verify your email address: " + emailBody))
	}
}

// Verify verifies that the user is authentic
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

		t, err := s.DB.FetchToken(email)
		if err != nil {
			err := fmt.Errorf("Email ID %s does not exist", email)
			log.Println(err)
			http.Error(w, err.Error(), 500)
			return
		}
		if token == t {
			log.Println("Verified")
			if err := s.DB.UpdateValidity(email); err != nil {
				er := fmt.Errorf("Failed to update the validity %s", err.Error())
				log.Println(er)
			}
			w.Write([]byte("Successfully Verified."))
			return
		}
		// delete the user details from the db
		if err := s.DB.DeleteUserDetails(email); err != nil {
			log.Println("Error in deleting user details upon failed verification: ", err)
		}
		w.Write([]byte("Verification Unsuccessfull"))
		return
	}
}

// Login logs in the user if the credentials are correct
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
	s := newServer()
	// check if email client is setup
	if os.Getenv("AW_EMAIL") == "" || os.Getenv("AW_PASSWORD") == "" {
		err := fmt.Errorf(`Sender email and password is not exported. Please export them. Eg. export AW_PASSWORD=<password> export AW_EMAIL=<email>`)
		log.Fatal(err)
	}
	s.Router.HandleFunc("/signup", s.RegisterNewUser()).Methods("POST")
	s.Router.HandleFunc("/verify", s.Verify()).Methods("GET")
	s.Router.HandleFunc("/login", s.Login()).Methods("POST")
	log.Println("Server running on: 8192")
	log.Fatal(http.ListenAndServe(":8192", s))

}
