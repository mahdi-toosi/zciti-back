package repository

import (
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/post/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils/paginator"
)

type IRepository interface {
	GetAll(req request.PostsRequest) (posts []*schema.Post, paging paginator.Pagination, err error)
	GetOne(businessID uint64, id uint64) (post *schema.Post, err error)
	Create(post *schema.Post) (result *schema.Post, err error)
	Update(id uint64, post *schema.Post) error
	UpdateCommentCount(id uint64, num string) (err error)
	Delete(businessID uint64, id uint64) error
	DeleteTaxonomies(req request.PostTaxonomies) error
	InsertTaxonomies(req request.PostTaxonomies) error
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.PostsRequest) (posts []*schema.Post, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.Post{}).
		Where("business_id = ?", req.BusinessID).
		Where("type = ?", schema.PostTypePost)

	if req.Keyword != "" {
		query.Where("title Like ?", "%"+req.Keyword+"%")
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.
		Preload("Author").
		Preload("Business").
		Preload("Taxonomies").
		Order("created_at desc").Find(&posts).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(businessID uint64, id uint64) (post *schema.Post, err error) {
	err = _i.DB.Main.
		Preload("Author").
		Preload("Business").
		Preload("Taxonomies").
		Where("business_id = ?", businessID).
		Where("type = ?", schema.PostTypePost).
		First(&post, id).Error
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (_i *repo) Create(post *schema.Post) (result *schema.Post, err error) {
	err = _i.DB.Main.Create(&post).Error
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (_i *repo) Update(id uint64, post *schema.Post) error {
	return _i.DB.Main.Model(&schema.Post{}).
		Where(&schema.Post{ID: id}).
		Updates(post).Error
}

func (_i *repo) UpdateCommentCount(id uint64, num string) (err error) {
	return _i.DB.Main.Model(&schema.Comment{}).
		Where(&schema.Comment{ID: id}).
		Update("meta", _i.DB.Main.Raw("jsonb_set(meta, '{CommentsCount}', (meta->'CommentsCount')::int ?)", num)).
		Error
}

func (_i *repo) Delete(businessID uint64, id uint64) error {
	return _i.DB.Main.
		Where(&schema.Post{BusinessID: businessID}).
		Delete(&schema.Post{}, id).Error
}

func (_i *repo) DeleteTaxonomies(req request.PostTaxonomies) error {
	return _i.DB.Main.Exec(
		"DELETE FROM posts_taxonomies WHERE post_id = ? AND taxonomy_id IN (?)",
		req.PostID,
		req.IDs,
	).Error
}

func (_i *repo) InsertTaxonomies(req request.PostTaxonomies) error {
	for _, taxonomyID := range req.IDs {
		_i.DB.Main.Exec(
			"INSERT INTO posts_taxonomies (post_id, taxonomy_id) VALUES (?, ?)",
			req.PostID,
			taxonomyID,
		)
	}
	return nil
}
