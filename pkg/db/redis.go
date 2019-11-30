package db

import (
	"encoding/json"
	"fmt"
	"log"

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
	db.Redis.Set(email, bodyJSON, 0)
	return signupToken, nil
}

func (db *DB) FetchToken(email string) (string, error) {
	val, err := db.Redis.Get(email).Result()
	if err != nil {
		return "", err
	}
	if val == "" {
		err := fmt.Errorf("Error in fetching the signup token")
		return "", err
	}
	var ud UserDetails
	if err := json.Unmarshal([]byte(val), &ud); err != nil {
		return "", err
	}
	log.Println(ud)
	return ud.SignUpToken, nil
}

func (db *DB) UpdateValidity(email string) error {
	val, err := db.Redis.Get(email).Result()
	if err != nil {
		return err
	}

	if val == "" {
		err := fmt.Errorf("Error in fetching the signup token")
		return err
	}
	var ud UserDetails

	if err := json.Unmarshal([]byte(val), &ud); err != nil {
		return err
	}
	ud.Verified = true
	bodyJSON, err := json.Marshal(&ud)
	if err != nil {
		return err
	}
	db.Redis.Set(email, bodyJSON, 0)
	log.Println(ud)
	return nil
}

func (db *DB) FetchPasswordAndStatus(email string) (string, bool, error) {
	val, err := db.Redis.Get(email).Result()
	if err != nil {
		return "", false, err
	}

	if val == "" {
		err := fmt.Errorf("Error in fetching the signup token")
		return "", false, err
	}
	var ud UserDetails

	if err := json.Unmarshal([]byte(val), &ud); err != nil {
		return "", false, err
	}
	return ud.Password, ud.Verified, nil
}
