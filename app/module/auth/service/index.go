package service

import (
	"errors"
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/auth/request"
	"go-fiber-starter/app/module/auth/response"
	usersRepo "go-fiber-starter/app/module/user/repository"
	userResponse "go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils/helpers"
)

//go:generate mockgen -destination=article_service_mock.go -package=service . AuthService
type IService interface {
	Login(req request.Login) (res response.Login, err error)
	Register(req *request.Register) (res response.Register, err error)
}

func Service(Repo usersRepo.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo usersRepo.IRepository
}

func (_i *service) Login(req request.Login) (res response.Login, err error) {
	// check user by username
	user, err := _i.Repo.FindUserByMobile(req.Mobile)
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
	token, err := middleware.GenerateTokenAccess(*userResponse.FromDomain(user))
	if err != nil {
		return
	}

	res.User = *userResponse.FromDomain(user)
	res.Token = token

	return
}

func (_i *service) Register(req *request.Register) (res response.Register, err error) {
	// check user by username
	user := &schema.User{
		Mobile:    req.Mobile,
		LastName:  req.LastName,
		FirstName: req.FirstName,
		Roles:     []string{"user"},
		Password:  helpers.Hash([]byte(req.Password)),
	}

	err = _i.Repo.Create(user)
	if err != nil {
		return response.Register{}, err
	}

	log.Debug().Msgf("%+v", user)

	// do create token
	token, err := middleware.GenerateTokenAccess(*userResponse.FromDomain(user))
	if err != nil {
		return
	}

	res.User = *userResponse.FromDomain(user)
	res.Token = token

	return
}
