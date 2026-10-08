package server

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/hrnnsx/libra/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	//* NOTES: ada kemungkinan untuk penambahan role, jadi sementara menggunakan multiple routing chains

	// health
	r.GET("/", s.HelloWorldHandler)
	r.GET("/health", s.healthHandler)

	protected := r.Group("")

	// auth
	auth := r.Group("/auth")
	{
		auth.POST("/register", s.authHandler.Register)
		auth.POST("/login", s.authHandler.Login)
	}

	// user
	user := protected.Group("/users")
	user.Use(middleware.AuthMiddleware())
	{
		user.GET("/me", s.userHandler.GetProfile)
		user.PATCH("/me", s.userHandler.UpdateProfile)
	}

	// anime
	animes := r.Group("/animes")
	{
		animes.GET("", s.animeHandler.BrowseAnime)
		animes.GET("/search", s.animeHandler.SearchAnime)
		animes.GET("/:external_id", s.animeHandler.GetAnime)
		animes.POST("/identify", s.animeHandler.IdentifyAnime)
	}

	// library
	library := r.Group("/me/library")
	library.Use(middleware.AuthMiddleware())
	{
		library.POST("/", s.libraryHandler.AddAnime)
		library.GET("/", s.libraryHandler.GetLibrary)

		library.GET("/:id", s.libraryHandler.GetAnime)
		library.PATCH("/:id", s.libraryHandler.UpdateAnime)
		library.DELETE("/:id", s.libraryHandler.DeleteAnime)
	}

	// groups
	groups := r.Group("/me/groups")
	groups.Use(middleware.AuthMiddleware())
	{
		groups.POST("", s.groupHandler.CreateGroup)
		groups.GET("", s.groupHandler.GetGroups)

		groups.GET("/:id", s.groupHandler.GetGroup)
		groups.PATCH("/:id", s.groupHandler.UpdateGroup)
		groups.DELETE("/:id", s.groupHandler.DeleteGroup)

		groups.POST("/:id/animes", s.groupLibraryAnimeHandler.AddAnime)
		groups.GET("/:id/animes", s.groupLibraryAnimeHandler.GetAnimes)
		groups.DELETE("/:id/animes/:library_anime_id", s.groupLibraryAnimeHandler.DeleteAnime)
	}

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
