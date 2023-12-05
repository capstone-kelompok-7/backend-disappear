package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	_ = validate.RegisterValidation("noSpace", noSpace)
}

func noSpace(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return !strings.Contains(password, " ")
}

func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		var customErrors []string

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' wajib diisi", err.Field()))
			case "min":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' minimal harus memiliki panjang %s karakter", err.Field(), err.Param()))
			case "email":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' harus berupa alamat email yang valid", err.Field()))
			case "noSpace":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' tidak boleh mengandung spasi", err.Field()))
			default:
				customErrors = append(customErrors, fmt.Sprintf("Validasi field '%s' gagal dengan tag '%s'", err.Field(), err.Tag()))
			}
		}

		return errors.New(strings.Join(customErrors, "; "))
	}
	return nil
}
