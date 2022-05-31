package router

import (
	"database/sql"
	"time"

	"github.com/eciccone/rh/api/handler"
	"github.com/eciccone/rh/api/middleware"
	"github.com/eciccone/rh/api/repo/profile"
	"github.com/eciccone/rh/api/repo/recipe"
	"github.com/eciccone/rh/api/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func New() Router {
	r := Router{
		Engine: gin.New(),
	}

	r.setConfiguration()

	return r
}

func (r *Router) setConfiguration() {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowHeaders:     []string{"Access-Control-Allow-Origin", "*"},
		AllowMethods:     []string{"Access-Control-Allow-Methods", "*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Engine.Use(cors.New(config))
}

func (r *Router) Run(addr string) {
	r.Engine.Run(addr)
}

func (r *Router) BuildRoutes(db *sql.DB) {
	pr := profile.NewRepo(db)
	rr := recipe.NewRepo(db)

	ps := service.NewProfileService(pr)
	is := service.NewFileProcessor()
	rs := service.NewRecipeService(rr, is)

	ph := handler.NewProfileHandler(ps)
	rh := handler.NewRecipeHandler(rs)

	// all end points below must have a valid access token
	r.Engine.Use(middleware.Validate())

	// location for recipe image uploads
	r.Engine.Static("/static/images", "./static/images")

	// profile routes
	r.Engine.GET("/profile", handler.Handler(ph.GetProfile))
	r.Engine.POST("/profile", handler.Handler(ph.PostProfile))

	// all end points below must have already created a profile with recihub
	r.Engine.Use(middleware.Profile(ps))

	// recipe routes
	r.Engine.GET("/recipes/:id", handler.Handler(rh.GetRecipe))
	r.Engine.GET("/recipes", handler.Handler(rh.GetRecipes))
	r.Engine.POST("/recipes", handler.Handler(rh.PostRecipe))
	r.Engine.PUT("/recipes/:id", handler.Handler(rh.PutRecipe))
	r.Engine.PUT("/recipes/:id/image", handler.Handler(rh.PutRecipeImage))
	r.Engine.DELETE("/recipes/:id", handler.Handler(rh.DeleteRecipe))
}
