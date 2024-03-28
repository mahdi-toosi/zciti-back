package response

import (
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
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

func FromDomain(asset *schema.Asset, domain string) (res *Asset) {
	if asset == nil {
		return res
	}

	a := &Asset{
		ID:            asset.ID,
		Ext:           asset.Ext,
		Path:          asset.Path,
		Size:          asset.Size,
		Title:         asset.Title,
		UserID:        asset.UserID,
		IsPrivate:     asset.IsPrivate,
		BusinessID:    asset.BusinessID,
		BusinessTitle: asset.Business.Title,
		UserFullName:  asset.User.FirstName + " " + asset.User.LastName,
	}

	if strings.Contains(a.Path, "/private/") {
		a.Path = ""
	} else {
		a.Path = strings.ReplaceAll(a.Path, "storage/public", domain+"/asset")
	}

	return a
}
