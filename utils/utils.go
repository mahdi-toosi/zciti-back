package utils

import (
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/user/response"
)

func IsEnabled(key bool) func(c *fiber.Ctx) bool {
	if key {
		return nil
	}

	return func(c *fiber.Ctx) bool { return true }
}

func InlineCondition(condition bool, ifResult any, elseResult any) any {
	if condition {
		return ifResult
	}
	return elseResult
}

func GiveAuthUser(c *fiber.Ctx) response.User {
	return c.Locals("user").(response.User)
}
