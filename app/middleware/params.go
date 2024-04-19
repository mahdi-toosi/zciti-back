package middleware

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/utils"
)

func SetBusinessID(c *fiber.Ctx) error {
	businessID, err := utils.GetIntInParams(c, "businessID")
	if err != nil {
		return err
	}

	c.Locals("businessID", businessID)
	return c.Next()
}
