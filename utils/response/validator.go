package response

import (
	"reflect"

	"github.com/rs/zerolog/log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	fat "github.com/go-playground/validator/v10/translations/fa"
	"github.com/gofiber/fiber/v2"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	validate = validator.New()

	uni = ut.New(en.New())
	trans, _ = uni.GetTranslator("fa")

	if err := fat.RegisterDefaultTranslations(validate, trans); err != nil && !fiber.IsChild() {
		log.Panic().Err(err).Msg("")
	}
}

func ValidateStruct(input any) error {
	return validate.Struct(input)
}

func ParseBody(c *fiber.Ctx, body any) error {
	if err := c.BodyParser(body); err != nil {
		return err
	}

	return nil
}

func ParseAndValidate(c *fiber.Ctx, body any) error {
	v := reflect.ValueOf(body)

	switch v.Kind() {
	case reflect.Ptr:
		ParseBody(c, body)

		return ValidateStruct(v.Elem().Interface())
	case reflect.Struct:
		ParseBody(c, &body)

		return ValidateStruct(v)
	default:
		return nil
	}
}
