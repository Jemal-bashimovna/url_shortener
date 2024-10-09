package handler

import (
	"net/http"
	"urlshortener/models"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h *Handler) createUrl(ctx *gin.Context) {

	var input models.InputURL

	if err := ctx.BindJSON(&input); err != nil {
		NewErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateURL(input.OriginalURL); err != nil {
		NewErrorMessage(ctx, http.StatusBadRequest, err.Error())
		return
	}

	generatedShortURL, err := h.service.CreateURL(input.OriginalURL)
	if err != nil {
		NewErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	shortenURL := viper.GetString("domain") + generatedShortURL

	ctx.JSON(http.StatusOK, models.CreateResponse{
		ShortURL: shortenURL,
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

	http.Redirect(ctx.Writer, ctx.Request, originalURL, http.StatusMovedPermanently)
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

func (h *Handler) deleteUrl(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		NewErrorMessage(ctx, http.StatusBadRequest, "invalid id param")
		return
	}

	if err := h.service.DeleteURL(id); err != nil {
		NewErrorMessage(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.DeleteStatus{
		Status: "Success",
	})
}
