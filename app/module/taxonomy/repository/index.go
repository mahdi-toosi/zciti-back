package repository

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/taxonomy/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.Taxonomies) (taxonomies []*schema.Taxonomy, paging paginator.Pagination, err error)
	Search(req request.Taxonomies) (taxonomies []*schema.Taxonomy, paging paginator.Pagination, err error)
	GetOne(BusinessID uint64, id uint64) (taxonomy *schema.Taxonomy, err error)
	Create(taxonomy *schema.Taxonomy) (err error)
	Update(id uint64, taxonomy *schema.Taxonomy) (err error)
	Delete(BusinessID uint64, id *uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Taxonomies) (taxonomies []*schema.Taxonomy, paging paginator.Pagination, err error) {
	if req.Pagination.Page > 0 {
		var total int64
		err := _i.DB.Main.Model(&schema.Taxonomy{}).
			Where(&schema.Taxonomy{BusinessID: req.BusinessID, Type: req.Type, ParentID: nil}).
			Count(&total).Error
		if err != nil {
			return nil, paginator.Pagination{}, err
		}
		req.Pagination.Total = total
	}

	q := fmt.Sprintf(`
		WITH roots AS (
			SELECT id FROM taxonomies WHERE parent_id IS NULL
				AND business_id = ? AND %s type = ? AND %s deleted_at IS NULL 
				OFFSET ? LIMIT ?
		),
			 recursive AS (
				 WITH RECURSIVE taxonomy_tree(id, parent_id, title, type, domain, slug, description, created_at, depth) AS (
					 SELECT c.id, c.parent_id, c.title, c.type, c.domain, c.slug, c.description, c.created_at, 1
					 	FROM taxonomies c JOIN roots ON roots.id = c.id
						WHERE deleted_at IS NULL
					 UNION ALL
					 SELECT c.id, c.parent_id, c.title, c.type, c.domain, c.slug, c.description, c.created_at, p.depth + 1
					 	FROM taxonomies c JOIN taxonomy_tree p ON c.parent_id = p.id
						WHERE deleted_at IS NULL
				 )
				 SELECT * FROM taxonomy_tree
			 )
		SELECT r.* FROM recursive r ORDER BY r.created_at DESC;
		`,
		utils.InlineCondition(req.Domain != "", "domain = '"+req.Domain+"' AND", ""),
		utils.InlineCondition(req.Keyword != "", "title Like '%"+req.Keyword+"%' AND", ""),
	)

	query := _i.DB.Main.Raw(q, req.BusinessID, req.Type, req.Pagination.Offset, req.Pagination.Limit)

	err = query.Scan(&taxonomies).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) Search(req request.Taxonomies) (taxonomies []*schema.Taxonomy, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(schema.Taxonomy{}).
		Where(schema.Taxonomy{BusinessID: req.BusinessID})

	// TODO what if we doesnt want to add parent_id in wheres ?!!!!!!!
	if req.ParentID == 0 {
		query.Where("parent_id IS NULL")
	} else if req.ParentID > 0 {
		query.Where("parent_id = ?", req.ParentID)
	}

	if req.Keyword != "" {
		query.Where("title Like ?", "%"+req.Keyword+"%")
	}

	if req.Type != "" {
		query.Where("type = ?", req.Type)
	}

	if req.Domain != "" {
		query.Where("domain = ?", req.Domain)
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&taxonomies).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (taxonomy *schema.Taxonomy, err error) {
	if err := _i.DB.Main.
		Where(&schema.Taxonomy{BusinessID: businessID}).
		First(&taxonomy, id).Error; err != nil {
		return nil, err
	}

	return taxonomy, nil
}

func (_i *repo) Create(taxonomy *schema.Taxonomy) (err error) {
	return _i.DB.Main.Create(taxonomy).Error
}

func (_i *repo) Update(id uint64, taxonomy *schema.Taxonomy) (err error) {
	return _i.DB.Main.Model(&schema.Taxonomy{}).
		Where(&schema.Taxonomy{ID: id, BusinessID: taxonomy.BusinessID}).
		Updates(taxonomy).Error
}

func (_i *repo) Delete(businessID uint64, id *uint64) error {
	// check for child
	if err := _i.DB.Main.
		Where(&schema.Taxonomy{ParentID: id, BusinessID: businessID}).
		First(&schema.Taxonomy{}).
		Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		// does not have child
		return _i.DB.Main.
			Where(&schema.Taxonomy{ID: *id, BusinessID: businessID}).
			Delete(&schema.Taxonomy{ID: *id}).Error
	}

	return fiber.ErrFailedDependency

}
