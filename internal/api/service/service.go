package service

import (
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/api/service/chat"
	"socialAPI/internal/api/service/friendship"
	"socialAPI/internal/api/service/user"
	"socialAPI/internal/shared"
)

type Service interface {
	Auth() auth.AuthService
	Token() shared.TokenService
	User() user.UserService
	Friendship() friendship.FriendshipService
	Chat() chat.ChatService
}

type service struct {
	auth       auth.AuthService
	token      shared.TokenService
	user       user.UserService
	friendship friendship.FriendshipService
	chat       chat.ChatService
}

func NewService(a auth.AuthService, t shared.TokenService, u user.UserService, fr friendship.FriendshipService, c chat.ChatService) Service {
	return &service{auth: a, token: t, user: u, friendship: fr, chat: c}
}

func (s service) Auth() auth.AuthService {
	return s.auth
}

func (s service) Token() shared.TokenService {
	return s.token
}

func (s service) User() user.UserService {
	return s.user
}

func (s service) Friendship() friendship.FriendshipService {
	return s.friendship
}

func (s service) Chat() chat.ChatService {
	return s.chat
}
