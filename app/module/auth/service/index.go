package service

import (
	"errors"

	"go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/auth/request"
	"go-fiber-starter/app/module/auth/response"
	"go-fiber-starter/app/module/user/repository"
)

//go:generate mockgen -destination=article_service_mock.go -package=service . AuthService
type IService interface {
	Login(req request.LoginRequest) (res response.LoginResponse, err error)
}

func Service(Repo repository.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Login(req request.LoginRequest) (res response.LoginResponse, err error) {
	// check user by username
	user, err := _i.Repo.FindUserByUsername(req.Username)
	if err != nil {
		return
	}

	if user == nil {
		err = errors.New("user not found")
		return
	}

	// check password
	if !user.ComparePassword(req.Password) {
		err = errors.New("password not match")
		return
	}

	// do create token
	claims, err := middleware.GenerateTokenAccess(user.ID)
	if err != nil {
		return
	}

	res.Token = claims.Token
	res.Type = claims.Type
	res.ExpiresAt = claims.ExpiresAt

	return
}
