package service

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/repository"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils/paginator"
	"slices"
)

type IService interface {
	Index(req request.Users) (users []*response.User, paging paginator.Pagination, err error)
	Show(id uint64) (user *response.User, err error)
	Store(req request.User) (err error)
	Update(id uint64, req request.User) (err error)
	UpdateAccount(req request.UpdateUserAccount) (err error)
	Destroy(id uint64) error

	GetPostObservers(postId uint64) (users []*response.User, err error)
	Users(req request.BusinessUsers) (users []*response.User, paging paginator.Pagination, err error)
	InsertUser(businessID uint64, userID uint64) (err error)
	DeleteUser(businessID uint64, userID uint64) (err error)
	BusinessUsersAddRole(req request.BusinessUsersStoreRole) error
	BusinessUsersTogglePostToObserve(userID uint64, postID uint64) (err error)
	BusinessUsersToggleSuspense(req request.BusinessUsersToggleSuspense) error
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

func (_i *service) GetPostObservers(postId uint64) (users []*response.User, err error) {
	results, err := _i.Repo.GetPostObservers(postId)
	if err != nil {
		return
	}

	for _, result := range results {
		users = append(users, response.FromDomain(result, nil))
	}

	return
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

func (_i *service) BusinessUsersTogglePostToObserve(userID uint64, postID uint64) (err error) {
	user, err := _i.Repo.GetOne(userID)
	if err != nil {
		return
	}
	if slices.Contains(user.Meta.PostsToObserve, postID) {
		index := slices.Index(user.Meta.PostsToObserve, postID)
		user.Meta.PostsToObserve = slices.Delete(user.Meta.PostsToObserve, index, 1)
	} else {
		user.Meta.PostsToObserve = append(user.Meta.PostsToObserve, postID)
	}

	err = _i.Repo.Update(userID, user)
	if err != nil {
		return err
	}

	return nil
}

func (_i *service) BusinessUsersAddRole(req request.BusinessUsersStoreRole) error {
	//_, err := _i.Repo.GetUser(req)
	//if err != nil {
	//	return &fiber.Error{
	//		Code:    fiber.StatusBadRequest,
	//		Message: "فقط برای کاربر خود میتوانید نقش مشخص کنید",
	//	}
	//}

	user, err := _i.Repo.GetOne(req.UserID)
	if err != nil {
		return err
	}

	user.Permissions[req.BusinessID] = req.Roles
	if user.Meta == nil {
		user.Meta = &schema.UserMeta{
			PostsToObserve:      req.PostsToObserve,
			TaxonomiesToObserve: req.TaxonomiesToObserve,
		}
	} else {
		user.Meta.PostsToObserve = req.PostsToObserve
		user.Meta.TaxonomiesToObserve = req.TaxonomiesToObserve
	}

	if err = _i.Repo.Update(req.UserID, user); err != nil {
		return err
	}

	return nil
}

func (_i *service) BusinessUsersToggleSuspense(req request.BusinessUsersToggleSuspense) error {
	user, err := _i.Repo.GetOne(req.UserID)
	if err != nil {
		return err
	}

	if req.IsSuspended {
		user.IsSuspended = &req.IsSuspended
	} else {
		isSuspended := false
		user.IsSuspended = &isSuspended
	}

	if len(req.SuspenseReason) == 0 {
		var suspenseReason string
		user.SuspenseReason = &suspenseReason
	} else {
		user.SuspenseReason = &req.SuspenseReason
	}

	if err = _i.Repo.Update(req.UserID, user); err != nil {
		return err
	}

	return nil
}
