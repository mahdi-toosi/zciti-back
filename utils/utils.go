package utils

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go-fiber-starter/app/database/schema"
	"math/rand"
	"regexp"
	"strconv"
	"time"
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

func GetIntInParams(c *fiber.Ctx, key string) (uint64, error) {
	return strconv.ParseUint(c.Params(key), 10, 64)
}

func GetIntInQueries(c *fiber.Ctx, key string) (uint64, error) {
	return strconv.ParseUint(c.Query(key), 10, 64)
}

func GetAuthenticatedUser(c *fiber.Ctx) (schema.User, error) {
	user, ok := c.Locals("user").(schema.User)
	if ok {
		return user, nil
	}
	return schema.User{}, errors.New("user doesn't exist")
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

func Random(min int, max int) uint64 {
	return uint64(rand.Intn(max-min) + min) //nolint:gosec
}

func RandomDateTime() time.Time {
	year := rand.Intn(2023) //nolint:gosec
	month := time.Month(rand.Intn(12) + 1)
	day := rand.Intn(31) + 1 //nolint:gosec
	hour := rand.Intn(24)    //nolint:gosec
	min := rand.Intn(60)     //nolint:gosec
	sec := rand.Intn(60)     //nolint:gosec

	return time.Date(year, month, day, hour, min, sec, 0, time.UTC)
}
