package recipe

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/eciccone/rh/api/repo"
)

type RecipeRepository interface {
	InsertRecipe(recipe Recipe) (Recipe, error)
	SelectRecipeById(id int) (Recipe, error)
	SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]Recipe, error)
	UpdateRecipe(recipe Recipe) error
	UpdateRecipeImageName(id int, imageName string) error
	DeleteRecipe(id int) error
}

type recipeRepo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) RecipeRepository {
	return &recipeRepo{db}
}

// Inserts a recipe into the database.
func (r *recipeRepo) InsertRecipe(recipe Recipe) (Recipe, error) {
	var insertedRecipe Recipe

	fn := func(tx *sql.Tx) error {
		// insert into recipe table
		recipe, err := r.insertRecipe(tx, recipe)
		if err != nil {
			return err
		}

		// insert all the ingredients for the recipe
		ingredients, err := r.insertIngredients(tx, recipe.Ingredients, recipe.Id)
		if err != nil {
			return err
		}

		// set the result
		insertedRecipe.Id = recipe.Id
		insertedRecipe.Name = recipe.Name
		insertedRecipe.Username = recipe.Username
		insertedRecipe.Ingredients = ingredients

		return nil
	}

	return insertedRecipe, repo.Tx(r.db, fn)
}

// Inserts a recipe into the recipe table.
func (r *recipeRepo) insertRecipe(tx *sql.Tx, recipe Recipe) (Recipe, error) {
	result, err := tx.Exec("INSERT INTO RECIPE(name, username) VALUES (?, ?)", recipe.Name, recipe.Username)
	if err != nil {
		return Recipe{}, fmt.Errorf("recipe.InsertRecipe() failed to insert recipe: %v", err)
	}

	recipeId, _ := result.LastInsertId()
	if recipeId == 0 {
		return Recipe{}, errors.New("recipe.InsertRecipe() no id was generated for recipe")
	}

	recipe.Id = int(recipeId)

	return recipe, nil
}

// Inserts all the ingredients into the ingredient table.
func (r *recipeRepo) insertIngredients(tx *sql.Tx, ingredients []Ingredient, recipeId int) ([]Ingredient, error) {
	var result []Ingredient

	for _, ing := range ingredients {
		res, err := tx.Exec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)", ing.Name, ing.Amount, ing.Unit, recipeId)
		if err != nil {
			return nil, fmt.Errorf("recipe.InsertRecipe() failed to insert ingredient: %v", err)
		}

		ingId, _ := res.LastInsertId()
		if ingId == 0 {
			return nil, errors.New("recipe.InsertRecipe() no id was generated for ingredient")
		}

		ing.Id = int(ingId)
		result = append(result, ing)
	}

	return result, nil
}

func (r *recipeRepo) SelectRecipeById(id int) (Recipe, error) {
	return Recipe{}, nil
}

func (r *recipeRepo) SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]Recipe, error) {
	return []Recipe{}, nil
}

func (r *recipeRepo) UpdateRecipe(recipe Recipe) error {
	return nil
}

func (r *recipeRepo) UpdateRecipeImageName(id int, imageName string) error {
	return nil
}

func (r *recipeRepo) DeleteRecipe(id int) error {
	return nil
}
