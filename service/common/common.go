package common

import (
	"github.com/rafata1/feedies/pkg/custom_error"
	"net/http"
)

func WrapErrQueryDB(err error) error {
	return custom_error.Error{
		HTTPCode: http.StatusInternalServerError,
		Code:     501,
		Message:  err.Error(),
	}
}

func WrapErrSaveDB(err error) error {
	return custom_error.Error{
		HTTPCode: http.StatusInternalServerError,
		Code:     502,
		Message:  err.Error(),
	}
}
