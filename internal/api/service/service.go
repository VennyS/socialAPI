package service

import (
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/api/service/user"
	"socialAPI/internal/shared"
)

type Service interface {
	Auth() auth.AuthService
	Token() shared.TokenService
	User() user.UserService
}

type service struct {
	auth  auth.AuthService
	token shared.TokenService
	user  user.UserService
}

func NewService(a auth.AuthService, t shared.TokenService, u user.UserService) Service {
	return &service{auth: a, token: t, user: u}
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
