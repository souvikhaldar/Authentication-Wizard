package login

import (
	"github.com/souvikhaldar/Authentication-Wizard/pkg/signup"
	"log"
)

type Repository interface {
	FetchPasswordAndStatus(email string) (string, bool, error)
}

// IsRegistered checks if the user has already signup and verified
// his email or not
func IsRegistered(repo Repository, email, password string) bool {
	p, s, err := repo.FetchPasswordAndStatus(email)
	if err != nil {
		log.Println(err)
		return false
	}
	return signup.Hash(password) == p && s
}
