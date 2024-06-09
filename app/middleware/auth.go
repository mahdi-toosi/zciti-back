package middleware

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/golang-jwt/jwt/v4"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/utils/config"
	"time"
)

type JWTCustomClaim struct {
	User schema.User
	jwt.RegisteredClaims
}

func Protected(cfg *config.Config) fiber.Handler {
	return jwtware.New(jwtware.Config{
		ErrorHandler:   jwtError,
		SuccessHandler: jwtSuccess,
		Claims:         &JWTCustomClaim{},
		SigningKey:     []byte(cfg.Middleware.Jwt.Secret),
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
	jwtCustomClaim := c.Locals("user").(*jwt.Token).Claims.(*JWTCustomClaim)

	c.Locals("user", jwtCustomClaim.User)
	return c.Next()
}

func GenerateTokenAccess(user schema.User, jwtConfig config.Jwt) (token string, expiresAt *jwt.NumericDate, err error) {
	conf := config.NewConfig()
	expiresAt = jwt.NewNumericDate(time.Now().Add(jwtConfig.Expiration * time.Second))

	jwtCustomClaim := JWTCustomClaim{
		User: schema.User{
			ID:              user.ID,
			Mobile:          user.Mobile,
			LastName:        user.LastName,
			FirstName:       user.FirstName,
			Permissions:     user.Permissions,
			MobileConfirmed: user.MobileConfirmed,
		},

		RegisteredClaims: jwt.RegisteredClaims{ExpiresAt: expiresAt},
	}

	unSignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtCustomClaim)
	token, err = unSignedToken.SignedString([]byte(conf.Middleware.Jwt.Secret))
	if err != nil {
		return "", nil, err
	}

	return token, expiresAt, nil
}
