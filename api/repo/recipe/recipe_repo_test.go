package recipe

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

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
