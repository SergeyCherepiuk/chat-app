package authdomain

type LoginRequestBody struct {
	Username string `json:"username"`
	Password string `json:"password"`
}