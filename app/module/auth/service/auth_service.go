package service

import (
	"errors"

	"github.com/bangadam/go-fiber-starter/app/middleware"
	"github.com/bangadam/go-fiber-starter/app/module/auth/request"
	"github.com/bangadam/go-fiber-starter/app/module/auth/response"
	user_repo "github.com/bangadam/go-fiber-starter/app/module/user/repository"
)

//go:generate mockgen -destination=article_service_mock.go -package=service . AuthService
type IService interface {
	Login(req request.LoginRequest) (res response.LoginResponse, err error)
}

func Service(userRepo user_repo.IRepository) IService {
	return &service{
		userRepo,
	}
}

type service struct {
	userRepo user_repo.IRepository
}

func (_i *service) Login(req request.LoginRequest) (res response.LoginResponse, err error) {
	// check user by username
	user, err := _i.userRepo.FindUserByUsername(req.Username)
	if err != nil {
		return
	}

	if user == nil {
		err = errors.New("user not found")
		return
	}

	// password not null
	if user.Password == nil {
		err = errors.New("user password not found")
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
