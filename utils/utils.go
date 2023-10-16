package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/module/user/response"
	"regexp"
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

func GiveAuthUser(c *fiber.Ctx) (response.User, error) {
	user, ok := c.Locals("user").(response.User)
	if ok {
		return user, nil
	}
	return response.User{}, errors.New("user doesn't exist")
}

func ValidateMobileNumber(number string) error {
	// Define the regular expression pattern for a mobile number
	pattern := `9(1[0-9]|3[1-9]|2[1-9])-?[0-9]{3}-?[0-9]{4}`

	// Compile the regex pattern
	regex := regexp.MustCompile(pattern)

	// Match the number against the regex pattern
	valid := regex.MatchString(number)
	if !valid {
		return errors.New("تلفن همراه معتبر نمی باشد")
	}

	return nil
}
