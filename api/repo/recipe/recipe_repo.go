package recipe

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/eciccone/rh/api/repo"
)

type RecipeRepository interface {
	InsertRecipe(recipe Recipe) (int, error)
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

func (r *recipeRepo) InsertRecipe(recipe Recipe) (int, error) {
	id := 0

	fn := func(tx *sql.Tx) error {
		// insert into recipe table
		result, err := tx.Exec("INSERT INTO RECIPE(name, username) VALUES (?, ?)", recipe.Name, recipe.Username)
		if err != nil {
			return fmt.Errorf("recipe.InsertRecipe() failed to insert recipe: %v", err)
		}

		// get the id of the newly inserted recipe
		recipeId, _ := result.LastInsertId()
		if recipeId == 0 {
			return errors.New("recipe.InsertRecipe() no id was generated for recipe")
		}

		// insert all the ingredients for the recipe
		for _, ing := range recipe.Ingredients {
			_, err := tx.Exec("INSERT INTO INGREDIENT(name, amount, unit, recipeid) VALUES(?, ?, ?, ?)",
				ing.Name, ing.Amount, ing.Unit, recipeId)
			if err != nil {
				return fmt.Errorf("recipe.InsertRecipe() failed to insert ingredient: %v", err)
			}
		}

		// set the id to be returned
		id = int(recipeId)

		return nil
	}

	return id, repo.Tx(r.db, fn)
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
