package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/hrnnsx/libra/external/anilist"
	"github.com/hrnnsx/libra/internal/database"
	"github.com/hrnnsx/libra/internal/handler"
	"github.com/hrnnsx/libra/internal/repository"
	"github.com/hrnnsx/libra/internal/service"
)

type Server struct {
	port int
	db   database.Service

	// handler
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	animeHandler   *handler.AnimeHandler
	libraryHandler *handler.LibraryHandler
	groupHandler   *handler.GroupHandler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db := database.New()

	// client
	anilistClient := anilist.NewClient()

	// repository
	userRepository := repository.NewUserRepository(db.GORM())

	animeRepository := repository.NewAnimeRepository(db.GORM())
	libraryAnimeRepository := repository.NewLibraryAnimeRepository(db.GORM())
	groupRepository := repository.NewGroupRepository(db.GORM())

	// service
	authService := service.NewAuthService(userRepository)
	userService := service.NewUserService(userRepository)

	animeService := service.NewAnimeService(anilistClient)
	libraryService := service.NewLibraryService(animeRepository, libraryAnimeRepository, anilistClient)
	groupService := service.NewGroupService(groupRepository)

	// handler
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)

	animeHandler := handler.NewAnimeHandler(animeService)
	libraryHandler := handler.NewLibraryHandler(libraryService)
	groupHandler := handler.NewGroupHandler(groupService)

	newServer := &Server{
		port: port,
		db:   db,

		authHandler:    authHandler,
		userHandler:    userHandler,
		animeHandler:   animeHandler,
		libraryHandler: libraryHandler,
		groupHandler:   groupHandler,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
