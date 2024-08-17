package response

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type Response struct {
	Status int    `json:"status"`
	Error  string `json:"error,omitempty"`
}

func OK() Response {
	return Response{Status: http.StatusOK}
}

func Error(status int, msg string) Response {
	return Response{
		Status: status,
		Error:  msg,
	}
}

func ValidationErrors(errs validator.ValidationErrors) string {
	var msg []string
	for _, err := range errs {
		msg = append(msg, fmt.Sprintf("field %s not valid", err.Field()))
	}
	return fmt.Sprintf("Validation errors: %s", strings.Join(msg, ", "))
}
