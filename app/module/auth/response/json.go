package response

import (
	"go-fiber-starter/app/module/user/response"
)

type LoginResponse struct {
	User  response.User `json:"user"`
	Token string        `json:"token"`
}
