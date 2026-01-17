package validator

import (
	"github.com/go-playground/validator/v10"
	"reflect"
	"regexp"
	"strings"
)

const alphaSpaceRegexString = "^[a-zA-Z ]+$"

func New() *validator.Validate {
	validate := validator.New()
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	validate.RegisterValidation("alphaspace", isAlphaSpace)

	return validate
}

func isAlphaSpace(fl validator.FieldLevel) bool {
	reg := regexp.MustCompile(alphaSpaceRegexString)
	return reg.MatchString(fl.Field().String())
}
