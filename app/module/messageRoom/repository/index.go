package repository

import (
	"errors"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/messageRoom/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.MessageRooms) (messageRooms []*schema.MessageRoom, paging paginator.Pagination, err error)
	Delete(id uint64) (err error)
	GetOne(businessID uint64, userID uint64) (messageRoom *schema.MessageRoom, err error)
	Create(messageRoom *schema.MessageRoom) (err error)
	Update(id uint64, messageRoom *schema.MessageRoom) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.MessageRooms) (messageRooms []*schema.MessageRoom, paging paginator.Pagination, err error) {
	query := _i.DB.Chat.Model(&schema.MessageRoom{})

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at asc").Find(&messageRooms).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, userID uint64) (messageRoom *schema.MessageRoom, err error) {
	err = _i.DB.Chat.
		First(&messageRoom, "business_id = ? AND user_id = ?", businessID, userID).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			mr := schema.MessageRoom{
				UserID:     userID,
				Status:     "active",
				BusinessID: businessID,
			}

			if err := _i.Create(&mr); err != nil {
				return nil, err
			}

			return &mr, nil
		}

		return nil, err
	}

	return messageRoom, nil
}

func (_i *repo) Create(messageRoom *schema.MessageRoom) (err error) {
	return _i.DB.Chat.Create(&messageRoom).Error
}

func (_i *repo) Update(id uint64, messageRoom *schema.MessageRoom) (err error) {
	return _i.DB.Chat.Model(&schema.MessageRoom{}).
		Where(&schema.MessageRoom{ID: id}).
		Updates(messageRoom).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Chat.Delete(&schema.MessageRoom{}, id).Error
}
