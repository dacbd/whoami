package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/httplog/v2"
	_ "github.com/joho/godotenv/autoload"
	"golang.org/x/term"
)

func shouldLogJSON() bool {
	if term.IsTerminal(int(os.Stdout.Fd())) {
		return false
	}
	return true
}

func getPort() int {
	if value, ok := os.LookupEnv("PORT"); ok {
		if port, err := strconv.Atoi(value); err != nil {
			return port
		}
	}
	return 3000
}

type Server struct {
	port       int
	httpLogger *httplog.Logger
}

func NewServer() *http.Server {
	httpLogger := httplog.NewLogger("whoami", httplog.Options{
		JSON:             shouldLogJSON(),
		LogLevel:         slog.LevelDebug,
		Concise:          true,
		RequestHeaders:   true,
		MessageFieldName: "message",
		/*
			Tags: map[string]string{
				"version": "v1.0-81aa4244d9fc8076a",
				"env":     "dev",
			},
		*/
		QuietDownRoutes: []string{
			"/health",
		},
		QuietDownPeriod: 10 * time.Second,
		// SourceFieldName: "source",
	})
	NewServer := &Server{
		port:       getPort(),
		httpLogger: httpLogger,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	slog.Info("Server Initialized", "Addr", server.Addr)

	return server
}
