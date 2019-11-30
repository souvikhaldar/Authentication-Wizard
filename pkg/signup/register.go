package signup

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	user := new(SignupBody)
	if err := json.NewDecoder(r.Body).Decode(user); err != nil {
		e := fmt.Errorf("Error in parsing payload: %s", err)
		log.Println(e)
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

}

