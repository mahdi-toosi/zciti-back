package controller

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/transaction/request"
	"go-fiber-starter/app/module/transaction/service"
	walletService "go-fiber-starter/app/module/wallet/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	//Show(c *fiber.Ctx) error
	//Store(c *fiber.Ctx) error
	//Update(c *fiber.Ctx) error
	//Delete(c *fiber.Ctx) error
}

func RestController(s service.IService, w walletService.IService) IRestController {
	return &controller{s, w}
}

type controller struct {
	service       service.IService
	walletService walletService.IService
}

// Index all Transactions
// @Summary      Get all transactions
// @Tags         Transactions
// @Security     Bearer
// @Param        WalletID path int true "Wallet ID"
// @Router       /wallets/:walletID/transactions [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	walletID, err := utils.GetIntInParams(c, "walletID")
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Transactions
	req.WalletID = walletID
	req.Pagination = paginate

	// get Wallet and check the owner
	wallet, err := _i.walletService.Show(&walletID, nil, nil)
	if err != nil {
		return err
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	if wallet.UserID != nil {
		if *wallet.UserID != user.ID {
			return fiber.ErrForbidden
		}
	}

	if wallet.BusinessID != nil {
		if !user.IsBusinessOwner(*wallet.BusinessID) {
			return fiber.ErrForbidden
		}
	}

	transactions, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: transactions,
		Meta: paging,
	})
}
