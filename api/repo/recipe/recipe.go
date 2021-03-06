package recipe

type Recipe struct {
	Id          int          `json:"id"`
	Name        string       `json:"name"`
	Username    string       `json:"username"`
	ImageName   string       `json:"image"`
	Ingredients []Ingredient `json:"ingredients,omitempty"`
	Steps       []Step       `json:"steps,omitempty"`
}

type Ingredient struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Amount   string `json:"amount"`
	Unit     string `json:"unit"`
	RecipeId int    `json:"-"`
}

type Step struct {
	StepNumber  int    `json:"step_number"`
	Description string `json:"description"`
	RecipeId    int    `json:"-"`
}
