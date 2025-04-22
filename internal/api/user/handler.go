package user

import (
	"net/http"
	"socialAPI/internal/lib"
	"strconv"

	"github.com/go-chi/render"
)

func (c UserController) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметр excludeID из строки запроса (например, /users?excludeID=3)
		excludeIDParam := r.URL.Query().Get("excludeID")

		var excludeID *uint
		if excludeIDParam != "" {
			// Преобразуем excludeIDParam в uint
			id, err := strconv.Atoi(excludeIDParam)
			if err != nil {
				lib.SendMessage(w, r, http.StatusBadRequest, "Invalid excludeID parameter")
				return
			}
			// Преобразуем в указатель
			idUint := uint(id)
			excludeID = &idUint
		}

		// Теперь передаем excludeID в сервис
		users, err := c.userService.GetAllUsers(excludeID)
		if err != nil {
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		// Отправляем успешный ответ
		render.Status(r, http.StatusOK)
		render.JSON(w, r, users)
	}
}
