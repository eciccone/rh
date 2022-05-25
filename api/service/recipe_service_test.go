package service

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/stretchr/testify/assert"
)

type RecipeRepoMocker struct {
	InsertRecipeMock            func(recipe recipe.Recipe) (recipe.Recipe, error)
	SelectRecipeByIdMock        func(id int) (recipe.Recipe, error)
	SelectRecipesByUsernameMock func(username string, orderBy string, offset int, limit int) ([]recipe.Recipe, error)
	UpdateRecipeMock            func(recipe recipe.Recipe) (recipe.Recipe, error)
	UpdateRecipeImageNameMock   func(id int, imageName string) error
	DeleteRecipeMock            func(id int) error
}

func NewMockRecipeRepo() recipe.RecipeRepository {
	return &RecipeRepoMocker{}
}

func (r *RecipeRepoMocker) InsertRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	return r.InsertRecipeMock(args)
}

func (r *RecipeRepoMocker) SelectRecipeById(id int) (recipe.Recipe, error) {
	return r.SelectRecipeByIdMock(id)
}

func (r *RecipeRepoMocker) SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]recipe.Recipe, error) {
	return r.SelectRecipesByUsernameMock(username, orderBy, offset, limit)
}

func (r *RecipeRepoMocker) UpdateRecipe(args recipe.Recipe) (recipe.Recipe, error) {
	return r.UpdateRecipeMock(args)
}

func (r *RecipeRepoMocker) UpdateRecipeImageName(id int, imageName string) error {
	return r.UpdateRecipeImageNameMock(id, imageName)
}

func (r *RecipeRepoMocker) DeleteRecipe(id int) error {
	return r.DeleteRecipeMock(id)
}

func Test_CreateRecipe(t *testing.T) {
	td := []struct {
		Input    recipe.Recipe
		Expected recipe.Recipe
		InsertFn func(args recipe.Recipe) (recipe.Recipe, error)
		Assert   func(expected recipe.Recipe, actual recipe.Recipe, err error)
	}{
		{
			Input:    recipe.Recipe{Name: "Test Name", Username: "Test User"},
			Expected: recipe.Recipe{Id: 1, Name: "Test Name", Username: "Test User"},
			InsertFn: func(args recipe.Recipe) (recipe.Recipe, error) {
				args.Id = 1
				return args, nil
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    recipe.Recipe{Name: "Test Name", Username: "Test User"},
			Expected: recipe.Recipe{},
			InsertFn: func(args recipe.Recipe) (recipe.Recipe, error) {
				return recipe.Recipe{}, errors.New("failed")
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    recipe.Recipe{Name: "", Username: ""},
			Expected: recipe.Recipe{},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{InsertRecipeMock: tr.InsertFn}
		rs := NewRecipeService(rr)
		result, err := rs.CreateRecipe(tr.Input)
		tr.Assert(tr.Expected, result, err)
	}
}

func Test_GetRecipe(t *testing.T) {
	td := []struct {
		Input    int
		Expected recipe.Recipe
		SelectFn func(id int) (recipe.Recipe, error)
		Assert   func(expected recipe.Recipe, actual recipe.Recipe, err error)
	}{
		{
			Input:    1,
			Expected: recipe.Recipe{Id: 1, Name: "Test Name", Username: "Test User"},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Name", Username: "Test User"}, nil
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    1,
			Expected: recipe.Recipe{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{}, errors.New("fail")
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    1,
			Expected: recipe.Recipe{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{}, sql.ErrNoRows
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    0,
			Expected: recipe.Recipe{},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{SelectRecipeByIdMock: tr.SelectFn}
		rs := NewRecipeService(rr)
		result, err := rs.GetRecipe(tr.Input)
		tr.Assert(tr.Expected, result, err)
	}
}
