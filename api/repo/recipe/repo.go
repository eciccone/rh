package recipe

import "database/sql"

type RecipeRepository interface {
	InsertRecipe(recipe Recipe) (int, error)
	SelectRecipeById(id int) (Recipe, error)
	SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]Recipe, error)
	UpdateRecipe(recipe Recipe) error
	UpdateRecipeImageName(id int, imageName string) error
	DeleteRecipe(id int) error
}

type repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) RecipeRepository {
	return &repo{db}
}

func (r *repo) InsertRecipe(recipe Recipe) (int, error) {
	return 0, nil
}

func (r *repo) SelectRecipeById(id int) (Recipe, error) {
	return Recipe{}, nil
}

func (r *repo) SelectRecipesByUsername(username string, orderBy string, offset int, limit int) ([]Recipe, error) {
	return []Recipe{}, nil
}

func (r *repo) UpdateRecipe(recipe Recipe) error {
	return nil
}

func (r *repo) UpdateRecipeImageName(id int, imageName string) error {
	return nil
}

func (r *repo) DeleteRecipe(id int) error {
	return nil
}
