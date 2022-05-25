package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/eciccone/rh/api/rherr"
	"github.com/eciccone/rh/api/service"
	"github.com/gin-gonic/gin"
)

type recipeHandler struct {
	recipeService service.RecipeService
}

func NewRecipeHandler(recipeService service.RecipeService) recipeHandler {
	return recipeHandler{recipeService}
}

// /recipes/:id
func (h *recipeHandler) GetRecipe(c *gin.Context) {
	recipeId, _ := strconv.Atoi(c.Param("id"))

	recipe, err := h.recipeService.GetRecipe(recipeId)
	if err != nil {
		if errors.Is(err, rherr.ErrBadRequest) {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "invalid id",
			})
		} else if errors.Is(err, rherr.ErrNotFound) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"msg": "recipe not found",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "recipe found",
		"recipe": recipe,
	})
}
