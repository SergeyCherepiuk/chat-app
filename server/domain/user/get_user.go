package userdomain

type GetUserResponseBody struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Username    string `json:"username"`
	Description string `json:"description"`
}
