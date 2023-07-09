package domain

import (
	"errors"
	"strings"
)

type SignUpRequestBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Password  string `json:"password"`
}

func (body SignUpRequestBody) Validate() error {
	var err error

	if strings.TrimSpace(body.FirstName) == "" {
		err = errors.Join(err, errors.New("first name is empty"))
	}

	if strings.TrimSpace(body.LastName) == "" {
		err = errors.Join(err, errors.New("last name is empty"))
	}

	if strings.TrimSpace(body.Username) == "" {
		err = errors.Join(err, errors.New("username is empty"))
	}

	if strings.TrimSpace(body.Password) == "" {
		err = errors.Join(err, errors.New("password is empty"))
	} else if len(body.Password) < 8 {
		err = errors.Join(err, errors.New("password is too short"))
	}

	return err
}
