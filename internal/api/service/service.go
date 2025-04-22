package service

import (
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/api/service/friendship"
	"socialAPI/internal/api/service/user"
	"socialAPI/internal/shared"
)

type Service interface {
	Auth() auth.AuthService
	Token() shared.TokenService
	User() user.UserService
	Friendship() friendship.FriendshipService
}

type service struct {
	auth       auth.AuthService
	token      shared.TokenService
	user       user.UserService
	friendship friendship.FriendshipService
}

func NewService(a auth.AuthService, t shared.TokenService, u user.UserService, fr friendship.FriendshipService) Service {
	return &service{auth: a, token: t, user: u, friendship: fr}
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
