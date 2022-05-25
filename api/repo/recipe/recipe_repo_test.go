package recipe

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_SelectRecipesByUsername(t *testing.T) {
	data := []struct {
		Name        string
		R           []Recipe
		Username    string
		ExpectedSQL func(sqlmock.Sqlmock, []Recipe, string)
		Pass        bool
		Assert      func(sqlmock.Sqlmock, []Recipe, []Recipe, error)
	}{
		{
			Name: "select recipes by username",
			R: []Recipe{
				{1, "Test Name 1", "Test User", "test-img.png", nil},
				{2, "Test Name 2", "Test User", "test-img.jpg", nil},
				{3, "Test Name 3", "Test User", "test-img.png", nil},
			},
			Username: "Test User",
			ExpectedSQL: func(m sqlmock.Sqlmock, r []Recipe, username string) {
				recipeRow := sqlmock.NewRows([]string{"id", "name", "username", "imagename"})
				for _, rr := range r {
					recipeRow.AddRow(rr.Id, rr.Name, rr.Username, rr.ImageName)
				}

				m.ExpectQuery("SELECT id, name, username, imagename FROM recipe WHERE username = ? ORDER BY ? LIMIT ?, ?").
					WithArgs(username, "id desc", 0, 10).WillReturnRows(recipeRow)
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, expected, actual []Recipe, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual)
			},
		},
		{
			Name: "select recipes by username error",
			R: []Recipe{
				{1, "Test Name 1", "Test User", "test-img.png", nil},
				{2, "Test Name 2", "Test User", "test-img.jpg", nil},
				{3, "Test Name 3", "Test User", "test-img.png", nil},
			},
			Username: "Test User",
			ExpectedSQL: func(m sqlmock.Sqlmock, r []Recipe, username string) {
				m.ExpectQuery("SELECT id, name, username, imagename FROM recipe WHERE username = ? ORDER BY ? LIMIT ?, ?").
					WithArgs(username, "id desc", 0, 10).WillReturnError(errors.New("error selecting recipes by username"))
			},
			Pass: false,
			Assert: func(m sqlmock.Sqlmock, expected, actual []Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, []Recipe{}, actual)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		t.Log("TEST: ", d.Name)

		d.ExpectedSQL(mock, d.R, d.Username)
		rr := NewRepo(db)
		result, err := rr.SelectRecipesByUsername(d.Username, "id desc", 0, 10)
		d.Assert(mock, d.R, result, err)
	}
}

func Test_SelectRecipeById(t *testing.T) {
	data := []struct {
		Name        string
		R           Recipe
		ExpectedSQL func(sqlmock.Sqlmock, Recipe)
		Pass        bool
		Assert      func(sqlmock.Sqlmock, Recipe, Recipe, error)
	}{
		{
			Name: "select recipe",
			R: Recipe{
				Id:        1,
				Name:      "Test Recipe",
				Username:  "Test User",
				ImageName: "test-img.png",
				Ingredients: []Ingredient{
					{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp", RecipeId: 1},
					{Id: 2, Name: "Ingredient 2", Amount: "1", Unit: "cups", RecipeId: 1},
				},
			},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				recipeRow := sqlmock.NewRows([]string{"id", "name", "username", "imagename"}).
					AddRow(recipe.Id, recipe.Name, recipe.Username, recipe.ImageName)
				mock.ExpectQuery("SELECT id, name, username, imagename FROM recipe WHERE id = ?").
					WithArgs(recipe.Id).WillReturnRows(recipeRow)

				ingredientRows := sqlmock.NewRows([]string{"id", "name", "amount", "unit", "recipeid"})
				for _, i := range recipe.Ingredients {
					ingredientRows.AddRow(i.Id, i.Name, i.Amount, i.Unit, i.RecipeId)
				}
				mock.ExpectQuery("SELECT id, name, amount, unit, recipeid FROM ingredient WHERE recipeid = ?").
					WithArgs(recipe.Id).
					WillReturnRows(ingredientRows)
			},
			Pass: true,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, result)
			},
		},
		{
			Name: "select recipe error",
			R: Recipe{
				Id:          1,
				Ingredients: []Ingredient{},
			},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				mock.ExpectQuery("SELECT id, name, username, imagename FROM recipe WHERE id = ?").
					WithArgs(recipe.Id).WillReturnError(errors.New("error selecting recipe"))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, Recipe{}, result)
			},
		},
		{
			Name: "select recipe ingredient error",
			R: Recipe{
				Id:        1,
				Name:      "Test Recipe",
				Username:  "Test User",
				ImageName: "test-img.png",
				Ingredients: []Ingredient{
					{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp", RecipeId: 1},
					{Id: 2, Name: "Ingredient 2", Amount: "1", Unit: "cups", RecipeId: 1},
				},
			},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				recipeRow := sqlmock.NewRows([]string{"id", "name", "username", "imagename"}).
					AddRow(recipe.Id, recipe.Name, recipe.Username, recipe.ImageName)
				mock.ExpectQuery("SELECT id, name, username, imagename FROM recipe WHERE id = ?").
					WithArgs(recipe.Id).WillReturnRows(recipeRow)

				mock.ExpectQuery("SELECT id, name, amount, unit, recipeid FROM ingredient WHERE recipeid = ?").
					WithArgs(recipe.Id).WillReturnError(errors.New("error selecting ingredients"))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, Recipe{}, result)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		t.Log("TEST: ", d.Name)

		d.ExpectedSQL(mock, d.R)
		rr := NewRepo(db)
		result, err := rr.SelectRecipeById(d.R.Id)
		d.Assert(mock, d.R, result, err)
	}
}

func Test_InsertRecipe(t *testing.T) {
	data := []struct {
		Name        string
		R           Recipe
		ExpectedSQL func(sqlmock.Sqlmock, Recipe)
		Pass        bool
		Assert      func(sqlmock.Sqlmock, Recipe, Recipe, error)
	}{
		{
			Name: "insert recipe",
			R: Recipe{
				Id:       1,
				Name:     "Test Recipe",
				Username: "Test User",
				Ingredients: []Ingredient{
					{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp"},
					{Id: 2, Name: "Ingredient 2", Amount: "1", Unit: "cups"},
				},
			},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				mock.ExpectExec("INSERT INTO RECIPE(name, username) VALUES (?, ?)").
					WithArgs(recipe.Name, recipe.Username).
					WillReturnResult(sqlmock.NewResult(int64(recipe.Id), 1))
				for _, in := range recipe.Ingredients {
					mock.ExpectExec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)").
						WithArgs(in.Name, in.Amount, in.Unit, recipe.Id).
						WillReturnResult(sqlmock.NewResult(int64(in.Id), 1))
				}
			},
			Pass: true,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.NoError(t, err)
				assert.Equal(t, expected, result)
			},
		},
		{
			Name: "insert recipe no generated id",
			R:    Recipe{},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				mock.ExpectExec("INSERT INTO RECIPE(name, username) VALUES (?, ?)").
					WithArgs(recipe.Name, recipe.Username).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, result)
			},
		},
		{
			Name: "insert recipe error",
			R:    Recipe{},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				mock.ExpectExec("INSERT INTO RECIPE(name, username) VALUES (?, ?)").
					WithArgs(recipe.Name, recipe.Username).
					WillReturnError(errors.New("error inserting recipe"))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, expected, result)
			},
		},
		{
			Name: "insert recipe ingredient error",
			R: Recipe{
				Id:       1,
				Name:     "Test Recipe",
				Username: "Test User",
				Ingredients: []Ingredient{
					{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp"},
				},
			},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				mock.ExpectExec("INSERT INTO RECIPE(name, username) VALUES (?, ?)").
					WithArgs(recipe.Name, recipe.Username).
					WillReturnResult(sqlmock.NewResult(int64(recipe.Id), 1))
				mock.ExpectExec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)").
					WithArgs(recipe.Ingredients[0].Name, recipe.Ingredients[0].Amount, recipe.Ingredients[0].Unit, recipe.Id).
					WillReturnError(errors.New("error inserting ingredient"))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, Recipe{}, result)
			},
		},
		{
			Name: "insert recipe ingredient no generated id",
			R: Recipe{
				Id:       1,
				Name:     "Test Recipe",
				Username: "Test User",
				Ingredients: []Ingredient{
					{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp"},
				},
			},
			ExpectedSQL: func(mock sqlmock.Sqlmock, recipe Recipe) {
				mock.ExpectExec("INSERT INTO RECIPE(name, username) VALUES (?, ?)").
					WithArgs(recipe.Name, recipe.Username).
					WillReturnResult(sqlmock.NewResult(int64(recipe.Id), 1))
				mock.ExpectExec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)").
					WithArgs(recipe.Ingredients[0].Name, recipe.Ingredients[0].Amount, recipe.Ingredients[0].Unit, recipe.Id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Pass: false,
			Assert: func(mock sqlmock.Sqlmock, expected, result Recipe, err error) {
				assert.Error(t, err)
				assert.Equal(t, Recipe{}, result)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		t.Log("TEST: ", d.Name)

		mock.ExpectBegin()
		d.ExpectedSQL(mock, d.R)

		if d.Pass {
			mock.ExpectCommit()
		} else {
			mock.ExpectRollback()
		}

		rr := NewRepo(db)
		result, err := rr.InsertRecipe(d.R)

		d.Assert(mock, d.R, result, err)
	}
}
