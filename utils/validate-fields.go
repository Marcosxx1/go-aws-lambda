package utils

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// NewError creates a new HTTP error response and sends it as JSON through the given Gin context.
// It takes three parameters: the Gin context, the HTTP status code, and the error.
func NewError(context *gin.Context, status int, err error) {
	er := HTTPError{
		Code:    status,
		Message: err.Error(),
	}
	context.JSON(status, er)
}

// HTTPError represents an HTTP error response.
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// ValidateStruct validates the fields of a given struct using the validator package.
// It takes an interface representing the struct to be validated.
// It returns an error if validation fails, otherwise, it returns nil.
//
// Example:
//
//	type ExampleStruct struct {
//	    Name  string `validate:"required"`
//	    Age   int    `validate:"required,min=18"`
//	}
//
//	example := ExampleStruct{Name: "", Age: 15}
//	err := ValidateStruct(example)
//	if err != nil {
//	    fmt.Println("Validation errors:", err)
//	}
func ValidateStruct(obj interface{}) error {
	validate := validator.New()
	err := validate.Struct(obj)

	if err == nil {
		return nil
	}

	var errorMessages []string

	switch e := err.(type) {
	case validator.ValidationErrors:
		for _, validationError := range e {
			field := strings.ToLower(validationError.StructField())

			switch validationError.Tag() {
			case "required":
				errorMessages = append(errorMessages, field+" is required")
			case "min":
				errorMessages = append(errorMessages, field+" is required with min "+validationError.Param())
			case "oneof":
				errorMessages = append(errorMessages, "Invalid "+field+" format")
			case "gtfield":
				errorMessages = append(errorMessages, field+" must be greater than "+validationError.Param())
			}
		}
	case *validator.InvalidValidationError:
		errorMessages = append(errorMessages, "Invalid validation error")
	default:
		return e
	}

	return errors.New(strings.Join(errorMessages, "; "))
}
