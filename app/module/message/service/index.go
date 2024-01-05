package service

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/message/repository"
	"go-fiber-starter/app/module/message/request"
	"go-fiber-starter/app/module/message/response"
	messageRoomRepo "go-fiber-starter/app/module/messageRoom/repository"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Messages) (messages []*response.Message, paging paginator.Pagination, msgRoom *schema.MessageRoom, err error)
	//Show(id uint64) (message *response.Message, err error)
	Store(req request.Message) (err error)
	Update(id uint64, req request.Message) (err error)
	Destroy(id uint64) error
}

func Service(Repo repository.IRepository, messageRoomRepo messageRoomRepo.IRepository) IService {
	return &service{
		Repo:            Repo,
		messageRoomRepo: messageRoomRepo,
	}
}

type service struct {
	Repo            repository.IRepository
	messageRoomRepo messageRoomRepo.IRepository
}

func (_i *service) Index(req request.Messages) (messages []*response.Message, paging paginator.Pagination, msgRoom *schema.MessageRoom, err error) {
	msgRoom, err = _i.messageRoomRepo.GetOne(req.BusinessID, req.UserID)
	if err != nil {
		return
	}

	req.RoomID = msgRoom.ID
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		messages = append(messages, response.FromDomain(result))
	}

	return
}

//func (_i *service) Show(id uint64) (message *response.Message, err error) {
//	result, err := _i.Repo.GetOne(id)
//	if err != nil {
//		return nil, err
//	}
//
//	return response.FromDomain(result), nil
//}

func (_i *service) Store(req request.Message) (err error) {
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uint64, req request.Message) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}
