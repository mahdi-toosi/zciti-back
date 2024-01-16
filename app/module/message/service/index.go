package service

import (
	"github.com/rs/zerolog/log"
	"go-fiber-starter/app/database/schema"
	businessRepo "go-fiber-starter/app/module/business/repository"
	"go-fiber-starter/app/module/message/repository"
	"go-fiber-starter/app/module/message/request"
	"go-fiber-starter/app/module/message/response"
	messageRoomRepo "go-fiber-starter/app/module/messageRoom/repository"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Messages, roomID *uint64) (messages []*response.Message, paging paginator.Pagination, msgRoom *schema.MessageRoom, business *schema.Business, err error)
	Store(req request.Message) (err error)
	Update(id uint64, req request.Message) (err error)
	Destroy(id uint64) error
	//Show(id uint64) (message *response.Message, err error)
}

func Service(Repo repository.IRepository, messageRoomRepo messageRoomRepo.IRepository, businessRepo businessRepo.IRepository) IService {
	return &service{
		Repo:            Repo,
		BusinessRepo:    businessRepo,
		MessageRoomRepo: messageRoomRepo,
	}
}

type service struct {
	Repo            repository.IRepository
	BusinessRepo    businessRepo.IRepository
	MessageRoomRepo messageRoomRepo.IRepository
}

func (_i *service) Index(req request.Messages, roomID *uint64) (messages []*response.Message, paging paginator.Pagination, msgRoom *schema.MessageRoom, business *schema.Business, err error) {
	if roomID == nil {
		msgRoom, err = _i.MessageRoomRepo.GetOne(req.BusinessID, req.UserID)
		if err != nil {
			return
		}

		req.RoomID = &msgRoom.ID
		business, err = _i.BusinessRepo.GetOne(msgRoom.BusinessID)
		if err != nil {
			return
		}
		log.Debug().Msgf("msgRoom => %+v", msgRoom)
	} else {
		req.RoomID = roomID
	}

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
