package handler

import (
	"shotenedurl/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/shorten", h.createUrl)
	router.GET("/:short_url", h.redirectUrl)
	router.GET("/stats/:short_url", h.statsUrl)
	router.DELETE("/:id", h.deleteUrl)

	return router
}
