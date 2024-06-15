package middleware

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
)

func AdminPermission(c *fiber.Ctx) error {
	user, err := utils.GetAuthenticatedUser(c)
	if err != nil {
		return err
	}
	if user.IsAdmin() {
		return c.Next()
	}

	return c.Status(fiber.StatusForbidden).
		JSON(fiber.Map{"status": "error", "message": "you don't have permission", "data": nil})
}

func ForUser(c *fiber.Ctx) error {
	c.Locals("forUser", true)
	return c.Next()
}

func BusinessPermission(domain Domain, permission Permission) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetAuthenticatedUser(c)
		if err != nil {
			return err
		}

		// TODO should i remove it ?! 👇🏻
		if user.IsAdmin() {
			return c.Next()
		}

		businessID, err := utils.GetIntInParams(c, "businessID")
		if err != nil {
			return errors.New("مسیر وارد شده صحیح نیست")
		}

		if hasPermission(user.Permissions, businessID, domain, permission) {
			return c.Next()
		}
		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"status": "error", "message": "شما دسترسی لازم را ندارید", "data": nil})
	}
}

func hasPermission(userPermissions schema.UserPermissions, businessID uint64, domain Domain, permission Permission) bool {
	for _, role := range userPermissions[businessID] {
		r, ok1 := Permissions[role]
		if !ok1 {
			continue
		}

		d, ok2 := r[domain]
		if !ok2 {
			continue
		}

		p, ok3 := d[permission]
		if !ok3 {
			continue
		}

		if p {
			return true
		}
	}

	return false
}

// define Permissions

var Permissions = map[schema.UserRole]map[Domain]map[Permission]bool{
	schema.URAdmin: {
		DUser:                 {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DFile:                 {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DPost:                 {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DOrder:                {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DCoupon:               {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DProduct:              {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DComment:              {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DMessage:              {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DTaxonomy:             {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DBusiness:             {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DTransaction:          {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DMessageRoom:          {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DReservation:          {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DNotification:         {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DNotificationTemplate: {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
	},
	schema.URBusinessOwner: {
		DUser:                 {PReadAll: true, PReadSingle: true, PDelete: true},
		DFile:                 {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DPost:                 {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DOrder:                {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true},
		DProduct:              {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DComment:              {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DMessage:              {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DTaxonomy:             {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DBusiness:             {PCreate: true, PReadSingle: true, PUpdate: true},
		DMessageRoom:          {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DReservation:          {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DNotification:         {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DNotificationTemplate: {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
	},
	schema.URUser: {
		DTaxonomy:     {PReadAll: true},
		DBusiness:     {PReadSingle: true},
		DPost:         {PReadSingle: true},
		DOrder:        {PReadSingle: true, PCreate: true},
		DProduct:      {PReadAll: true, PReadSingle: true},
		DReservation:  {PReadAll: true, PReadSingle: true},
		DComment:      {PCreate: true, PReadAll: true, PUpdate: true, PDelete: true},
		DNotification: {PReadAll: true, PReadSingle: true},
	},
}

// end define Permissions

// define Domains

type Domain int

const (
	DPost Domain = iota
	DUser
	DOrder
	DFile
	DCoupon
	DProduct
	DComment
	DMessage
	DTaxonomy
	DBusiness
	DReservation
	DTransaction
	DNotification
	DMessageRoom
	DNotificationTemplate
)

// end define Domains

// define Permissions

type Permission int

const (
	PCreate Permission = iota
	PUpdate
	PDelete
	PReadAll
	PReadSingle
)

// end define Permissions

type account struct {
	Title           string
	AssetsSizeLimit uint64
}

const megabyte = 1000 * 1024

var Accounts = map[schema.BusinessAccount]account{
	schema.BusinessAccountDefault: {Title: "DefaultAccount", AssetsSizeLimit: 1 * megabyte},
}
