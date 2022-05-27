package recipe

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_SelectRecipeCountByUsername(t *testing.T) {
	data := []struct {
		Name        string
		Username    string
		ExpectedSQL func(sqlmock.Sqlmock, string)
		Pass        bool
		Assert      func(sqlmock.Sqlmock, int, error)
	}{
		{
			Name:     "select recipe count",
			Username: "Test User",
			ExpectedSQL: func(m sqlmock.Sqlmock, username string) {
				m.ExpectQuery("SELECT COUNT(*) FROM recipe WHERE username = ?").WithArgs(username).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(10))
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, count int, err error) {
				assert.NoError(t, err)
				assert.Equal(t, 10, count)
			},
		},
		{
			Name:     "select recipe count error",
			Username: "Test User",
			ExpectedSQL: func(m sqlmock.Sqlmock, username string) {
				m.ExpectQuery("SELECT COUNT(*) FROM recipe WHERE username = ?").WithArgs(username).
					WillReturnError(errors.New("failed"))
			},
			Pass: false,
			Assert: func(m sqlmock.Sqlmock, count int, err error) {
				assert.Error(t, err)
				assert.Equal(t, 0, count)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		t.Log("TEST: ", d.Name)
		d.ExpectedSQL(mock, d.Username)
		rr := NewRepo(db)
		count, err := rr.SelectRecipeCountByUsername(d.Username)
		d.Assert(mock, count, err)
	}
}

func Test_DeleteRecipe(t *testing.T) {
	data := []struct {
		Name        string
		Id          int
		ExpectedSQL func(sqlmock.Sqlmock, int)
		Pass        bool
		Assert      func(sqlmock.Sqlmock, error)
	}{
		{
			Name: "delete recipe",
			Id:   1,
			ExpectedSQL: func(m sqlmock.Sqlmock, id int) {
				m.ExpectExec("DELETE FROM recipe WHERE id = ?").WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, err error) {
				assert.NoError(t, err)
			},
		},
		{
			Name: "delete recipe error",
			Id:   1,
			ExpectedSQL: func(m sqlmock.Sqlmock, id int) {
				m.ExpectExec("DELETE FROM recipe WHERE id = ?").WithArgs(id).
					WillReturnError(errors.New("failed to delete recipe"))
			},
			Pass: false,
			Assert: func(m sqlmock.Sqlmock, err error) {
				assert.Error(t, err)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		t.Log("TEST: ", d.Name)
		d.ExpectedSQL(mock, d.Id)
		rr := NewRepo(db)
		err := rr.DeleteRecipe(d.Id)
		d.Assert(mock, err)
	}
}

func Test_UpdateRecipeImageName(t *testing.T) {
	data := []struct {
		Name        string
		ImageName   string
		Id          int
		ExpectedSQL func(sqlmock.Sqlmock, int, string)
		Pass        bool
		Assert      func(sqlmock.Sqlmock, error)
	}{
		{
			Name:      "update recipe image name",
			ImageName: "test-img.jpg",
			Id:        1,
			ExpectedSQL: func(m sqlmock.Sqlmock, id int, imageName string) {
				m.ExpectExec("UPDATE recipe SET imagename = ? WHERE id = ?").WithArgs(imageName, id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, err error) {
				assert.NoError(t, err)
			},
		},
		{
			Name:      "update recipe image name error",
			ImageName: "test-img.jpg",
			Id:        1,
			ExpectedSQL: func(m sqlmock.Sqlmock, id int, imageName string) {
				m.ExpectExec("UPDATE recipe SET imagename = ? WHERE id = ?").WithArgs(imageName, id).
					WillReturnError(errors.New("failed to update recipe imagename"))
			},
			Pass: false,
			Assert: func(m sqlmock.Sqlmock, err error) {
				assert.Error(t, err)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		t.Log("TEST: ", d.Name)
		d.ExpectedSQL(mock, d.Id, d.ImageName)
		rr := NewRepo(db)
		err := rr.UpdateRecipeImageName(d.Id, d.ImageName)
		d.Assert(mock, err)
	}
}

func Test_UpdateRecipe(t *testing.T) {
	data := []struct {
		Name                string
		R                   Recipe
		ExpectedIngredients []Ingredient
		ExpectedSQL         func(sqlmock.Sqlmock, Recipe)
		Pass                bool
		Assert              func(sqlmock.Sqlmock, Recipe, Recipe, []Ingredient, error)
	}{
		{
			Name: "update recipe with old and new ingredients",
			R: Recipe{
				Id:        1,
				Name:      "Test Recipe",
				Username:  "Test User",
				ImageName: "test-img.png",
				Ingredients: []Ingredient{
					{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp", RecipeId: 1},
					{Id: 0, Name: "Ingredient 2", Amount: "1", Unit: "cups", RecipeId: 1},
				},
			},
			ExpectedIngredients: []Ingredient{
				{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp", RecipeId: 1},
				{Id: 2, Name: "Ingredient 2", Amount: "1", Unit: "cups", RecipeId: 1},
			},
			ExpectedSQL: func(m sqlmock.Sqlmock, recipe Recipe) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE recipe SET name = ?, imagename = ? WHERE id = ?").
					WithArgs(recipe.Name, recipe.ImageName, recipe.Id).WillReturnResult(sqlmock.NewResult(0, 1))

				m.ExpectExec("DELETE FROM ingredient WHERE recipeid = 1 AND id NOT IN (?)").
					WithArgs(1).WillReturnResult(sqlmock.NewResult(0, 0))

				m.ExpectExec("UPDATE ingredient SET name = ?, amount = ?, unit = ? WHERE id = ?").
					WithArgs(recipe.Ingredients[0].Name, recipe.Ingredients[0].Amount, recipe.Ingredients[0].Unit, recipe.Ingredients[0].Id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				m.ExpectExec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)").
					WithArgs(recipe.Ingredients[1].Name, recipe.Ingredients[1].Amount, recipe.Ingredients[1].Unit, recipe.Ingredients[1].RecipeId).
					WillReturnResult(sqlmock.NewResult(2, 1))
				m.ExpectCommit()
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, expected, actual Recipe, expectedI []Ingredient, err error) {
				assert.NoError(t, err)
				expected.Ingredients = expectedI
				assert.Equal(t, expected, actual)
			},
		},
		{
			Name: "update recipe with only new ingredients",
			R: Recipe{
				Id:        1,
				Name:      "Test Recipe",
				Username:  "Test User",
				ImageName: "test-img.png",
				Ingredients: []Ingredient{
					{Id: 0, Name: "Ingredient 1", Amount: "1", Unit: "tbsp", RecipeId: 1},
					{Id: 0, Name: "Ingredient 2", Amount: "1", Unit: "cups", RecipeId: 1},
				},
			},
			ExpectedIngredients: []Ingredient{
				{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp", RecipeId: 1},
				{Id: 2, Name: "Ingredient 2", Amount: "1", Unit: "cups", RecipeId: 1},
			},
			ExpectedSQL: func(m sqlmock.Sqlmock, recipe Recipe) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE recipe SET name = ?, imagename = ? WHERE id = ?").
					WithArgs(recipe.Name, recipe.ImageName, recipe.Id).WillReturnResult(sqlmock.NewResult(0, 1))

				m.ExpectExec("DELETE FROM ingredient WHERE recipeid = ?").
					WithArgs(recipe.Id).WillReturnResult(sqlmock.NewResult(0, 0))

				m.ExpectExec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)").
					WithArgs(recipe.Ingredients[0].Name, recipe.Ingredients[0].Amount, recipe.Ingredients[0].Unit, recipe.Ingredients[0].RecipeId).
					WillReturnResult(sqlmock.NewResult(1, 1))

				m.ExpectExec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)").
					WithArgs(recipe.Ingredients[1].Name, recipe.Ingredients[1].Amount, recipe.Ingredients[1].Unit, recipe.Ingredients[1].RecipeId).
					WillReturnResult(sqlmock.NewResult(2, 1))
				m.ExpectCommit()
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, expected, actual Recipe, expectedI []Ingredient, err error) {
				assert.NoError(t, err)
				expected.Ingredients = expectedI
				assert.Equal(t, expected, actual)
			},
		},
		{
			Name: "update recipe with only old ingredients",
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
			ExpectedIngredients: []Ingredient{
				{Id: 1, Name: "Ingredient 1", Amount: "1", Unit: "tbsp", RecipeId: 1},
				{Id: 2, Name: "Ingredient 2", Amount: "1", Unit: "cups", RecipeId: 1},
			},
			ExpectedSQL: func(m sqlmock.Sqlmock, recipe Recipe) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE recipe SET name = ?, imagename = ? WHERE id = ?").
					WithArgs(recipe.Name, recipe.ImageName, recipe.Id).WillReturnResult(sqlmock.NewResult(0, 1))

				m.ExpectExec("DELETE FROM ingredient WHERE recipeid = 1 AND id NOT IN (?, ?)").
					WithArgs(1, 2).WillReturnResult(sqlmock.NewResult(0, 0))

				m.ExpectExec("UPDATE ingredient SET name = ?, amount = ?, unit = ? WHERE id = ?").
					WithArgs(recipe.Ingredients[0].Name, recipe.Ingredients[0].Amount, recipe.Ingredients[0].Unit, recipe.Ingredients[0].Id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				m.ExpectExec("UPDATE ingredient SET name = ?, amount = ?, unit = ? WHERE id = ?").
					WithArgs(recipe.Ingredients[1].Name, recipe.Ingredients[1].Amount, recipe.Ingredients[1].Unit, recipe.Ingredients[1].Id).
					WillReturnResult(sqlmock.NewResult(0, 1))

				m.ExpectCommit()
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, expected, actual Recipe, expectedI []Ingredient, err error) {
				assert.NoError(t, err)
				expected.Ingredients = expectedI
				assert.Equal(t, expected, actual)
			},
		},
		{
			Name: "update recipe error",
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
			ExpectedIngredients: []Ingredient{},
			ExpectedSQL: func(m sqlmock.Sqlmock, recipe Recipe) {
				m.ExpectBegin()
				m.ExpectExec("UPDATE recipe SET name = ?, imagename = ? WHERE id = ?").
					WithArgs(recipe.Name, recipe.ImageName, recipe.Id).
					WillReturnError(errors.New("error updating recipe"))
				m.ExpectRollback()
			},
			Pass: true,
			Assert: func(m sqlmock.Sqlmock, expected, actual Recipe, expectedI []Ingredient, err error) {
				assert.Error(t, err)
				assert.Equal(t, Recipe{}, actual)
			},
		},
	}

	db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))

	for _, d := range data {
		t.Log("TEST: ", d.Name)

		d.ExpectedSQL(mock, d.R)
		rr := NewRepo(db)
		result, err := rr.UpdateRecipe(d.R)
		d.Assert(mock, d.R, result, d.ExpectedIngredients, err)
	}
}

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
				{1, "Test Name 1", "Test User", "test-img.png", nil, nil},
				{2, "Test Name 2", "Test User", "test-img.jpg", nil, nil},
				{3, "Test Name 3", "Test User", "test-img.png", nil, nil},
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
				{1, "Test Name 1", "Test User", "test-img.png", nil, nil},
				{2, "Test Name 2", "Test User", "test-img.jpg", nil, nil},
				{3, "Test Name 3", "Test User", "test-img.png", nil, nil},
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
