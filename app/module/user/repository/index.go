package repository

import (
	"errors"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/internal/bootstrap/database"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

type IRepository interface {
	GetAll(req request.Users) (users []*schema.User, paging paginator.Pagination, err error)
	GetOne(id uint64) (user *schema.User, err error)
	Create(user *schema.User) (err error)
	Update(id uint64, user *schema.User) (err error)
	Delete(id uint64) (err error)

	FindUserByMobile(mobile uint64) (user *schema.User, err error)

	GetUsers(req request.BusinessUsers) (users []*schema.User, paging paginator.Pagination, err error)
	GetPostObservers(postID uint64) (users []*schema.User, err error)
	GetUser(req request.BusinessUsersStoreRole) (user *schema.User, err error)
	InsertUser(businessID uint64, userID uint64) (err error)
	DeleteUser(businessID uint64, userID uint64) (err error)
}

func Repository(DB *database.Database) IRepository {
	return &repo{
		DB,
	}
}

type repo struct {
	DB *database.Database
}

func (_i *repo) GetAll(req request.Users) (users []*schema.User, paging paginator.Pagination, err error) {
	query := _i.DB.Main.Model(&schema.User{})

	if req.Keyword != "" {
		// Check if keyword is numeric for mobile search
		_, isNumeric := strconv.ParseUint(req.Keyword, 10, 64)
		if isNumeric == nil {
			query.Where("first_name Like ? OR last_name Like ? OR mobile = ?",
				"%"+req.Keyword+"%", "%"+req.Keyword+"%", req.Keyword)
		} else {
			query.Where("first_name Like ? OR last_name Like ?",
				"%"+req.Keyword+"%", "%"+req.Keyword+"%")
		}
	}

	if req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Order("created_at desc").Find(&users).Error
	if err != nil {
		return
	}

	paging = *req.Pagination

	return
}

func (_i *repo) GetOne(id uint64) (user *schema.User, err error) {
	if err := _i.DB.Main.First(&user, id).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (_i *repo) Create(user *schema.User) (err error) {
	return _i.DB.Main.Create(user).Error
}

func (_i *repo) Update(id uint64, user *schema.User) (err error) {
	return _i.DB.Main.Model(&schema.User{}).
		Where(&schema.User{ID: id}).
		Updates(user).Error
}

func (_i *repo) Delete(id uint64) error {
	return _i.DB.Main.Delete(&schema.User{}, id).Error
}

func (_i *repo) FindUserByMobile(mobile uint64) (user *schema.User, err error) {
	if err := _i.DB.Main.Where("mobile = ?", mobile).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("شماره تلفن همراه معتبر نمی باشد.")
		}
		return nil, err
	}

	return user, nil
}

func (_i *repo) GetUsers(req request.BusinessUsers) (users []*schema.User, paging paginator.Pagination, err error) {
	query := _i.DB.Main.
		Model(&[]schema.User{}).
		Select("users.*, COUNT(reservations.id) as reservation_count").
		Group("users.id").
		Order("users.created_at ASC")

	if req.CountUsing != 0 {
		query.Having("COUNT(reservations.id) >= ?", req.CountUsing)
	}

	if req.Role != "" {
		query.Where(`permissions -> '2' ? 'businessObserver'`)
	}

	// Build the LEFT JOIN condition dynamically with parameterized queries
	joinCondition := "LEFT JOIN reservations ON reservations.user_id = users.id AND reservations.deleted_at IS NULL AND reservations.status = 'reserved'"
	joinArgs := []interface{}{}

	if req.StartTime != nil && !req.StartTime.IsZero() {
		joinCondition += " AND reservations.start_time >= ?"
		joinArgs = append(joinArgs, utils.StartOfDayString(*req.StartTime))
	}

	if req.EndTime != nil && !req.EndTime.IsZero() {
		joinCondition += " AND reservations.end_time <= ?"
		joinArgs = append(joinArgs, utils.EndOfDayString(*req.EndTime))
	}

	query.Joins(joinCondition, joinArgs...)

	if len(req.UserIDs) > 0 {
		query.Where("users.id IN (?)", req.UserIDs)
	}

	if len(req.Username) > 0 {
		num, _ := strconv.ParseUint(req.Username, 10, 64)
		query.Where("users.mobile = ?", num)
	}

	if len(req.FullName) > 0 {
		query.Where(
			"CONCAT(users.first_name, ' ', users.last_name) LIKE ?",
			"%"+strings.TrimSpace(req.FullName)+"%",
		)
	}

	if req.CityID > 0 {
		query.Where("users.city_id = ?", req.CityID)
	}

	if req.WorkspaceID > 0 {
		query.Where("users.workspace_id = ?", req.WorkspaceID)
	}

	if req.DormitoryID > 0 {
		query.Where("users.dormitory_id = ?", req.DormitoryID)
	}

	if req.IsSuspended != "" {
		isSuspended := false
		if req.IsSuspended == "1" {
			isSuspended = true
		}
		query.Where("users.is_suspended = ?", isSuspended)
	}

	if req.Pagination != nil && req.Pagination.Page > 0 {
		var total int64
		query.Count(&total)
		req.Pagination.Total = total

		query.Offset(req.Pagination.Offset)
		query.Limit(req.Pagination.Limit)
	}

	err = query.Preload("Dormitory").
		Preload("Workspace").
		Preload("City").
		Find(&users).Error
	if err != nil {
		return
	}

	if req.Pagination != nil {
		paging = *req.Pagination
	}

	return
}

func (_i *repo) GetPostObservers(postID uint64) (users []*schema.User, err error) {
	err = _i.DB.Main.
		Model(&[]schema.User{}).
		Where(`meta->'PostsToObserve' @> ?::jsonb`, []uint64{postID}).
		Find(&users).Error
	if err != nil {
		// Handle error
	}

	return users, nil
}

func (_i *repo) GetUser(req request.BusinessUsersStoreRole) (user *schema.User, err error) {
	if err := _i.DB.Main.Exec(
		`SELECT FROM business_users WHERE user_id = ? AND business_id = ?`,
		req.UserID, req.BusinessID,
	).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (_i *repo) InsertUser(businessID uint64, userID uint64) (err error) {
	err = _i.DB.Main.
		Exec(`INSERT INTO business_users (user_id, business_id) VALUES (?, ?)`,
			userID, businessID,
		).Error

	if err != nil {
		return err
	}

	return nil
}

func (_i *repo) DeleteUser(businessID uint64, userID uint64) (err error) {
	err = _i.DB.Main.
		Exec(
			`DELETE FROM business_users WHERE user_id = ? AND business_id = ?`,
			userID, businessID,
		).Error

	if err != nil {
		return err
	}

	return nil
}
