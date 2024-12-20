package requests

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

// This is a global validator instance
var Validate = validator.New()

// Every request struct could use this function to format the error message
func FormatError(err error, customErrorMsg map[string]string) error {

	if err == nil {
		return nil
	}

	for _, e := range (err).(validator.ValidationErrors) {
		fieldTag := fmt.Sprintf("%s.%s", e.Field(), e.Tag())
		// log.Println(fieldTag)
		if msg, exists := customErrorMsg[fieldTag]; exists {
			return errors.New(msg)
		}
	}

	// Fallback if no custom message is defined
	return errors.New("invalid request params")
}
