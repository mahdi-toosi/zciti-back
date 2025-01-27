package service

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Users) (users []*response.User, paging paginator.Pagination, err error)
	Show(id uint64) (user *response.User, err error)
	Store(req request.User) (err error)
	Update(id uint64, req request.User) (err error)
	UpdateAccount(req request.UpdateUserAccount) (err error)
	Destroy(id uint64) error

	Users(req request.BusinessUsers) (users []*response.User, paging paginator.Pagination, err error)
	InsertUser(businessID uint64, userID uint64) (err error)
	DeleteUser(businessID uint64, userID uint64) (err error)
	BusinessUsersAddRole(req request.BusinessUsersStoreRole) error
}

func Service(repo repository.IRepository) IService {
	return &service{repo}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.Users) (users []*response.User, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		users = append(users, response.FromDomain(result, nil))
	}

	return
}

func (_i *service) Show(id uint64) (user *response.User, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result, nil), nil
}

func (_i *service) Store(req request.User) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.User) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) UpdateAccount(req request.UpdateUserAccount) (err error) {
	u := req.ToDomain()
	return _i.Repo.Update(u.ID, u)
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}

func (_i *service) Users(req request.BusinessUsers) (users []*response.User, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetUsers(req)
	if err != nil {
		return
	}

	for _, result := range results {
		users = append(users, response.FromDomain(result, &req.BusinessID))
	}

	return
}

func (_i *service) InsertUser(businessID uint64, userID uint64) (err error) {
	err = _i.Repo.InsertUser(businessID, userID)
	if err != nil {
		return
	}

	return nil
}

func (_i *service) DeleteUser(businessID uint64, userID uint64) (err error) {
	err = _i.Repo.DeleteUser(businessID, userID)
	if err != nil {
		return
	}

	return nil
}
func (_i *service) BusinessUsersAddRole(req request.BusinessUsersStoreRole) error {
	_, err := _i.Repo.GetUser(req)
	if err != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "فقط برای کاربر خود میتوانید نقش مشخص کنید",
		}
	}

	user, err := _i.Repo.GetOne(req.UserID)
	if err != nil {
		return err
	}

	user.Permissions[req.BusinessID] = req.Roles

	if err = _i.Repo.Update(req.UserID, user); err != nil {
		return err
	}

	return nil
}
