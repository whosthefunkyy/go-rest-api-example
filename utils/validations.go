
package utils

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/whosthefunkyy/go-rest-api-example/models"
)

var (
	ErrInvalidBody = errors.New("invalid request body")
	ErrEmptyName   = errors.New("name cannot be empty")
)

func ParseAndValidateUser(r *http.Request) (*models.User, error) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		return nil, ErrInvalidBody
	}
	if user.Name == "" {
		return nil, ErrEmptyName
	}
	return &user, nil
}
