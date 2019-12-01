package signup

import (
	"github.com/souvikhaldar/gorand"
)

type Repository interface {
	AddUser(u *UserDetails, email string) error
}

func SignupUser(repo Repository, EmailID, Password string) (string, error) {
	signupToken := gorand.RandStr(5)
	UserDetails := &UserDetails{
		Hash(Password),
		signupToken,
		false,
	}
	if err := repo.AddUser(UserDetails, EmailID); err != nil {
		return signupToken, err
	}
	return signupToken, nil
}
