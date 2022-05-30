package main

import (
	"log"
	"time"

	"github.com/eciccone/rh/api/handler"
	"github.com/eciccone/rh/api/middleware"
	"github.com/eciccone/rh/api/repo/profile"
	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/eciccone/rh/api/service"
	"github.com/eciccone/rh/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	if godotenv.Load() != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := database.Open()
	if err != nil {
		log.Fatalf("failed to open database: %s", err)
	}
	defer db.Close()

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
