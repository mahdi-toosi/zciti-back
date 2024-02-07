package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/comment/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(postID uint64, req request.Comments) (comments []map[string]interface{}, paging paginator.Pagination, err error)
	GetOne(id uint64) (comment *schema.Comment, err error)
	Create(comment *schema.Comment) (err error)
	Update(id uint64, comment *schema.Comment) (err error)
	Delete(id uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(postID uint64, req request.Comments) (comments []map[string]interface{}, paging paginator.Pagination, err error) {
	if req.Pagination.Page > 0 {
		var total int64
		_i.DB.Main.Model(&schema.Comment{}).Where("post_id = ? and parent_id is null", postID).Count(&total)
		req.Pagination.Total = total
	}

	err = _i.DB.Main.Raw(`
		WITH roots AS (
			SELECT id FROM comments WHERE parent_id IS NULL AND post_id = ? offset ? limit ?
		),
			 recursive AS (
				 WITH RECURSIVE comment_tree(id, parent_id, content, author_id, status, is_business_owner, created_at, depth) AS (
					 SELECT c.id, c.parent_id, c.content, c.author_id, c.status, c.is_business_owner, c.created_at, 1
					 FROM comments c JOIN roots ON roots.id = c.id
		
					 UNION ALL
		
					 SELECT c.id, c.parent_id, c.content, c.author_id, c.status, c.is_business_owner, c.created_at, p.depth + 1
					 FROM comments c JOIN comment_tree p ON c.parent_id = p.id
				 )
				 SELECT * FROM comment_tree
			 )
		SELECT r.*, u.first_name || ' ' || u.last_name as author_full_name FROM recursive r
		left join users u on u.id = r.author_id
		order by r.created_at desc;
		`, postID, req.Pagination.Offset, req.Pagination.Limit).Scan(&comments).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (comment *schema.Comment, err error) {
	if err := _i.DB.Main.First(&comment, id).Error; err != nil {
		return nil, err
	}

	return comment, nil
}

func (_i *repo) Create(comment *schema.Comment) (err error) {
	return _i.DB.Main.Create(comment).Error
}

func (_i *repo) Update(id uint64, comment *schema.Comment) (err error) {
	return _i.DB.Main.Model(&schema.Comment{}).
		Where(&schema.Comment{ID: id}).
		Updates(comment).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.Comment{}, id).Error
}
