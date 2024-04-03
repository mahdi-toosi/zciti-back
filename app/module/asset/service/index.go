package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/middleware"
	"go-fiber-starter/app/module/asset/repository"
	"go-fiber-starter/app/module/asset/request"
	"go-fiber-starter/app/module/asset/response"
	businessRepo "go-fiber-starter/app/module/business/repository"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/paginator"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type IService interface {
	Index(req request.Assets) (assets []*response.Asset, assetsSize uint64, paging paginator.Pagination, err error)
	Store(ctx *fiber.Ctx, req request.Asset) (err error)
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

func (_i *service) Index(req request.Assets) (assets []*response.Asset, assetsSize uint64, paging paginator.Pagination, err error) {
	results, assetsSize, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		assets = append(assets, response.FromDomain(result, _i.config.ParseAddress()))
	}

	return
}

func (_i *service) Store(c *fiber.Ctx, req request.Asset) (err error) {
	business, err := _i.BusinessRepo.GetOne(req.BusinessID)
	if err != nil {
		return err
	}

	accountAssetLimit := middleware.Accounts[business.Account].AssetsSizeLimit
	business.AssetsSize += uint64(req.Asset.Size)

	if business.AssetsSize > accountAssetLimit {
		return fiber.ErrForbidden
	}

	var folder string
	if req.IsPrivate {
		folder = "private"
	} else {
		folder = "public"
	}
	path := filepath.Join("./storage", folder, time.Now().Format(time.DateOnly))
	_ = os.MkdirAll(path, 0755)
	prefix := strconv.FormatInt(time.Now().UnixMilli(), 10) + "-"
	fileName := prefix + strings.ReplaceAll(req.Asset.Filename, " ", "-")
	pathWithFileName := filepath.Join(path, fileName)
	req.Path = pathWithFileName
	req.Size = uint64(req.Asset.Size)

	err = c.SaveFile(&req.Asset, pathWithFileName)
	if err != nil {
		return err
	}

	// TODO: store different sizes of images
	if !req.AlsoOptimize {
		resultSize, err := utils.StoreImageOptimizedVersions(path, fileName)
		if err != nil {
			return err
		}
		business.AssetsSize += uint64(resultSize)
	}

	err = _i.BusinessRepo.Update(business.ID, business)
	if err != nil {

		return err
	}

	return _i.Repo.Create(req.ToDomain())
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

	err = _i.Repo.Delete(id)
	if err != nil {
		return err
	}

	err = utils.DeleteFile(asset.Path)
	if err != nil {
		return err
	}

	business.AssetsSize -= asset.Size
	err = _i.BusinessRepo.Update(business.ID, business)
	if err != nil {
		return err
	}

	return nil
}
