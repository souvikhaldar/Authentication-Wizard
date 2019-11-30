package login

import (
	"github.com/souvikhaldar/Authentication-Wizard/pkg/signup"
	"log"
)

type Repository interface {
	FetchPasswordAndStatus(email string) (string, bool, error)
}

func IsRegistered(repo Repository, email, password string) bool {
	p, s, err := repo.FetchPasswordAndStatus(email)
	if err != nil {
		log.Println(err)
		return false
	}
	return signup.Hash(password) == p && s
}
