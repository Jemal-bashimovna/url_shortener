package handler

import (
	"net/http"
	urlshortener "shotenedurl"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h *Handler) createUrl(ctx *gin.Context) {

	var input urlshortener.InputURL
	if err := ctx.BindJSON(&input); err != nil {
		NewErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := input.ValidateURL(input.OriginalURL); err != nil {
		NewErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	str, err := h.service.CreateURL(input.OriginalURL)
	if err != nil {
		NewErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	shortenURL := viper.GetString("domain") + str
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"shortened_url": shortenURL,
	})

}

func (h *Handler) GetAll(ctx *gin.Context) {
	list, err := h.service.GetAll()
	if err != nil {
		NewErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})

}
func (h *Handler) redirectUrl(ctx *gin.Context) {
	short_url := ctx.Param("short_url")

	originalURL, err := h.service.RedirectURL(short_url)

	if err != nil {
		NewErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if len(originalURL) == 0 {
		NewErrorMessage(ctx, http.StatusBadRequest, "short URL not found")
		return
	}

	ctx.JSON(http.StatusMovedPermanently, map[string]interface{}{
		"original_URL": originalURL,
	})
}

func (h *Handler) statsUrl(ctx *gin.Context) {

	shortURL := ctx.Param("short_url")

	stats, err := h.service.GetStatsURL(shortURL)

	if err != nil {
		NewErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, map[string]interface{}{
		"stats": stats,
	})
}

func (h *Handler) deleteUrl(ctx *gin.Context) {}
