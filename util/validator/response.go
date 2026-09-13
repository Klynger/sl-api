package validator

import (
	"fmt"

	"github.com/go-playground/validator/v10"

	errorsCommon "sl-api/api/routes/common/errors"
)

// ToErrResponse converts validator errors into error items with per-tag codes.
// It returns nil when err is not a validator.ValidationErrors, which signals
// an unexpected validation failure to the caller.
func ToErrResponse(err error) []errorsCommon.ErrorItem {
	fieldErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	items := make([]errorsCommon.ErrorItem, len(fieldErrors))

	for i, fe := range fieldErrors {
		code, detail := mapValidationTag(fe)
		items[i] = errorsCommon.NewErrorWithMeta(code, detail, map[string]any{
			"field": fe.Field(),
		})
	}

	return items
}

func mapValidationTag(fe validator.FieldError) (code, detail string) {
	switch fe.Tag() {
	case "required":
		return errorsCommon.CodeRequired, fmt.Sprintf("%s is a required field", fe.Field())
	case "max":
		return errorsCommon.CodeMaxLength, fmt.Sprintf("%s must be a maximum of %s in length", fe.Field(), fe.Param())
	case "url":
		return errorsCommon.CodeInvalidURL, fmt.Sprintf("%s must be a valid URL", fe.Field())
	case "alphaspace":
		return errorsCommon.CodeAlphaSpace, fmt.Sprintf("%s can only contain alphabetic and space characters", fe.Field())
	case "datetime":
		if fe.Param() == "2006-01-02" {
			return errorsCommon.CodeInvalidDate, fmt.Sprintf("%s must be a valid date", fe.Field())
		}
		return errorsCommon.CodeInvalidDate, fmt.Sprintf("%s must follow %s format", fe.Field(), fe.Param())
	default:
		return errorsCommon.CodeValidationFailed, fmt.Sprintf("something wrong on %s; %s", fe.Field(), fe.Tag())
	}
}
