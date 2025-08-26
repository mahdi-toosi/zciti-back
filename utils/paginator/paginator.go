package paginator

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/utils"
)

const (
	defaultLimit = 10
)

type Pagination struct {
	Limit  int   `json:"itemPerPage,omitempty"`
	Offset int   `json:"-"`
	Page   int   `json:"page,omitempty"`
	Total  int64 `json:"total,omitempty"`
}

func Paginate(c *fiber.Ctx) (*Pagination, error) {
	limit, err := utils.GetUintInQueries(c, "itemPerPage")
	if err != nil {
		limit = defaultLimit
	}
	page, err := utils.GetUintInQueries(c, "page")
	if err != nil {
		page = 0
		limit = 0
	}
	p := &Pagination{
		Limit: int(limit),
		Page:  int(page),
	}
	if p.Page > 0 {
		p.Offset = (p.Page - 1) * p.Limit
	}

	return p, nil
}
