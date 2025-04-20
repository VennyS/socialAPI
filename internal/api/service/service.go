package service

import (
	"socialAPI/internal/api/service/auth"
	"socialAPI/internal/shared"
)

type Service interface {
	Auth() auth.AuthService
	Token() shared.TokenService
}

type service struct {
	auth  auth.AuthService
	token shared.TokenService
}

func NewService(a auth.AuthService, t shared.TokenService) Service {
	return &service{auth: a, token: t}
}

func (s service) Auth() auth.AuthService {
	return s.auth
}

func (s service) Token() shared.TokenService {
	return s.token
}
