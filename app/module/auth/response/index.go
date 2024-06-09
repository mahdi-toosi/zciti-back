package response

import (
	"github.com/golang-jwt/jwt/v4"
	"go-fiber-starter/app/module/user/response"
)

type Auth struct {
	Token     string
	User      response.User
	ExpiresAt *jwt.NumericDate
}
