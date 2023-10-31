package paginator

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

const (
	defaultLimit = 10
)

type Pagination struct {
	Limit  int   `json:"itemPerPage,omitempty"`
	Offset int   `json:"-"`
	Page   int   `json:"-"`
	Total  int64 `json:"total,omitempty"`
}

func Paginate(c *fiber.Ctx) (*Pagination, error) {
	limit, err := strconv.Atoi(c.Query("itemPerPage"))
	if err != nil {
		limit = defaultLimit
	}
	page, err := strconv.Atoi(c.Query("page"))
	if err != nil {
		page = 0
		limit = 0
	}
	p := &Pagination{
		Limit: limit,
		Page:  page,
	}
	if p.Page > 0 {
		p.Offset = (p.Page - 1) * p.Limit
	}

	return p, nil
}
