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
	var result Recipe

	fn := func(tx *sql.Tx) error {
		recipe, err := r.insertRecipe(tx, recipe)
		if err != nil {
			return err
		}

		ingredients, err := r.insertIngredients(tx, recipe.Ingredients, recipe.Id)
		if err != nil {
			return err
		}

		result = Recipe{recipe.Id, recipe.Name, recipe.Username, recipe.ImageName, ingredients}
		return nil
	}

	return result, repo.Tx(r.db, fn)
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
			return nil, fmt.Errorf("insertIngredients() failed to insert ingredient: %v", err)
		}

		ingId, _ := res.LastInsertId()
		if ingId == 0 {
			return nil, errors.New("insertIngredients() no id was generated for ingredient")
		}

		ing.Id = int(ingId)
		result = append(result, ing)
	}

	return result, nil
}

// Selects a recipe from the database
func (r *recipeRepo) SelectRecipeById(id int) (Recipe, error) {
	var result Recipe

	row := r.db.QueryRow("SELECT id, name, username, imagename FROM recipe WHERE id = ?", id)
	if err := row.Scan(&result.Id, &result.Name, &result.Username, &result.ImageName); err != nil {
		return Recipe{}, fmt.Errorf("recipe.SelectRecipeById() failed to select recipe: %v", err)
	}

	ingredients, err := r.selectIngredients(id)
	if err != nil {
		return Recipe{}, err
	}

	result.Ingredients = ingredients

	return result, nil
}

// Selects ingredients for a recipe from the database
func (r *recipeRepo) selectIngredients(recipeId int) ([]Ingredient, error) {
	var result []Ingredient

	rows, err := r.db.Query("SELECT id, name, amount, unit, recipeid FROM ingredient WHERE recipeid = ?", recipeId)
	if err != nil {
		return []Ingredient{}, fmt.Errorf("selectIngredients() failed to select ingredients: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var i Ingredient
		if err := rows.Scan(&i.Id, &i.Name, &i.Amount, &i.Unit, &i.RecipeId); err != nil {
			return []Ingredient{}, fmt.Errorf("selectIngredients() failed to scan row: %v", err)
		}
		result = append(result, i)
	}

	return result, nil
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
