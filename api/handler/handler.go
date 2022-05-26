package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/eciccone/rh/api/service"
	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidJSON = errors.New("invalid json data")
)

func Handler(h func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c)
		if err == nil {
			return
		}

		// handle 400
		if errors.Is(err, ErrInvalidJSON) ||
			errors.Is(err, service.ErrProfileExists) ||
			errors.Is(err, service.ErrProfileData) ||
			errors.Is(err, service.ErrRecipeData) ||
			errors.Is(err, service.ErrIngredientData) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// handle 403
		if errors.Is(err, service.ErrRecipeForbidden) || errors.Is(err, service.ErrUsernameForbidden) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"msg": err.Error(),
			})
			return
		}

		// handle 404
		if errors.Is(err, service.ErrNoRecipe) || errors.Is(err, service.ErrNoProfile) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		log.Printf("internal server error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
	}
}
