package recipe

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/eciccone/rh/api/repo"
)

type RecipeRepository interface {
	InsertRecipe(recipe Recipe) (Recipe, error)
	SelectRecipeById(id int) (Recipe, error)
	SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]Recipe, error)
	SelectRecipeCountByUsername(username string) (int, error)
	UpdateRecipe(recipe Recipe) (Recipe, error)
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

		steps, err := r.insertSteps(tx, recipe.Steps, recipe.Id)
		if err != nil {
			return err
		}

		result = Recipe{recipe.Id, recipe.Name, recipe.Username, recipe.ImageName, ingredients, steps}
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
		ing.RecipeId = recipeId
		result = append(result, ing)
	}

	return result, nil
}

// Inserts all the steps into the step table.
func (r *recipeRepo) insertSteps(tx *sql.Tx, steps []Step, recipeId int) ([]Step, error) {
	var result []Step

	for _, s := range steps {
		_, err := tx.Exec("INSERT INTO STEP(stepnumber, description, recipeid) VALUES(?, ?, ?)", s.StepNumber, s.Description, recipeId)
		if err != nil {
			return nil, fmt.Errorf("insertSteps failed to insert step: %v", err)
		}
		s.RecipeId = recipeId
		result = append(result, s)
	}

	return result, nil
}

// Selects a recipe from the database
func (r *recipeRepo) SelectRecipeById(id int) (Recipe, error) {
	var result Recipe

	row := r.db.QueryRow("SELECT id, name, username, imagename FROM recipe WHERE id = ?", id)
	if err := row.Scan(&result.Id, &result.Name, &result.Username, &result.ImageName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Recipe{}, err
		}

		return Recipe{}, fmt.Errorf("recipe.SelectRecipeById() failed to select recipe: %v", err)
	}

	ingredients, err := r.selectIngredients(id)
	if err != nil {
		return Recipe{}, err
	}

	steps, err := r.selectSteps(id)
	if err != nil {
		return Recipe{}, err
	}

	result.Ingredients = ingredients
	result.Steps = steps

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

func (r *recipeRepo) selectSteps(recipeId int) ([]Step, error) {
	result := []Step{}

	rows, err := r.db.Query("SELECT stepnumber, description, recipeid FROM step WHERE recipeid = ?", recipeId)
	if err != nil {
		return []Step{}, fmt.Errorf("selectSteps failed to select steps: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s Step
		if err := rows.Scan(&s.StepNumber, &s.Description, &s.RecipeId); err != nil {
			return []Step{}, fmt.Errorf("selectSteps failed to scan row: %v", err)
		}
		result = append(result, s)
	}

	return result, nil
}

// Selects a page of recipes for a user. Does not include ingredients with recipes.
func (r *recipeRepo) SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]Recipe, error) {
	var result []Recipe

	sql := "SELECT id, name, username, imagename FROM recipe WHERE username = ? ORDER BY ? LIMIT ?, ?"
	rows, err := r.db.Query(sql, username, orderBy, offset, limit)
	if err != nil {
		return []Recipe{}, fmt.Errorf("SelectRecipesByUsername() failed to select recipes: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r Recipe
		if err := rows.Scan(&r.Id, &r.Name, &r.Username, &r.ImageName); err != nil {
			return []Recipe{}, fmt.Errorf("SelectRecipesByUsername() failed to scan row: %v", err)
		}
		result = append(result, r)
	}

	return result, nil
}

func (r *recipeRepo) SelectRecipeCountByUsername(username string) (int, error) {
	rows, err := r.db.Query("SELECT COUNT(*) FROM recipe WHERE username = ?", username)
	if err != nil {
		return 0, fmt.Errorf("SelectRecipeCountByUsername failed to select count: %w", err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, fmt.Errorf("SelectRecipeCountByUsername failed to scan row: %v", err)
		}
	}

	return count, nil
}

// Updates a recipe in the database
func (r *recipeRepo) UpdateRecipe(recipe Recipe) (Recipe, error) {
	var result Recipe

	err := repo.Tx(r.db, func(tx *sql.Tx) error {
		_, err := tx.Exec("UPDATE recipe SET name = ?, imagename = ? WHERE id = ?", recipe.Name, recipe.ImageName, recipe.Id)
		if err != nil {
			return err
		}

		ingredients, err := r.updateIngredients(tx, recipe.Ingredients, recipe.Id)
		if err != nil {
			return fmt.Errorf("UpdateRecipe() failed to update ingredients: %v", err)
		}

		steps, err := r.upsertStep(tx, recipe.Steps, recipe.Id)
		if err != nil {
			return fmt.Errorf("UpdateRecipe failed to update steps: %w", err)
		}

		result = Recipe{recipe.Id, recipe.Name, recipe.Username, recipe.ImageName, ingredients, steps}

		return nil
	})

	return result, err
}

// Updates the ingredients associated with a recipe
func (r *recipeRepo) updateIngredients(tx *sql.Tx, ingredients []Ingredient, recipeId int) ([]Ingredient, error) {
	var result []Ingredient
	var existingIngredients []Ingredient
	var existingIngredientIds []interface{}
	var newIngredients []Ingredient
	for _, i := range ingredients {
		if i.Id == 0 {
			newIngredients = append(newIngredients, i)
		} else {
			existingIngredients = append(existingIngredients, i)
			existingIngredientIds = append(existingIngredientIds, i.Id)
		}
	}

	if len(existingIngredientIds) == 0 {
		// delete ingredients if none are being used
		if err := r.deleteIngredients(tx, recipeId); err != nil {
			return []Ingredient{}, fmt.Errorf("updateIngredients() failed to delete ingredients: %v", err)
		}
	} else {
		// delete ingredients that are no longer needed but keep the ones that are
		questionMarks := "?" + strings.Repeat(", ?", len(existingIngredientIds)-1)
		sql := fmt.Sprintf("DELETE FROM ingredient WHERE recipeid = %d AND id NOT IN (%s)", recipeId, questionMarks)
		_, err := tx.Exec(sql, existingIngredientIds...)
		if err != nil {
			return []Ingredient{}, fmt.Errorf("updateIngredients() error deleting unnecessary ingredients: %v", err)
		}
		// update the kept ingredients
		for _, i := range existingIngredients {
			_, err := tx.Exec("UPDATE ingredient SET name = ?, amount = ?, unit = ? WHERE id = ?", i.Name, i.Amount, i.Unit, i.Id)
			if err != nil {
				return []Ingredient{}, fmt.Errorf("updateIngredients() error updating ingredient: %v", err)
			}
			result = append(result, i)
		}
	}

	// insert the new ingredients
	i, err := r.insertIngredients(tx, newIngredients, recipeId)
	if err != nil {
		return []Ingredient{}, fmt.Errorf("updateIngredients() failed to insert ingredient: %v", err)
	}
	result = append(result, i...)

	return result, nil
}

func (r *recipeRepo) upsertStep(tx *sql.Tx, steps []Step, recipeId int) ([]Step, error) {
	if len(steps) == 0 {
		_, err := tx.Exec("DELETE FROM step WHERE recipeid = ?", recipeId)
		if err != nil {
			return []Step{}, fmt.Errorf("upsertStep failed to delete steps: %w", err)
		}
		return []Step{}, nil
	}

	var stepNumbers []interface{}

	for _, s := range steps {
		stepNumbers = append(stepNumbers, s.StepNumber)

		res, err := tx.Exec("UPDATE step SET description = ? WHERE stepnumber = ? AND recipeid = ?", s.Description, s.StepNumber, recipeId)
		if err != nil {
			return []Step{}, fmt.Errorf("upsertStep failed to update step: %w", err)
		}

		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return []Step{}, fmt.Errorf("upsertStep failed to get rows affected; %w", err)
		}
		if rowsAffected == 0 {
			_, err = tx.Exec("INSERT INTO STEP(stepnumber, description, recipeid) VALUES(?, ?, ?)", s.StepNumber, s.Description, recipeId)
			if err != nil {
				return []Step{}, fmt.Errorf("upsertStep failed to insert step: %w", err)
			}
		}
	}

	questionMarks := "?" + strings.Repeat(", ?", len(stepNumbers)-1)
	sql := fmt.Sprintf("DELETE FROM step WHERE recipeid = %d AND stepnumber NOT IN (%s)", recipeId, questionMarks)
	_, err := tx.Exec(sql, stepNumbers...)
	if err != nil {
		return []Step{}, fmt.Errorf("upsertStep failed to delete old steps: %w", err)
	}

	return steps, nil
}

// Deletes all ingredients associated with a recipe
func (r *recipeRepo) deleteIngredients(tx *sql.Tx, recipeId int) error {
	_, err := tx.Exec("DELETE FROM ingredient WHERE recipeid = ?", recipeId)
	if err != nil {
		return fmt.Errorf("deleteIngredients() failed to delete ingredient: %v", err)
	}

	return nil
}

// Updates a image name for a recipe
func (r *recipeRepo) UpdateRecipeImageName(id int, imageName string) error {
	_, err := r.db.Exec("UPDATE recipe SET imagename = ? WHERE id = ?", imageName, id)
	if err != nil {
		return fmt.Errorf("UpdateRecipeImageName() failed to update imagename: %v", err)
	}

	return nil
}

func (r *recipeRepo) DeleteRecipe(id int) error {
	_, err := r.db.Exec("DELETE FROM recipe WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("DeleteRecipe() failed to delete recipe: %v", err)
	}

	return nil
}
