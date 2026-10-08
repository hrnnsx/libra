package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"github.com/hrnnsx/libra/internal/database"
	"github.com/hrnnsx/libra/internal/handler"
	"github.com/hrnnsx/libra/internal/repository"
	"github.com/hrnnsx/libra/internal/service"
)

type Server struct {
	port int
	db   database.Service

	// handler
	authHandler *handler.AuthHandler
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	db := database.New()

	// repository
	userRepository := repository.NewUserRepository(db.GORM())

	// service
	authService := service.NewAuthService(userRepository)

	// handler
	authHandler := handler.NewAuthHandler(authService)

	newServer := &Server{
		port:        port,
		db:          db,
		authHandler: authHandler,
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
