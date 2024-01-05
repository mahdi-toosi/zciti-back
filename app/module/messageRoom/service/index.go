package service

import (
	"go-fiber-starter/app/module/messageRoom/repository"
	"go-fiber-starter/app/module/messageRoom/request"
	"go-fiber-starter/app/module/messageRoom/response"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.MessageRooms) (messageRooms []*response.MessageRoom, paging paginator.Pagination, err error)
	Destroy(id uint64) error
	Update(id uint64, req request.MessageRoom) (err error)
	//Show(id uint64) (messageRoom *response.MessageRoom, err error)
	// Store(req request.MessageRoom) (err error)
}

func Service(Repo repository.IRepository) IService {
	return &service{
		Repo,
	}
}

type service struct {
	Repo repository.IRepository
}

func (_i *service) Index(req request.MessageRooms) (messageRooms []*response.MessageRoom, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		messageRooms = append(messageRooms, response.FromDomain(result))
	}

	return
}

//func (_i *service) Show(id uint64) (messageRoom *response.MessageRoom, err error) {
//	result, err := _i.Repo.GetOne(id)
//	if err != nil {
//		return nil, err
//	}
//
//	return response.FromDomain(result), nil
//}

// func (_i *service) Store(req request.MessageRoom) (err error) {
//	return _i.Repo.Create(req.ToDomain())
// }

func (_i *service) Update(id uint64, req request.MessageRoom) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
