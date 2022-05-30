package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/eciccone/rh/api/service"
	"github.com/gin-gonic/gin"
)

type RecipeHandler struct {
	recipeService service.RecipeService
}

func NewRecipeHandler(recipeService service.RecipeService) RecipeHandler {
	return RecipeHandler{recipeService}
}

// delete /recipes/:id
func (h *RecipeHandler) DeleteRecipe(c *gin.Context) error {
	username := c.GetString("username")
	recipeId, _ := strconv.Atoi(c.Param("id"))

	err := h.recipeService.RemoveRecipe(recipeId, username)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "recipe deleted",
	})

	return nil
}

// put /recipes/:id
func (h *RecipeHandler) PutRecipe(c *gin.Context) error {
	var input recipe.Recipe
	if err := c.ShouldBindJSON(&input); err != nil {
		return ErrInvalidJSON
	}

	input.Username = c.GetString("username")
	if input.Username == "" {
		return errors.New("PutRecipe failed to get username, should have been set in middleware")
	}

	recipeId, _ := strconv.Atoi(c.Param("id"))
	input.Id = recipeId

	result, err := h.recipeService.UpdateRecipe(input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "recipe updated",
		"recipe": result,
	})

	return nil
}

// put /recipes/:id/image
func (h *RecipeHandler) PutRecipeImage(c *gin.Context) error {
	recipeId, _ := strconv.Atoi(c.Param("id"))

	username := c.GetString("username")
	if username == "" {
		return errors.New("PutRecipeImage failed to get username, should have been set in middleware")
	}

	file, err := c.FormFile("image")
	if errors.Is(err, http.ErrMissingFile) {
		return ErrMissingFile
	}
	if err != nil {
		return fmt.Errorf("PutRecipeImage failed to get file: %w", err)
	}

	imagename, err := h.recipeService.UpdateRecipeImage(recipeId, username, file)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "recipe image updated",
		"recipe": imagename,
	})

	return nil
}

// post /recipes
func (h *RecipeHandler) PostRecipe(c *gin.Context) error {
	var input recipe.Recipe
	if err := c.ShouldBindJSON(&input); err != nil {
		return ErrInvalidJSON
	}

	input.Username = c.GetString("username")
	if input.Username == "" {
		return errors.New("PostRecipe failed to get username, should have been set in middleware")
	}

	result, err := h.recipeService.CreateRecipe(input)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "recipe created",
		"recipe": result,
	})

	return nil
}

// get /recipes/:id
func (h *RecipeHandler) GetRecipe(c *gin.Context) error {
	recipeId, _ := strconv.Atoi(c.Param("id"))

	recipe, err := h.recipeService.GetRecipe(recipeId)
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "recipe found",
		"recipe": recipe,
	})

	return nil
}

// get /recipes[&limit=][&offset=]
func (h *RecipeHandler) GetRecipes(c *gin.Context) error {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 64)

	username := c.GetString("username")
	if username == "" {
		return errors.New("GetRecipes failed to get username, should have been set in middleware")
	}

	recipePage, err := h.recipeService.GetRecipesForUsername(username, "", int(offset), int(limit))
	if err != nil {
		return err
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":     "recipe found",
		"recipes": recipePage.Recipes,
		"limit":   recipePage.Limit,
		"offset":  recipePage.Offset,
		"total":   recipePage.Total,
	})

	return nil
}
