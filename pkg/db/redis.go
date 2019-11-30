package db

import (
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/souvikhaldar/Authentication-Wizard/pkg/signup"
	"github.com/souvikhaldar/gorand"
)

type DB struct {
	Redis *redis.Client
}

func (db *DB) AddUser(email string, password string) (string, error) {
	signupToken := gorand.RandStr(5)
	UserDetails := &UserDetails{
		signup.Hash(password),
		signupToken,
		false,
	}
	db.Redis.Set(signupToken, email, 0)
	bodyJSON, err := json.Marshal(UserDetails)
	if err != nil {
		return "", err
	}
	fmt.Println("Body JSON:", bodyJSON)
	db.Redis.Set(email, bodyJSON, 0)
	return signupToken, nil
}

func (db *DB) FetchToken(email string) (string, error) {
	val := db.Redis.Get(email).String()

	if val == "" {
		err := fmt.Errorf("Error in fetching the signup token")
		return "", err
	}
	var ud UserDetails

	if err := json.Unmarshal([]byte(val), &ud); err != nil {
		return "", err
	}
	return ud.SignUpToken, nil
}

func (db *DB) UpdateValidity(email string) error {
	val := db.Redis.Get(email).String()

	if val == "" {
		err := fmt.Errorf("Error in fetching the signup token")
		return err
	}
	var ud UserDetails

	if err := json.Unmarshal([]byte(val), &ud); err != nil {
		return err
	}
	ud.Verified = true
	db.Redis.Set(email, ud, 0)
	return nil
}
