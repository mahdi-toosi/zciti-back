package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils"
)

func Permission(
	domain DomainType,
	permission PermissionType,
) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user, err := utils.GetAuthenticatedUser(c)
		if err != nil {
			return err
		}

		for _, role := range user.Roles {
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
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).
			JSON(fiber.Map{"status": "error", "message": "you don't have permission", "data": nil})
	}
}

// define Permissions

var Permissions = map[string]map[DomainType]map[PermissionType]bool{
	schema.RAdmin: {
		DUser:                 {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DPost:                 {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DBusiness:             {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DNotification:         {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DNotificationTemplate: {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DMessage:              {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
		DMessageRoom:          {PCreate: true, PReadAll: true, PReadSingle: true, PUpdate: true, PDelete: true},
	},
	schema.RUser: {},
}

// end define Permissions

// define Roles

// end define Roles

// define Domains

type DomainType int

const (
	DPost DomainType = iota
	DUser
	DNotification
	DNotificationTemplate
	DBusiness
	DMessage
	DMessageRoom
)

// end define Domains

// define Permissions

type PermissionType int

const (
	PCreate PermissionType = iota
	PReadAll
	PReadSingle
	PUpdate
	PDelete
)

// end define Permissions
