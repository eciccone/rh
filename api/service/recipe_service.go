package service

import (
	"database/sql"
	"errors"

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
		return recipe.Recipe{}, rherr.ErrInternal
	}

	return result, nil
}

func (s *recipeService) GetRecipe(id int) (recipe.Recipe, error) {
	result, err := s.recipeRepo.SelectRecipeById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recipe.Recipe{}, rherr.ErrNotFound
		}

		return recipe.Recipe{}, rherr.ErrInternal
	}

	return result, nil
}
