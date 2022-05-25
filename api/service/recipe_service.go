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
	GetRecipesForUsername(username string, orderBy string, offset int, limit int) (UsernameRecipePage, error)
	UpdateRecipe(args recipe.Recipe) (recipe.Recipe, error)
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
		return recipe.Recipe{}, fmt.Errorf("CreateRecipe failed to create recipe: %w", err)
	}

	return result, nil
}

func (s *recipeService) GetRecipe(id int) (recipe.Recipe, error) {
	if id <= 0 {
		return recipe.Recipe{}, rherr.ErrBadRequest
	}

	result, err := s.recipeRepo.SelectRecipeById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recipe.Recipe{}, rherr.ErrNotFound
		}

		return recipe.Recipe{}, fmt.Errorf("GetRecipe failed to get recipe: %w", err)
	}

	return result, nil
}

type UsernameRecipePage struct {
	Recipes []recipe.Recipe `json:"recipes"`
	Offset  int             `json:"offset"`
	Limit   int             `json:"limit"`
	Total   int             `json:"total"`
}

func (s *recipeService) GetRecipesForUsername(username string, orderBy string, offset int, limit int) (UsernameRecipePage, error) {
	if username == "" {
		return UsernameRecipePage{}, rherr.ErrBadRequest
	}

	if orderBy == "" {
		orderBy = "id desc"
	}

	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 10
	}

	recipes, err := s.recipeRepo.SelectRecipesByUsername(username, orderBy, offset, limit)
	if err != nil {
		return UsernameRecipePage{}, fmt.Errorf("GetRecipesForUsername failed to get recipes for username: %w", err)
	}

	total, err := s.recipeRepo.SelectRecipeCountByUsername(username)
	if err != nil {
		return UsernameRecipePage{}, fmt.Errorf("GetRecipesForUsername failed to get total recipe count: %w", err)
	}

	return UsernameRecipePage{
		Recipes: recipes,
		Offset:  offset,
		Limit:   limit,
		Total:   total,
	}, nil
}

func (s *recipeService) UpdateRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	if args.Name == "" || args.Username == "" {
		return recipe.Recipe{}, rherr.ErrBadRequest
	}

	// make sure recipe exists
	old, err := s.GetRecipe(args.Id)
	if err != nil {
		return old, err
	}

	if old.Username != args.Username {
		return recipe.Recipe{}, rherr.ErrForbidden
	}

	// don't update imagename, seperate func for this
	args.ImageName = old.ImageName

	result, err := s.recipeRepo.UpdateRecipe(args)
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("UpdateRecipe failed to update recipe: %w", err)
	}

	return result, nil
}
