package db

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-redis/redis/v7"
	"github.com/souvikhaldar/Authentication-Wizard/pkg/signup"
)

type DB struct {
	Redis *redis.Client
}

// AddUser first checks if the user is already registered. If not
// it stores the user data in the db and stores the token for
// further verification by the user by emailing a unique link
func (db *DB) AddUser(u *signup.UserDetails, email string) error {
	// check if user is already registered
	exists, err := db.Redis.Exists(email).Result()
	if err != nil {
		return err
	}
	if exists == 1 {
		return fmt.Errorf("User exists")
	}

	// new user, so store the data and put verified=false ATM
	bodyJSON, err := json.Marshal(u)
	if err != nil {
		return err
	}
	db.Redis.Set(email, bodyJSON, 0)

	// store token:email for verification purpose
	db.Redis.Set(u.SignUpToken, email, 0)

	return nil
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
	var ud signup.UserDetails
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
	var ud signup.UserDetails

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
	var ud signup.UserDetails

	if err := json.Unmarshal([]byte(val), &ud); err != nil {
		return "", false, err
	}
	return ud.Password, ud.Verified, nil
}

func (db *DB) DeleteUserDetails(email string) error {
	val, err := db.Redis.Get(email).Result()
	if err != nil {
		return err
	}

	if val == "" {
		err := fmt.Errorf("Error in fetching the signup token")
		return err
	}
	var ud signup.UserDetails

	if err := json.Unmarshal([]byte(val), &ud); err != nil {
		return err
	}
	if err := db.Redis.Del(email, ud.SignUpToken).Err(); err != nil {
		return err
	}
	return nil
}
