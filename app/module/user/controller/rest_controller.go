package controller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	bService "go-fiber-starter/app/module/business/service"
	orequest "go-fiber-starter/app/module/order/request"
	oService "go-fiber-starter/app/module/order/service"
	"go-fiber-starter/app/module/user/request"
	"go-fiber-starter/app/module/user/service"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/paginator"
	"go-fiber-starter/utils/response"
	"golang.org/x/exp/slices"
	"strconv"
)

type IRestController interface {
	Index(c *fiber.Ctx) error
	Show(c *fiber.Ctx) error
	Store(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	UpdateAccount(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error

	BusinessUsers(c *fiber.Ctx) error
	InsertUser(c *fiber.Ctx) error
	DeleteUser(c *fiber.Ctx) error
	BusinessUsersAddRole(c *fiber.Ctx) error

	Orders(c *fiber.Ctx) error
	OrderStore(c *fiber.Ctx) error
	OrderStatus(c *fiber.Ctx) error
}

func RestController(s service.IService, b bService.IService, o oService.IService, config *config.Config) IRestController {
	return &controller{s, b, o, config}
}

type controller struct {
	service  service.IService
	bService bService.IService
	oService oService.IService
	Config   *config.Config
}

// Index
// @Summary      Get all users
// @Tags         Users
// @Security     Bearer
// @Router       /users [get]
func (_i *controller) Index(c *fiber.Ctx) error {
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.Users
	req.Pagination = paginate
	req.Keyword = c.Query("Keyword")

	users, paging, err := _i.service.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: users,
		Meta: paging,
	})
}

// Show
// @Summary      Get one user
// @Tags         Users
// @Security     Bearer
// @Param        id path int true "User ID"
// @Router       /users/:id [get]
func (_i *controller) Show(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	user, err := _i.service.Show(id)
	if err != nil {
		return err
	}

	return c.JSON(user)
}

// Store
// @Summary      Create user
// @Tags         Users
// @Param 		 user body request.User true "User details"
// @Router       /users [post]
func (_i *controller) Store(c *fiber.Ctx) error {
	req := new(request.User)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	err := _i.service.Store(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Update
// @Summary      update user
// @Security     Bearer
// @Tags         Users
// @Param 		 user body request.User true "User details"
// @Param        id path int true "User ID"
// @Router       /users/:id [put]
func (_i *controller) Update(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	req := new(request.User)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	err = _i.service.Update(id, *req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// UpdateAccount
// @Summary      update user
// @Security     Bearer
// @Tags         Users
// @Param 		 user body request.UpdateUserAccount true "User details"
// @Param        id path int true "User ID"
// @Router       /users/user/account [put]
func (_i *controller) UpdateAccount(c *fiber.Ctx) error {
	req := new(request.UpdateUserAccount)
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}
	if err := utils.ValidateMobileNumber(strconv.FormatUint(req.Mobile, 10)); err != nil {
		return err
	}

	err := _i.service.UpdateAccount(*req)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// Delete
// @Summary      delete user
// @Tags         Users
// @Security     Bearer
// @Param        id path int true "User ID"
// @Router       /users/:id [delete]
func (_i *controller) Delete(c *fiber.Ctx) error {
	id, err := utils.GetIntInParams(c, "id")
	if err != nil {
		return err
	}

	err = _i.service.Destroy(id)
	if err != nil {
		return err
	}

	return c.JSON("success")
}

// BusinessUsers
// @Summary      Get one business users
// @Tags         Users
// @Security     Bearer
// @Param        id path int true "Business ID"
// @Router       /businesses/:businessID/users [get]
func (_i *controller) BusinessUsers(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	business, err := _i.bService.Show(businessID, schema.URUser)
	if err != nil {
		return fiber.ErrNotFound
	}

	if user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req request.BusinessUsers
	req.Pagination = paginate
	req.BusinessID = businessID

	users, paging, err := _i.service.Users(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: users,
		Meta: paging,
	})
}

// InsertUser
// @Summary      Insert one business user
// @Tags         Users
// @Security     Bearer
// @Param        businessId path int true "Business ID" ,userId path int true "User ID"
// @Router       /businesses/:businessID/users/:userID [post]
func (_i *controller) InsertUser(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := utils.GetIntInParams(c, "userID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	business, err := _i.bService.Show(businessID, schema.URUser)
	if err != nil {
		return fiber.ErrNotFound
	}

	if user.ID != userID && user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	err = _i.service.InsertUser(businessID, userID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"success"},
	})
}

// DeleteUser
// @Summary      Delete one business user
// @Tags         Users
// @Security     Bearer
// @Param        businessId path int true "Business ID" ,userId path int true "User ID"
// @Router       /businesses/:businessID/users/:userID [delete]
func (_i *controller) DeleteUser(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	userID, err := utils.GetIntInParams(c, "userID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	business, err := _i.bService.Show(businessID, schema.URUser)
	if err != nil {
		return fiber.ErrNotFound
	}

	if user.ID != userID && user.ID != business.OwnerID && !user.IsAdmin() {
		return fiber.ErrForbidden
	}

	err = _i.service.DeleteUser(businessID, userID)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"success"},
	})
}

// BusinessUsersAddRole
// @Summary      BusinessUsersAddRole add role to business user
// @Tags         Users
// @Security     Bearer
// @Param        businessId path int true "Business ID"
// @Router       /businesses/:businessID/users-add-role [post]
func (_i *controller) BusinessUsersAddRole(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return fiber.ErrBadRequest
	}

	req := new(request.BusinessUsersStoreRole)
	req.BusinessID = businessID
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	validRoles := []schema.UserRole{
		schema.URUser,
		schema.URAdmin,
		schema.URBusinessOwner,
		schema.URReservationViewer,
	}

	for _, role := range req.Roles {
		if !slices.Contains(validRoles, role) {
			return &fiber.Error{Code: fiber.StatusBadRequest, Message: "درخواست شما معتبر نمی باشد."}
		}
	}

	if err := _i.service.BusinessUsersAddRole(*req); err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Messages: response.Messages{"success"},
	})
}

// Orders
// @Summary      Get all orders
// @Tags         Users
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /user/orders [get]
func (_i *controller) Orders(c *fiber.Ctx) error {
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}
	paginate, err := paginator.Paginate(c)
	if err != nil {
		return err
	}

	var req orequest.Orders
	req.UserID = user.ID
	req.Pagination = paginate

	orders, paging, err := _i.oService.Index(req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: orders,
		Meta: paging,
	})
}

// OrderStatus
// @Summary      Get all orders
// @Tags         Users
// @Security     Bearer
// @Param        businessID path int true "Business ID"
// @Router       /user/orders/status [get]
func (_i *controller) OrderStatus(c *fiber.Ctx) error {
	userID, err := utils.GetIntInQueries(c, "UserID")
	if err != nil {
		return err
	}
	orderID, err := utils.GetIntInQueries(c, "OrderID")
	if err != nil {
		return err
	}
	authority := c.Query("Authority")
	status := c.Query("Status")

	if status == "OK" {
		status, err = _i.oService.Status(userID, orderID, authority)
		if err != nil {
			return err
		}
	}

	url := fmt.Sprintf(
		"%s/card/payment/result?Status=%s",
		_i.Config.App.FrontendDomain,
		status,
	)

	return c.Redirect(url)
}

// OrderStore
// @Summary      Create order
// @Tags         Orders
// @Param 		 order body request.Order true "Order details"
// @Param        businessID path int true "Business ID"
// @Router       /business/:businessID/orders [post]
// @Router       /user/orders [post]
func (_i *controller) OrderStore(c *fiber.Ctx) error {
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}

	req := new(orequest.Order)
	req.User = user
	if err := response.ParseAndValidate(c, req); err != nil {
		return err
	}

	orderID, paymentURL, err := _i.oService.Store(*req)
	if err != nil {
		return err
	}

	return response.Resp(c, response.Response{
		Data: map[string]any{"paymentUrl": paymentURL, "orderID": orderID},
	})
}
