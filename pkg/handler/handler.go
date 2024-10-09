package handler

import (
	"urlshortener/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/shorten", h.createUrl)
	router.GET("/:short_url", h.redirectUrl)
	router.GET("/stats/:short_url", h.statsUrl)
	router.DELETE("/:id", h.deleteUrl)

	return router
}
