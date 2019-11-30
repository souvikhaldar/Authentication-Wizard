package signup

import (
	"crypto/sha1"
	"fmt"
	"regexp"
)

func ValidateEmail(email string) bool {
	var rxEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return rxEmail.MatchString(email)
}

func Hash(password string) string {
	h := sha1.New()
	h.Write([]byte(password))
	return fmt.Sprintf("%x", bs)
}
