package lib

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

func ComparePasswords(hashedPwd string, plainPwd string) error {
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

func SendError(w http.ResponseWriter, r *http.Request, statusCode int, message string) {
	render.Status(r, statusCode)
	render.JSON(w, r, map[string]string{"message": message})
}
