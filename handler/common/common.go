package common

import (
	"github.com/gin-gonic/gin"
	"github.com/rafata1/feedies/pkg/custom_error"
	"net/http"
)

func WriteError(c *gin.Context, err error) {
	customErr := err.(custom_error.Error)
	c.JSON(customErr.HTTPCode, baseResponse{Code: customErr.Code, Message: customErr.Message})
}

func WriteData(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, baseResponse{Code: successCode, Message: successMessage, Data: data})
}
