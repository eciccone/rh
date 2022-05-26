package service

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/eciccone/rh/api/repo/recipe"
)

var (
	ErrRecipeData      = errors.New("must provide name for recipe")
	ErrIngredientData  = errors.New("must provide name, amount, and unit for ingredient")
	ErrNoRecipe        = errors.New("recipe not found")
	ErrRecipeForbidden = errors.New("recipe access not allowed")
)

type RecipeService interface {
	// Creates a new recipe.
	// Returns ErrRecipeData if recipe name is empty.
	CreateRecipe(recipe.Recipe) (recipe.Recipe, error)

	// Gets a recipe by id.
	// Returns ErrNoRecipe if recipe does not exist.
	GetRecipe(id int) (recipe.Recipe, error)

	// Gets a page of recipes given the username, order (defaults to id desc), offset and limit.
	GetRecipesForUsername(username string, orderBy string, offset int, limit int) (UsernameRecipePage, error)

	// Updates a recipe.
	// Returns ErrRecipeData if recipe name is empty.
	// Returns ErrNoRecipe if recipe does not exist.
	// Returns ErrRecipeForbidden if recipe does not belong to user.
	UpdateRecipe(args recipe.Recipe) (recipe.Recipe, error)

	// Removes a recipe.
	// Returns ErrNoRecipe if recipe does not exist.
	// Returns ErrRecipeForbidden if recipe does not belong to user.
	RemoveRecipe(id int, username string) error
}

type recipeService struct {
	recipeRepo   recipe.RecipeRepository
	imageService ImageService
}

func NewRecipeService(recipeRepo recipe.RecipeRepository, imageService ImageService) RecipeService {
	return &recipeService{recipeRepo, imageService}
}

// Creates a new recipe.
// Returns ErrRecipeData if recipe name is empty.
func (s *recipeService) CreateRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	if args.Name == "" {
		return recipe.Recipe{}, ErrRecipeData
	}

	result, err := s.recipeRepo.InsertRecipe(args)
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("CreateRecipe failed to create recipe: %w", err)
	}

	return result, nil
}

// Gets a recipe by id.
// Returns ErrNoRecipe if recipe does not exist.
func (s *recipeService) GetRecipe(id int) (recipe.Recipe, error) {
	result, err := s.recipeRepo.SelectRecipeById(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return recipe.Recipe{}, ErrNoRecipe
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

// Gets a page of recipes given the username, order (defaults to id desc), offset and limit.
func (s *recipeService) GetRecipesForUsername(username string, orderBy string, offset int, limit int) (UsernameRecipePage, error) {
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

// Updates a recipe.
// Returns ErrRecipeData if recipe name is empty.
// Returns ErrNoRecipe if recipe does not exist.
// Returns ErrRecipeForbidden if recipe does not belong to user.
func (s *recipeService) UpdateRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	if args.Name == "" {
		return recipe.Recipe{}, ErrRecipeData
	}

	// make sure recipe exists
	old, err := s.GetRecipe(args.Id)
	if err != nil {
		return old, err
	}

	if old.Username != args.Username {
		return recipe.Recipe{}, ErrRecipeForbidden
	}

	// don't update imagename, seperate func for this
	args.ImageName = old.ImageName

	result, err := s.recipeRepo.UpdateRecipe(args)
	if err != nil {
		return recipe.Recipe{}, fmt.Errorf("UpdateRecipe failed to update recipe: %w", err)
	}

	return result, nil
}

// Removes a recipe.
// Returns ErrNoRecipe if recipe does not exist.
// Returns ErrRecipeForbidden if recipe does not belong to user.
func (s *recipeService) RemoveRecipe(id int, username string) error {
	// select recipe by id to make sure it exists
	r, err := s.GetRecipe(id)
	if err != nil {
		return err
	}

	// make sure user deleting the recipe owns the recipe
	if r.Username != username {
		return ErrRecipeForbidden
	}

	if r.ImageName != "" {
		err := s.imageService.DeleteImage(os.Getenv("IMAGE_PATH"), r.ImageName)
		if err != nil {
			return err
		}
	}

	err = s.recipeRepo.DeleteRecipe(id)
	if err != nil {
		return fmt.Errorf("RemoveRecipe failed to delete recipe: %w", err)
	}

	return nil
}
