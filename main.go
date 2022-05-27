package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/eciccone/rh/api/handler"
	"github.com/eciccone/rh/api/middleware"
	"github.com/eciccone/rh/api/repo/profile"
	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/eciccone/rh/api/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

var (
	dbfile = "recihub.db"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	connName := fmt.Sprintf("%v?_foreign_keys=on", dbfile)
	db, err := sql.Open("sqlite3", connName)
	if err != nil {
		log.Fatal("failed to connect to database: " + dbfile)
	}
	defer db.Close()
	CreateSQLiteTables(db)

	pr := profile.NewRepo(db)
	rr := recipe.NewRepo(db)

	ps := service.NewProfileService(pr)
	is := service.NewFileProcessor()
	rs := service.NewRecipeService(rr, is)

	ph := handler.NewProfileHandler(ps)
	rh := handler.NewRecipeHandler(rs)

	router := setupRouter()

	router.Use(middleware.Validate())
	router.Static("/static/images", "./static/images")

	router.GET("/profile", handler.Handler(ph.GetProfile))
	router.POST("/profile", handler.Handler(ph.PostProfile))

	router.Use(middleware.Profile(ps))
	router.GET("/recipes/:id", handler.Handler(rh.GetRecipe))
	router.GET("/recipes", handler.Handler(rh.GetRecipes))
	router.POST("/recipes", handler.Handler(rh.PostRecipe))
	router.PUT("/recipes/:id", handler.Handler(rh.PutRecipe))
	router.PUT("/recipes/:id/image", handler.Handler(rh.PutRecipeImage))
	router.DELETE("/recipes/:id", handler.Handler(rh.DeleteRecipe))

	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.New()

	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))

	return r
}

const createProfileTable = `
	CREATE TABLE IF NOT EXISTS profile (
  	id TEXT NOT NULL PRIMARY KEY,
  	username TEXT NOT NULL UNIQUE
  );`

const createRecipeTable = `
	CREATE TABLE IF NOT EXISTS recipe (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		username TEXT NOT NULL,
		imagename TEXT default "",
		CHECK (name <> '' AND username <> '')
	);`

const createIngredientTable = `
	CREATE TABLE IF NOT EXISTS ingredient (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		amount TEXT NOT NULL,
		unit TEXT NOT NULL,
		recipeid INTEGER NOT NULL,
		FOREIGN KEY(recipeid) REFERENCES recipe(id) ON DELETE CASCADE
	);`

const createStepTable = `
	CREATE TABLE IF NOT EXISTS step (
		stepnumber INTEGER NOT NULL,
		description TEXT NOT NULL,
		recipeid INTEGER NOT NULL,
		PRIMARY KEY(stepnumber, recipeid),
		FOREIGN KEY(recipeid) REFERENCES recipe(id) ON DELETE CASCADE
	);`

func CreateSQLiteTables(conn *sql.DB) {
	if _, err := conn.Exec(createProfileTable); err != nil {
		log.Fatal("failed to create PROFILE table")
	}

	if _, err := conn.Exec(createRecipeTable); err != nil {
		log.Fatal("failed to create RECIPE table")
	}

	if _, err := conn.Exec(createIngredientTable); err != nil {
		log.Fatal("failed to create INGREDIENT table")
	}

	if _, err := conn.Exec(createStepTable); err != nil {
		log.Fatal("failed to create STEP table")
	}
}
