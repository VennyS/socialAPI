package chat

type CreateRequest struct {
	UserIDs []uint  `json:"userIDs" validate:"required,min=2"`
	Name    *string `json:"name"`
}
