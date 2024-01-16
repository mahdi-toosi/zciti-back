package service

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go-fiber-starter/app/database/schema"
	"go-fiber-starter/app/module/messageRoom/repository"
	"go-fiber-starter/app/module/messageRoom/request"
	"go-fiber-starter/app/module/messageRoom/response"
	"go-fiber-starter/utils/config"
	"go-fiber-starter/utils/paginator"
	"time"
)

type IService interface {
	Index(req request.MessageRooms) (messageRooms []*response.MessageRoom, paging paginator.Pagination, err error)
	Destroy(id uint64) error
	Update(id uint64, req request.MessageRoom) (err error)
	GenerateToken(messageRoom *schema.MessageRoom, business *schema.Business) (token string, err error)
	IsTokenValid(token string) (tokenData *response.MessageRoomToken, err error)
	// Show(id uint64) (messageRoom *response.MessageRoom, err error)
	// Store(req request.MessageRoom) (err error)
}

func Service(Repo repository.IRepository, config *config.Config) IService {
	return &service{
		Repo,
		config,
	}
}

type service struct {
	Repo   repository.IRepository
	config *config.Config
}

func (_i *service) Index(req request.MessageRooms) (messageRooms []*response.MessageRoom, paging paginator.Pagination, err error) {
	results, paging, err := _i.Repo.GetAll(req)
	if err != nil {
		return
	}

	for _, result := range results {
		messageRooms = append(messageRooms, response.FromDomain(result, nil))
	}

	return
}

// func (_i *service) Show(id uint64) (messageRoom *response.MessageRoom, err error) {
//	result, err := _i.Repo.GetOne(id)
//	if err != nil {
//		return nil, err
//	}
//
//	return response.FromDomain(result), nil
//}

// func (_i *service) Store(req request.MessageRoom) (err error) {
//	return _i.Repo.Create(req.ToDomain())
// }

func (_i *service) Update(id uint64, req request.MessageRoom) (err error) {
	return _i.Repo.Update(id, req.ToDomain())
}

func (_i *service) Destroy(id uint64) error {
	return _i.Repo.Delete(id)
}

func (_i *service) GenerateToken(messageRoom *schema.MessageRoom, business *schema.Business) (token string, err error) {
	ExpiresAt := jwt.NewNumericDate(time.Now().Add(14400 /*60 * 60 * 4 */ * time.Second))

	unSignedToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"exp": ExpiresAt,
			"messageRoom": response.MessageRoomToken{
				ID:              messageRoom.ID,
				UserID:          messageRoom.UserID,
				Status:          messageRoom.Status,
				BusinessID:      messageRoom.BusinessID,
				MembersAsString: fmt.Sprintf("%d,%d", business.OwnerID, messageRoom.UserID),
			},
		},
	)
	token, err = unSignedToken.SignedString([]byte(_i.config.Middleware.Jwt.Secret))

	if err != nil {
		return "", err
	}

	return token, nil
}

func (_i *service) IsTokenValid(token string) (tokenData *response.MessageRoomToken, err error) {
	_token, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if token.Method.Alg() != jwt.SigningMethodHS256.Name {
			return nil, err
		}
		// Return the key for validation
		return []byte(_i.config.Middleware.Jwt.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := _token.Claims.(jwt.MapClaims)
	if ok && _token.Valid {
		data := claims["messageRoom"].(map[string]interface{})

		tokenData = &response.MessageRoomToken{
			Status:     data["Status"].(string),
			ID:         uint64(data["ID"].(float64)),
			UserID:     uint64(data["UserID"].(float64)),
			BusinessID: uint64(data["BusinessID"].(float64)),
			//Members:    strings.Split(data["MembersAsString"].(string), ","),
		}
		return
	}
	return nil, err
}
