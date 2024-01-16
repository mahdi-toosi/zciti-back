package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/user/response"
	"go-fiber-starter/utils"
	"go-fiber-starter/utils/config"
	"strconv"
	"strings"
	"time"
)

func Protected(cfg *config.Config) fiber.Handler {
	if cfg.Middleware.Jwt.Secret == "" {
		panic("JWT secret is not set")
	}

	return jwtware.New(jwtware.Config{
		SigningKey:     []byte(cfg.Middleware.Jwt.Secret),
		ErrorHandler:   jwtError,
		SuccessHandler: jwtSuccess,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}

	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
}

func jwtSuccess(c *fiber.Ctx) error {
	userJSON := c.Locals("user").(*jwt.Token).Claims.(jwt.MapClaims)["user"].(map[string]interface{})

	ID, _ := strconv.ParseUint(userJSON["ID"].(string), 10, 64)
	Mobile, _ := strconv.ParseUint(userJSON["Mobile"].(string), 10, 64)

	c.Locals("user", schema.User{
		ID:              ID,
		FirstName:       userJSON["FirstName"].(string),
		LastName:        userJSON["LastName"].(string),
		Mobile:          Mobile,
		MobileConfirmed: utils.InlineCondition(userJSON["MobileConfirmed"] == "true", true, false).(bool),
		Roles:           strings.Split(userJSON["Roles"].(string), ","),
	})

	return c.Next()
}

func GenerateTokenAccess(user response.User, jwtConfig config.Jwt) (token string, err error) {
	conf := config.NewConfig()

	ExpiresAt := jwt.NewNumericDate(time.Now().Add(jwtConfig.Expiration * time.Second))

	unSignedToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": ExpiresAt,
			"user": map[string]string{
				"ID":              strconv.FormatUint(user.ID, 10),
				"FirstName":       user.FirstName,
				"LastName":        user.LastName,
				"Mobile":          strconv.FormatUint(user.Mobile, 10),
				"MobileConfirmed": fmt.Sprintf("%t", user.MobileConfirmed),
				"Roles":           strings.Join(user.Roles, ","),
			},
		},
	)

	token, err = unSignedToken.SignedString([]byte(conf.Middleware.Jwt.Secret))

	if err != nil {
		return "", err
	}

	return token, nil
}
