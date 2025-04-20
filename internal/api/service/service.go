package service

import "socialAPI/internal/api/service/auth"

type Service interface {
	Auth() auth.AuthService
}

type service struct {
	auth auth.AuthService
}

func NewService(a auth.AuthService) Service {
	return &service{auth: a}
}

func (s service) Auth() auth.AuthService {
	return s.auth
}
