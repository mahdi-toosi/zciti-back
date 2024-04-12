package service

import (
	"errors"
	"fmt"
	MessageWay "github.com/MessageWay/MessageWayGolang"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/auth/request"
	"go-fiber-starter/app/module/auth/response"
	usersRepo "go-fiber-starter/app/module/user/repository"
	userResponse "go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/helpers"
)

//go:generate mockgen -destination=article_service_mock.go -package=service . AuthService
type IService interface {
	Login(req request.Login, jwtConfig config.Jwt) (res response.Login, err error)
	Register(req *request.Register, jwtConfig config.Jwt) (res response.Register, err error)
	SendOtp(req *request.SendOtp) error
	ResetPass(req *request.ResetPass) error
}

func Service(Repo usersRepo.IRepository, MessageWay *MessageWay.App) IService {
	return &service{
		Repo,
		MessageWay,
	}
}

type service struct {
	Repo       usersRepo.IRepository
	MessageWay *MessageWay.App
}

func (_i *service) Login(req request.Login, jwtConfig config.Jwt) (res response.Login, err error) {
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
	token, err := middleware.GenerateTokenAccess(*user, jwtConfig)
	if err != nil {
		return
	}

	res.User = *userResponse.FromDomain(user)
	res.Token = token

	return
}

func (_i *service) Register(req *request.Register, jwtConfig config.Jwt) (res response.Register, err error) {
	// check user by username
	user := &schema.User{
		Mobile:      req.Mobile,
		LastName:    req.LastName,
		FirstName:   req.FirstName,
		Permissions: schema.UserPermissionsMap{},
		Password:    helpers.Hash([]byte(req.Password)),
	}

	err = _i.Repo.Create(user)
	if err != nil {
		return response.Register{}, err
	}

	// do create token
	token, err := middleware.GenerateTokenAccess(*user, jwtConfig)
	if err != nil {
		return
	}

	res.User = *userResponse.FromDomain(user)
	res.Token = token

	return
}

func (_i *service) SendOtp(req *request.SendOtp) error {
	user, err := _i.Repo.FindUserByMobile(req.Mobile)
	if err != nil {
		return err
	}

	_, err = _i.MessageWay.Send(MessageWay.Message{
		TemplateID: 3,
		Method:     "sms",
		Mobile:     fmt.Sprint("0", user.Mobile),
	})

	if err != nil {
		return err
	}
	return nil
}

func (_i *service) ResetPass(req *request.ResetPass) error {
	user, err := _i.Repo.FindUserByMobile(req.Mobile)
	if err != nil {
		return err
	}

	_, err = _i.MessageWay.Verify(MessageWay.OTPVerifyRequest{
		OTP:    req.Otp,
		Mobile: fmt.Sprint("0", user.Mobile),
	})

	if err != nil {
		return err
	}

	err = _i.Repo.Update(user.ID, &schema.User{Password: helpers.Hash([]byte(req.Password))})
	if err != nil {
		return err
	}

	return nil
}
