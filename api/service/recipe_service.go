package service

import (
	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/eciccone/rh/api/rherr"
)

type RecipeService interface {
	CreateRecipe(recipe.Recipe) (recipe.Recipe, error)
}

type recipeService struct {
	recipeRepo recipe.RecipeRepository
}

func NewRecipeService(recipeRepo recipe.RecipeRepository) RecipeService {
	return &recipeService{recipeRepo}
}

func (s *recipeService) CreateRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	if args.Name == "" || args.Username == "" {
		return recipe.Recipe{}, rherr.ErrBadRequest
	}

	result, err := s.recipeRepo.InsertRecipe(args)
	if err != nil {
		return recipe.Recipe{}, rherr.ErrInternal
	}

	return result, nil
}
