package signup

type SignupBody struct {
	EmailID  string `json:"email_id"`
	Password string `json:"password"`
}
