package service

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/eciccone/rh/api/repo/recipe"
)

var (
	ErrRecipeData      = errors.New("must provide name for recipe")
	ErrIngredientData  = errors.New("must provide name, amount, and unit for ingredient")
	ErrNoRecipe        = errors.New("recipe not found")
	ErrRecipeForbidden = errors.New("recipe access not allowed")
)

type RecipeService interface {
	CreateRecipe(recipe.Recipe) (recipe.Recipe, error)
	GetRecipe(id int) (recipe.Recipe, error)
	GetRecipesForUsername(username string, orderBy string, offset int, limit int) (UsernameRecipePage, error)
	UpdateRecipe(args recipe.Recipe) (recipe.Recipe, error)
	RemoveRecipe(id int, username string) error
}

type recipeService struct {
	recipeRepo recipe.RecipeRepository
}

func NewRecipeService(recipeRepo recipe.RecipeRepository) RecipeService {
	return &recipeService{recipeRepo}
}

// Creates a new recipe.
// If no name is set for the recipe ErrRecipeData is returned.
// If recipe failed to be inserted, an error is returned.
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
// If no row is returned from the database ErrNoRecipe is returned.
// If recipe failed to be selected from database, an error is returned.
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

// Gets a page of recipes given the username, order (defaults to id desc), offset and limit. If recipes
// fail to be selected an error is returned. If getting the total amount of recipes for the user fails,
// an error is returned.
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
// If no name is set for the recipe ErrRecipeData is returned.
// If no row is returnedmfrom the database when selecting the recipe ErrNoRecipe is returned.
// If recipe failed to be selected from database, an error is returned.
// If the updated recipe does not belong to the user ErrRecipeForbidden is returned.
// If updating recipe fails, an error is returned.
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
// If no row is returned from the database when selecting the recipe ErrNoRecipe is returned.
// If recipe failed to be selected an error is returned.
// If recipe being deleted does not belong to user ErrRecipeForbidden is returned.
// If recipe fails to be deleted an error is returned.
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
		// err := DeleteImage("./images", r.ImageName)
		// if err != nil {
		// 	return err
		// }
	}

	err = s.recipeRepo.DeleteRecipe(id)
	if err != nil {
		return fmt.Errorf("RemoveRecipe failed to delete recipe: %w", err)
	}

	return nil
}
