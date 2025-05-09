package lib

import (
	"fmt"
	"net/http"
	"socialAPI/internal/storage/repository"
	"strings"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePasswords(hashedPwd string, plainPwd string) error
}

type BcryptHasher struct{}

func (b *BcryptHasher) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (b *BcryptHasher) ComparePasswords(hashedPwd string, plainPwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(plainPwd))
}

func ValidateFields(fields map[string]string) error {
	var missing []string

	for field, value := range fields {
		if value == "" {
			missing = append(missing, field)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}

func SendMessage(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	render.Status(r, statusCode)
	render.JSON(w, r, map[string]string{"message": message})
}

func IsValidStatus(status repository.FriendshipStatus) bool {
	switch status {
	case repository.StatusPending, repository.StatusRejected, repository.StatusFriendship:
		return true
	default:
		return false
	}
}
