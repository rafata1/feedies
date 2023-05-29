package feed_handler

import (
	"github.com/gin-gonic/gin"
	"github.com/rafata1/feedies/config"
	"github.com/rafata1/feedies/handler/common"
	"github.com/rafata1/feedies/service/feed_service"
)

type IHandler interface {
	GetNews(c *gin.Context)
}

type handler struct {
	service feed_service.IService
}

func (h *handler) GetNews(c *gin.Context) {
	output, err := h.service.GetNews(c)
	if err != nil {
		common.WriteError(c, err)
		return
	}
	common.WriteData(c, output)
}

func Init(conf config.Config) IHandler {
	db := conf.MySQL.MustConnect()
	service := feed_service.Init(db)
	return &handler{
		service: service,
	}
}
