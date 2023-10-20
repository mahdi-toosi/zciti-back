package response

import (
	"go-fiber-starter/app/module/user/response"
)

type Login struct {
	User  response.User
	Token string
}

type Register struct {
	User  response.User
	Token string
}
