package db

type UserDetails struct {
	Password    string `json:"password"`
	SignUpToken string `json:"signup_token"`
	Verified    bool   `json:"verified"`
}
