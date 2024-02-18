package request

import (
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"mime/multipart"
)

type Asset struct {
	ID         uuid.UUID
	Asset      multipart.FileHeader `example:"file" validate:"required,file"`
	Title      string
	Path       string
	Ext        string
	IsPrivate  bool `example:"true" validate:"boolean"`
	UserID     uint64
	BusinessID uint64
}

type UpdateAsset struct {
	ID    uuid.UUID
	Title string `gorm:"not null"`
}

type Assets struct {
	Pagination *paginator.Pagination
	Keyword    string
}

func (req *Asset) ToDomain() *schema.Asset {
	return &schema.Asset{
		ID:         req.ID,
		Title:      req.Title,
		Path:       req.Path,
		Ext:        req.Ext,
		IsPrivate:  req.IsPrivate,
		UserID:     req.UserID,
		BusinessID: req.BusinessID,
	}
}
