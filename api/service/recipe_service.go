package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/eciccone/rh/api/rherr"
)

type RecipeService interface {
	CreateRecipe(recipe.Recipe) (recipe.Recipe, error)
	GetRecipe(id int) (recipe.Recipe, error)
}

type recipeService struct {
	recipeRepo recipe.RecipeRepository
}

func NewRecipeService(recipeRepo recipe.RecipeRepository) RecipeService {
	return &recipeService{recipeRepo}
}

func (s *recipeService) CreateRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	result, err := s.recipeRepo.InsertRecipe(args)
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("CreateRecipe failed to create recipe: %w", err)
	}

	return result, nil
}

func (s *recipeService) GetRecipe(id int) (recipe.Recipe, error) {
	result, err := s.recipeRepo.SelectRecipeById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recipe.Recipe{}, rherr.ErrNotFound
		}

		return recipe.Recipe{}, fmt.Errorf("GetRecipe failed to get recipe: %v", err)
	}

	return result, nil
}
