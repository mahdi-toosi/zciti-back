package controller

import (
	"go-fiber-starter/app/module/wallet/request"
	"go-fiber-starter/app/module/wallet/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"

	"github.com/gofiber/fiber/v2"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
}

func RestController(s service.IService) IRestController {
	return &controller{s}
}

type controller struct {
	service service.IService
}

// Index all Wallets
// @Summary      Get all wallets
// @Tags         Wallets
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /wallets [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Wallets
	req.Pagination = paginate

	wallets, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: wallets,
		Meta: paging,
	})
}

// Show one Wallet
// @Summary      Get one wallet
// @Tags         Wallets
// @Security     Bearer
// @Router       /wallets/wallet [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	queryUserID, userErr := utils.GetUintInQueries(c, "UserID")
	queryBusinessID, businessErr := utils.GetUintInQueries(c, "BusinessID")
	if userErr != nil && businessErr != nil {
		return &fiber.Error{
			Code:    fiber.StatusBadRequest,
			Message: "id ارسال نشده است",
		}
	}
	var userID *uint64
	var businessID *uint64

	if queryUserID > 0 {
		userID = &queryUserID

		if queryUserID != user.ID {
			return fiber.ErrForbidden
		}
	}

	if queryBusinessID > 0 {
		businessID = &queryBusinessID

		if !user.IsBusinessOwner(queryBusinessID) && !user.IsBusinessObserver(queryBusinessID) {
			return fiber.ErrForbidden
		}
	}

	wallet, err := _i.service.GetOrCreateWallet(userID, businessID, nil)
	if err != nil {
		return err
	}

	return c.JSON(wallet)
}
