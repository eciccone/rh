package service

import (
	"database/sql"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/stretchr/testify/assert"
)

type ImageServiceMocker struct {
	SaveImageMock   func() error
	DeleteImageMock func() error
}

func (s *ImageServiceMocker) SaveImage(file *multipart.FileHeader, path string, filename string) error {
	return s.SaveImageMock()
}

func (s *ImageServiceMocker) DeleteImage(path string, filename string) error {
	return s.DeleteImageMock()
}

type RecipeRepoMocker struct {
	InsertRecipeMock                func(recipe recipe.Recipe) (recipe.Recipe, error)
	SelectRecipeByIdMock            func(id int) (recipe.Recipe, error)
	SelectRecipesByUsernameMock     func(username string, orderBy string, offset int, limit int) ([]recipe.Recipe, error)
	SelectRecipeCountByUsernameMock func(username string) (int, error)
	UpdateRecipeMock                func(recipe recipe.Recipe) (recipe.Recipe, error)
	UpdateRecipeImageNameMock       func(id int, imageName string) error
	DeleteRecipeMock                func(id int) error
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

func (r *RecipeRepoMocker) SelectRecipeCountByUsername(username string) (int, error) {
	return r.SelectRecipeCountByUsernameMock(username)
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
		rs := NewRecipeService(rr, &ImageServiceMocker{})
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
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{SelectRecipeByIdMock: tr.SelectFn}
		rs := NewRecipeService(rr, &ImageServiceMocker{})
		result, err := rs.GetRecipe(tr.Input)
		tr.Assert(tr.Expected, result, err)
	}
}

func Test_GetRecipesForUsername(t *testing.T) {
	td := []struct {
		Username        string
		Expected        UsernameRecipePage
		SelectRecipesFn func(username string, orderBy string, offset int, limit int) ([]recipe.Recipe, error)
		SelectCountFn   func(username string) (int, error)
		Assert          func(expected UsernameRecipePage, actual UsernameRecipePage, err error)
	}{
		{
			Username: "Test User",
			Expected: UsernameRecipePage{
				Recipes: []recipe.Recipe{
					{Id: 1, Name: "Recipe 1", Username: "Test User"},
				},
				Offset: 0,
				Limit:  2,
				Total:  1,
			},
			SelectRecipesFn: func(username, orderBy string, offset, limit int) ([]recipe.Recipe, error) {
				return []recipe.Recipe{{Id: 1, Name: "Recipe 1", Username: "Test User"}}, nil
			},
			SelectCountFn: func(username string) (int, error) {
				return 1, nil
			},
			Assert: func(expected, actual UsernameRecipePage, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Username: "Test User",
			Expected: UsernameRecipePage{},
			SelectRecipesFn: func(username, orderBy string, offset, limit int) ([]recipe.Recipe, error) {
				return nil, errors.New("failed")
			},
			SelectCountFn: func(username string) (int, error) {
				return 1, nil
			},
			Assert: func(expected, actual UsernameRecipePage, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Username: "Test User",
			Expected: UsernameRecipePage{},
			SelectRecipesFn: func(username, orderBy string, offset, limit int) ([]recipe.Recipe, error) {
				return []recipe.Recipe{{Id: 1, Name: "Recipe 1", Username: "Test User"}}, nil
			},
			SelectCountFn: func(username string) (int, error) {
				return 0, errors.New("failed")
			},
			Assert: func(expected, actual UsernameRecipePage, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{SelectRecipesByUsernameMock: tr.SelectRecipesFn, SelectRecipeCountByUsernameMock: tr.SelectCountFn}
		rs := NewRecipeService(rr, &ImageServiceMocker{})
		result, err := rs.GetRecipesForUsername(tr.Username, "", tr.Expected.Offset, tr.Expected.Limit)
		tr.Assert(tr.Expected, result, err)
	}
}

func Test_UpdateRecipe(t *testing.T) {
	td := []struct {
		Input    recipe.Recipe
		Expected recipe.Recipe
		SelectFn func(id int) (recipe.Recipe, error)
		UpdateFn func(input recipe.Recipe) (recipe.Recipe, error)
		Assert   func(expected recipe.Recipe, actual recipe.Recipe, err error)
	}{
		{
			Input:    recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			Expected: recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"}, nil
			},
			UpdateFn: func(input recipe.Recipe) (recipe.Recipe, error) {
				return input, nil
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			Expected: recipe.Recipe{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User 1"}, nil
			},
			UpdateFn: func(input recipe.Recipe) (recipe.Recipe, error) {
				return input, nil
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			Expected: recipe.Recipe{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{}, errors.New("failed")
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Input:    recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			Expected: recipe.Recipe{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"}, nil
			},
			UpdateFn: func(input recipe.Recipe) (recipe.Recipe, error) {
				return recipe.Recipe{}, errors.New("failed")
			},
			Assert: func(expected, actual recipe.Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, actual)
			},
		},
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{SelectRecipeByIdMock: tr.SelectFn, UpdateRecipeMock: tr.UpdateFn}
		rs := NewRecipeService(rr, &ImageServiceMocker{})
		result, err := rs.UpdateRecipe(tr.Input)
		tr.Assert(tr.Expected, result, err)
	}
}

func Test_RemoveRecipe(t *testing.T) {
	td := []struct {
		Id       int
		Username string
		Expected recipe.Recipe
		SelectFn func(id int) (recipe.Recipe, error)
		DeleteFn func(id int) error
		Assert   func(err error)
	}{
		{
			Id:       1,
			Username: "Test User",
			Expected: recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"}, nil
			},
			DeleteFn: func(id int) error {
				return nil
			},
			Assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			Id:       1,
			Username: "Test User",
			Expected: recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{}, errors.New("failed")
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			Id:       1,
			Username: "Test User",
			Expected: recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User 1"}, nil
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			Id:       1,
			Username: "Test User",
			Expected: recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User"}, nil
			},
			DeleteFn: func(id int) error {
				return errors.New("failed")
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{SelectRecipeByIdMock: tr.SelectFn, DeleteRecipeMock: tr.DeleteFn}
		rs := NewRecipeService(rr, &ImageServiceMocker{})
		err := rs.RemoveRecipe(tr.Id, tr.Username)
		tr.Assert(err)
	}
}

func Test_RemoveRecipeWithImage(t *testing.T) {
	td := []struct {
		Id          int
		Username    string
		SelectFn    func(id int) (recipe.Recipe, error)
		DeleteFn    func(id int) error
		DeleteImgFn func() error
		Assert      func(err error)
	}{
		{
			Id:       1,
			Username: "Test User",
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User", ImageName: "test.file"}, nil
			},
			DeleteImgFn: func() error {
				return nil
			},
			DeleteFn: func(id int) error {
				return nil
			},
			Assert: func(err error) {
				assert.NoError(t, err)
			},
		},
		{
			Id:       1,
			Username: "Test User",
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User", ImageName: "test.file"}, nil
			},
			DeleteImgFn: func() error {
				return errors.New("failed")
			},
			DeleteFn: func(id int) error {
				return nil
			},
			Assert: func(err error) {
				assert.Error(t, err)
			},
		},
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{SelectRecipeByIdMock: tr.SelectFn, DeleteRecipeMock: tr.DeleteFn}
		rs := NewRecipeService(rr, &ImageServiceMocker{DeleteImageMock: tr.DeleteImgFn})
		err := rs.RemoveRecipe(tr.Id, tr.Username)
		tr.Assert(err)
	}
}

func Test_UpdateRecipeImage(t *testing.T) {
	td := []struct {
		Id              int
		Username        string
		MockFile        *multipart.FileHeader
		SelectFn        func(id int) (recipe.Recipe, error)
		UpdateImgNameFn func(id int, imagename string) error
		SaveImgFn       func() error
		Assert          func(result string, err error)
	}{
		{
			Id:       1,
			Username: "Test User",
			MockFile: &multipart.FileHeader{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User", ImageName: "test.file"}, nil
			},
			SaveImgFn: func() error {
				return nil
			},
			UpdateImgNameFn: func(id int, imagename string) error {
				return nil
			},
			Assert: func(result string, err error) {
				assert.NoError(t, err)
				assert.Equal(t, "test.file", result)
			},
		},
		{
			Id:       1,
			Username: "Test User",
			MockFile: &multipart.FileHeader{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User", ImageName: ""}, nil
			},
			SaveImgFn: func() error {
				return nil
			},
			UpdateImgNameFn: func(id int, imagename string) error {
				return nil
			},
			Assert: func(result string, err error) {
				assert.NoError(t, err)
				assert.NotEqual(t, "", result)
			},
		},
		{
			Id:       1,
			Username: "Test User",
			MockFile: &multipart.FileHeader{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User", ImageName: ""}, nil
			},
			SaveImgFn: func() error {
				return errors.New("failed")
			},
			UpdateImgNameFn: func(id int, imagename string) error {
				return nil
			},
			Assert: func(result string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "", result)
			},
		},
		{
			Id:       1,
			Username: "Test User",
			MockFile: &multipart.FileHeader{},
			SelectFn: func(id int) (recipe.Recipe, error) {
				return recipe.Recipe{Id: 1, Name: "Test Recipe", Username: "Test User", ImageName: ""}, nil
			},
			SaveImgFn: func() error {
				return nil
			},
			UpdateImgNameFn: func(id int, imagename string) error {
				return errors.New("failed")
			},
			Assert: func(result string, err error) {
				assert.Error(t, err)
				assert.Equal(t, "", result)
			},
		},
	}

	for _, tr := range td {
		rr := &RecipeRepoMocker{SelectRecipeByIdMock: tr.SelectFn, UpdateRecipeImageNameMock: tr.UpdateImgNameFn}
		rs := NewRecipeService(rr, &ImageServiceMocker{SaveImageMock: tr.SaveImgFn})
		result, err := rs.UpdateRecipeImage(tr.Id, tr.Username, tr.MockFile)
		tr.Assert(result, err)
	}
}
