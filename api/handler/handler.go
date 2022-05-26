package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/eciccone/rh/api/service"
	"github.com/gin-gonic/gin"
)

func Handler(h func(c *gin.Context) error) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := h(c)
		if err == nil {
			return
		}

		log.Println(err)

		if errors.Is(err, service.ErrRecipeData) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if errors.Is(err, service.ErrIngredientData) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if errors.Is(err, service.ErrNoRecipe) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"msg": err.Error(),
			})
			return
		}

		if errors.Is(err, service.ErrRecipeForbidden) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"msg": err.Error(),
			})
			return
		}

		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
	}
}
