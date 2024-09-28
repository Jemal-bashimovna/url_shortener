package handler

import (
	services "shotenedurl/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *services.Service
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.POST("/shorten", h.createUrl)
	router.GET("", h.GetAll)
	router.GET("/:short_url", h.redirectUrl)
	router.GET("/stats/:short_url", h.statsUrl)
	router.DELETE("/:id", h.deleteUrl)

	return router
}
