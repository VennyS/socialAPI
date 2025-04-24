package user

import (
	"net/http"
	"socialAPI/internal/api/middleware"
	"socialAPI/internal/lib"

	"github.com/go-chi/render"
)

func (c UserController) GetAllHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		excludeMe := r.URL.Query().Get("excludeMe")

		var excludeID *uint

		if excludeMe == "true" {
			userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
			if !ok {
				c.logger.Warnw("User ID not found in context", "error", "User ID not found")
				lib.SendMessage(w, r, http.StatusUnauthorized, "User ID not found in context")
				return
			}

			excludeID = &userID
		}

		c.logger.Infow("Get all users request", "excludeMe", excludeMe, "excludeID", excludeID)

		users, err := c.userService.GetAllUsers(excludeID)
		if err != nil {
			c.logger.Warnw("Failed to retrieve users", "error", err.Error())
			lib.SendMessage(w, r, err.StatusCode, err.Error())
			return
		}

		c.logger.Infow("Successfully retrieved users")

		render.Status(r, http.StatusOK)
		render.JSON(w, r, users)
	}
}
