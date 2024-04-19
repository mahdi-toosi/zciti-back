package response

import (
	"mime/multipart"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/locales/fa_IR"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	fat "github.com/go-playground/validator/v10/translations/fa"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var (
	validate *validator.Validate
	uni      *ut.UniversalTranslator
	trans    ut.Translator
)

func init() {
	validate = validator.New()

	err := validate.RegisterValidation("validate_file", validateFile)
	if err != nil {
		return
	}

	uni = ut.New(fa_IR.New())
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
		if err := ParseBody(c, body); err != nil {
			return err
		}

		return ValidateStruct(v.Elem().Interface())
	case reflect.Struct:
		if err := ParseBody(c, &body); err != nil {
			return err
		}

		return ValidateStruct(v)
	default:
		return nil
	}
}

func validateFile(fl validator.FieldLevel) bool {
	file, ok := fl.Field().Interface().(multipart.FileHeader)
	if !ok {
		return false
	}

	// Get the extension-size pairs from the validation tag
	pairs := strings.Split(fl.Param(), " ")

	// Check if the file extension is allowed and the file size is within the specified limit
	fileExtension := filepath.Ext(file.Filename)

	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		ext := strings.TrimSpace(parts[0])
		sizeStr := strings.TrimSpace(parts[1])

		if ext != fileExtension {
			continue
		}

		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			// Return false if the size parameter is not a valid integer
			return false
		}
		maxSizeInBytes := int64(size * 1024 * 1024) // Assuming the size parameter is in MB

		// Check if the file size is less than or equal to the maxSize
		return file.Size <= maxSizeInBytes
	}

	// If the file extension is not found in the allowed list, return false
	return false
}
