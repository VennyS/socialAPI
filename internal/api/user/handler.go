package user

import (
	"net/http"
	"socialAPI/internal/api"
	"socialAPI/internal/lib"

	"github.com/go-chi/render"
)

func (c UserController) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем параметр excludeMe из строки запроса (например, /users?excludeMe=true)
		excludeMe := r.URL.Query().Get("excludeMe")

		// Флаг для исключения пользователя из выборки
		var excludeID *uint

		// Если excludeMe равно "true", исключаем текущего пользователя
		if excludeMe == "true" {
			// Извлекаем userID из контекста (считаем, что userID установлен в контексте с помощью авторизации)
			userID, ok := r.Context().Value(api.UserIDKey).(uint)
			if !ok {
				lib.SendMessage(w, r, http.StatusUnauthorized, "User ID not found in context")
				return
			}

			// Преобразуем ID в указатель
			excludeID = &userID
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
