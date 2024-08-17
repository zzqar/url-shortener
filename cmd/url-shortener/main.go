package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	// Init config: cleanenv /настройки сервера
	cfg := config.MustLoadConfig()

	// Init logger: slog /логировать ошибки
	log := setupLogger(cfg.Env)
	log.Info(fmt.Sprintln("Server started ENV:", cfg.Env))
	log.Debug("Debug message activated")

	// Init storage: sqlite / база
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Error initializing storage", sl.Err(err))
		os.Exit(1)
	}

	//  init router: chi, 'che render' /распределить запросов
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(logger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Get("/save_url", save.New(log, storage, cfg))

	log.Info(fmt.Sprintln("Starting server in address:", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IddleTimeout,
		Handler:      router,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Error("Error starting server", sl.Err(err))
	}
	log.Info("Server stopped")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	options := slogpretty.PrettyHandlerOptions{SlogOpts: &slog.HandlerOptions{Level: slog.LevelDebug}}
	handler := options.NewPrettyHandler(os.Stdout)
	return slog.New(handler)
}
