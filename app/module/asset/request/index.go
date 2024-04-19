package request

import (
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"mime/multipart"
)

type Asset struct {
	ID           uuid.UUID
	Asset        multipart.FileHeader `example:"file" validate:"required,validate_file=.pdf:10 .doc:10 .docx:10 .jpeg:1 .jpg:1 .png:1"`
	Title        string               `example:"title" validate:"omitempty,min=2,max=255"`
	Path         string               ``
	Ext          string               ``
	AlsoOptimize bool                 `example:"title" validate:"omitempty,boolean"`
	Size         uint64               ``
	IsPrivate    bool                 `example:"true" validate:"boolean"`
	UserID       uint64               ``
	BusinessID   uint64               ``
}

type UpdateAsset struct {
	ID    uuid.UUID
	Title string `gorm:"not null"`
}

type Assets struct {
	BusinessID uint64
	Keyword    string
	Pagination *paginator.Pagination
}

func (req *Asset) ToDomain() *schema.Asset {
	return &schema.Asset{
		ID:         req.ID,
		Title:      req.Title,
		Path:       req.Path,
		Ext:        req.Ext,
		Size:       req.Size,
		IsPrivate:  req.IsPrivate,
		UserID:     req.UserID,
		BusinessID: req.BusinessID,
	}
}
