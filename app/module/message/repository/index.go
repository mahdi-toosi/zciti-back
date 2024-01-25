package repository

import (
	"encoding/json"
	"fmt"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/message/request"
	"go-fiber-starter/internal/bootstrap"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.Messages) (messages []*schema.Message, paging paginator.Pagination, err error)
	CheckForNewMessages(roomId uint64) (messages []*schema.Message) //GetOne(id uint64) (message *schema.Message, err error)
	Create(msg *schema.Message) (message *schema.Message, err error)
	Update(id uint64, message *schema.Message) (err error)
	Delete(id uint64) (err error)
}

func Repository(DB *database.Database, Redis *bootstrap.Redis) IRepository {
	return &repo{
		DB:    DB,
		Redis: Redis,
	}
}

type repo struct {
	Redis *bootstrap.Redis
	DB    *database.Database
}

func (_i *repo) GetAll(req request.Messages) (messages []*schema.Message, paging paginator.Pagination, err error) {
	query := _i.DB.Chat.Model(&schema.Message{}).Where("room_id = ?", req.RoomID)

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&messages).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func generateRoomIDForRedis(id uint64) string {
	return fmt.Sprintf("RoomID=%d", id)
}

func (_i *repo) CheckForNewMessages(roomId uint64) (messages []*schema.Message) {
	data, err := _i.Redis.Storage.Get(generateRoomIDForRedis(roomId))
	if err != nil {
		return nil
	}
	if err := json.Unmarshal(data, &messages); err != nil {
		return nil
	}
	return
}

func (_i *repo) GetOne(id uint64) (message *schema.Message, err error) {
	if err := _i.DB.Chat.First(&message, id).Error; err != nil {
		return nil, err
	}

	return message, nil
}

func (_i *repo) Create(msg *schema.Message) (message *schema.Message, err error) {
	message = msg
	err = _i.DB.Chat.Create(&message).Error
	return message, err
}

func (_i *repo) Update(id uint64, message *schema.Message) (err error) {
	return _i.DB.Chat.Model(&schema.Message{}).
		Where(&schema.Message{ID: id}).
		Updates(message).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Chat.Delete(&schema.Message{}, id).Error
}
