package recipe

type Recipe struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Username    string       `json:"username"`
	ImageName   string       `json:"image"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
}

func BuildRecipe(id int, name string, username string, imageName string, ingredients []Ingredient) Recipe {
	return Recipe{id, name, username, imageName, ingredients}
}

type Ingredient struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Unit     string `json:"unit"`
	RecipeId int    `json:"-"`
}
