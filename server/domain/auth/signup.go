package authdomain

import (
	"errors"
	"regexp"
)

type SignUpRequestBody struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	Username       string `json:"username"`
	Password       string `json:"password"`
	Description    string `json:"description"`
	ProfilePicture []byte `json:"profile_picture"`
}

var firstNameRegexp = regexp.MustCompile(`^[A-Z][a-zA-Z]+(?:-[A-Z][A-Za-z]+)?$`)
var lastNameRegexp = regexp.MustCompile(`^[A-Z][a-zA-Z]+(?: [A-Z][A-Za-z]+)?$`)
var usernameRegexp = regexp.MustCompile(`^[a-zA-Z0-9_]{5,20}$`)
var lowercase = regexp.MustCompile(`[a-z]`)
var uppercase = regexp.MustCompile(`[A-Z]`)
var digit = regexp.MustCompile(`[0-9]`)
var special = regexp.MustCompile(`[!@#$%&*]`)
var eightOrMore = regexp.MustCompile(`.{8,}`)

func (body SignUpRequestBody) Validate() error {
	var err error

	if !firstNameRegexp.Match([]byte(body.FirstName)) {
		err = errors.Join(err, errors.New("invalid first name"))
	}
	if !lastNameRegexp.Match([]byte(body.LastName)) {
		err = errors.Join(err, errors.New("invalid last name"))
	}
	if !usernameRegexp.Match([]byte(body.Username)) {
		err = errors.Join(err, errors.New("invalid username"))
	}

	password := []byte(body.Password)
	if eightOrMore.Find(password) == nil {
		err = errors.Join(err, errors.New("password must be at least 8 characters long"))
	} else if lowercase.Find(password) == nil {
		err = errors.Join(err, errors.New("password must contain at least one lowercase character"))
	} else if uppercase.Find(password) == nil {
		err = errors.Join(err, errors.New("password must contain at least one uppercase character"))
	} else if digit.Find(password) == nil {
		err = errors.Join(err, errors.New("password must contain at least one digit"))
	} else if special.Find(password) == nil {
		err = errors.Join(err, errors.New("password must contain at least one special character"))
	}

	return err
}
