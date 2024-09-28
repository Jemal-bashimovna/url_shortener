package handler

import (
	"net/http"
	urlshortener "shotenedurl"

	"github.com/gin-gonic/gin"
)

func (h *Handler) createUrl(ctx *gin.Context) {

	var url urlshortener.URL
	if err := ctx.BindJSON(&url); err != nil {
		NewErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := url.ValidateURL(url.OriginalURL); err != nil {
		NewErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	shortenURL, err := h.services.CreateURL(url)
	if err != nil {
		NewErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"shortened_url": shortenURL,
	})

}

func (h *Handler) redirectUrl(ctx *gin.Context) {}

func (h *Handler) statsUrl(ctx *gin.Context) {}

func (h *Handler) deleteUrl(ctx *gin.Context) {}
