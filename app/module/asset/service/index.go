package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/asset/repository"
	"go-fiber-starter/app/module/asset/request"
	"go-fiber-starter/app/module/asset/response"
	businessRepo "go-fiber-starter/app/module/business/repository"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/paginator"
)

type IService interface {
	Index(req request.Assets) (assets []*response.Asset, paging paginator.Pagination, err error)
	Show(id uuid.UUID) (asset *response.Asset, err error)
	Store(req request.Asset) (err error)
	Update(id uuid.UUID, req request.Asset) (err error)
	Destroy(user schema.User, id uuid.UUID) error
}

func Service(Repo repository.IRepository, BusinessRepo businessRepo.IRepository, config *config.Config) IService {
	return &service{
		Repo,
		config,
		BusinessRepo,
	}
}

type service struct {
	Repo         repository.IRepository
	config       *config.Config
	BusinessRepo businessRepo.IRepository
}

func (_i *service) Index(req request.Assets) (assets []*response.Asset, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		assets = append(assets, response.FromDomain(result, _i.config.ParseAddress()))
	}

	return
}

func (_i *service) Show(id uuid.UUID) (asset *response.Asset, err error) {
	result, err := _i.Repo.GetOne(id)
	if err != nil {
		return nil, err
	}

	return response.FromDomain(result, _i.config.ParseAddress()), nil
}

func (_i *service) Store(req request.Asset) (err error) {
	// TODO => add asset size to business meta
	return _i.Repo.Create(req.ToDomain())
}

func (_i *service) Update(id uuid.UUID, req request.Asset) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(user schema.User, id uuid.UUID) error {
	asset, err := _i.Repo.GetOne(id)
	if err != nil {
		return err
	}

	business, err := _i.BusinessRepo.GetOne(asset.BusinessID)
	if err != nil {
		return err
	}

	if business.OwnerID != user.ID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	return _i.Repo.Delete(id)
}
