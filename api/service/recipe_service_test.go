package service

import (
	"errors"
	"testing"

	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/stretchr/testify/assert"
)

type mockRecipeRepo struct{}

func NewMockRecipeRepo() recipe.RecipeRepository {
	return &mockRecipeRepo{}
}

func (r *mockRecipeRepo) InsertRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	if args.Id == -1 {
		return recipe.Recipe{}, errors.New("failed to insert recipe")
	}

	return recipe.Recipe{
		Id:       1,
		Name:     args.Name,
		Username: args.Username,
	}, nil
}

func (r *mockRecipeRepo) SelectRecipeById(id int) (recipe.Recipe, error) {
	var result recipe.Recipe

	return result, nil
}

func (r *mockRecipeRepo) SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]recipe.Recipe, error) {
	var result []recipe.Recipe

	return result, nil
}

func (r *mockRecipeRepo) UpdateRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	var result recipe.Recipe

	return result, nil
}

func (r *mockRecipeRepo) UpdateRecipeImageName(id int, imageName string) error {
	return nil
}

func (r *mockRecipeRepo) DeleteRecipe(id int) error {
	return nil
}

func Test_CreateRecipe(t *testing.T) {
	data := []struct {
		Id          int
		Name        string
		Username    string
		Ingredients []recipe.Ingredient
		Test        func(string, string, []recipe.Ingredient)
		Assert      func(recipe.Recipe, recipe.Recipe, error)
	}{
		{
			Name:     "Test Name",
			Username: "Test User",
			Ingredients: []recipe.Ingredient{
				{Name: "ingredient 1", Amount: "2", Unit: "tbsp"},
			},
			Assert: func(r1, r2 recipe.Recipe, err error) {
				assert.NoError(t, err)
				assert.NotZero(t, r2.Id)
				assert.Equal(t, r1.Name, r2.Name)
				assert.Equal(t, r1.Username, r2.Username)
			},
		},
		{
			Name:     "",
			Username: "Test User",
			Ingredients: []recipe.Ingredient{
				{Name: "ingredient 1", Amount: "2", Unit: "tbsp"},
			},
			Assert: func(r1, r2 recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Empty(t, r2)
			},
		},
		{
			Name:     "Test Name",
			Username: "",
			Ingredients: []recipe.Ingredient{
				{Name: "ingredient 1", Amount: "2", Unit: "tbsp"},
			},
			Assert: func(r1, r2 recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Empty(t, r2)
			},
		},
		{
			Id:       -1,
			Name:     "Test Name",
			Username: "Test User",
			Ingredients: []recipe.Ingredient{
				{Name: "ingredient 1", Amount: "2", Unit: "tbsp"},
			},
			Assert: func(r1, r2 recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Empty(t, r2)
			},
		},
	}

	for _, d := range data {
		rs := NewRecipeService(NewMockRecipeRepo())
		r1 := recipe.Recipe{
			Id:          d.Id,
			Name:        d.Name,
			Username:    d.Username,
			Ingredients: d.Ingredients,
		}

		r2, err := rs.CreateRecipe(r1)
		d.Assert(r1, r2, err)
	}
}
