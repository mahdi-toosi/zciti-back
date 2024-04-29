package response

import (
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/paginator"
	"strings"
)

type Asset struct {
	ID            uuid.UUID `json:",omitempty"`
	Title         string    `json:",omitempty"`
	Path          string    `json:",omitempty"`
	Ext           string    `json:",omitempty"`
	IsPrivate     bool      `json:",omitempty"`
	Size          uint64    `json:",omitempty"`
	UserID        uint64    `json:",omitempty"`
	BusinessID    uint64    `json:",omitempty"`
	BusinessTitle string    `json:",omitempty"`
	UserFullName  string    `json:",omitempty"`
}

type Assets struct {
	Data       []*Asset `json:",omitempty"`
	AssetsSize uint64
	Meta       paginator.Pagination `json:",omitempty"`
}

func FromDomain(item *schema.Asset, domain string) (res *Asset) {
	if item == nil {
		return res
	}

	a := &Asset{
		ID:            item.ID,
		Ext:           item.Ext,
		Path:          item.Path,
		Size:          item.Size,
		Title:         item.Title,
		UserID:        item.UserID,
		IsPrivate:     item.IsPrivate,
		BusinessID:    item.BusinessID,
		BusinessTitle: item.Business.Title,
		UserFullName:  item.User.FullName(),
	}

	if strings.Contains(a.Path, "/private/") {
		a.Path = ""
	} else {
		a.Path = strings.ReplaceAll(a.Path, "storage/public", domain+"/asset")
	}

	return a
}
