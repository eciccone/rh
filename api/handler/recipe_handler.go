package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/eciccone/rh/api/repo/recipe"
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

func (h *recipeHandler) DeleteRecipe(c *gin.Context) {
	username := c.GetString("username")
	recipeId, _ := strconv.Atoi(c.Param("id"))

	err := h.recipeService.RemoveRecipe(recipeId, username)
	if err != nil {
		if errors.Is(err, rherr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "recipe does not exist",
			})
		} else if errors.Is(err, rherr.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "invalid id",
			})
		} else if errors.Is(err, rherr.ErrForbidden) {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "forbidden",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "recipe deleted",
	})
}

func (h *recipeHandler) PutRecipe(c *gin.Context) {
	var input recipe.Recipe
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid json data",
		})
		return
	}
	input.Username = c.GetString("username")
	recipeId, _ := strconv.Atoi(c.Param("id"))
	input.Id = recipeId

	result, err := h.recipeService.UpdateRecipe(input)
	if err != nil {
		if errors.Is(err, rherr.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "must provide a name for recipe and ingredients",
			})
		} else if errors.Is(err, rherr.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{
				"msg": "recipe does not exist",
			})
		} else if errors.Is(err, rherr.ErrForbidden) {
			c.JSON(http.StatusForbidden, gin.H{
				"msg": "forbidden",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "recipe updated",
		"recipe": result,
	})
}

func (h *recipeHandler) PostRecipe(c *gin.Context) {
	var input recipe.Recipe
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": "invalid json data",
		})
		return
	}
	input.Username = c.GetString("username")

	result, err := h.recipeService.CreateRecipe(input)
	if err != nil {
		if errors.Is(err, rherr.ErrBadRequest) {
			c.JSON(http.StatusBadRequest, gin.H{
				"msg": "must provide a recipe name",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"msg": "internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "recipe created",
		"recipe": result,
	})
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

func (h *recipeHandler) GetRecipes(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)
	username := c.GetString("username")

	recipePage, err := h.recipeService.GetRecipesForUsername(username, "", int(offset), int(limit))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "recipe found",
		"recipes": recipePage.Recipes,
		"limit":   recipePage.Limit,
		"offset":  recipePage.Offset,
		"total":   recipePage.Total,
	})
}
