package main

import (
	"fmt"
	"log/slog"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envProd  = "prod"
	envDev   = "dev"
)

func main() {
	//TODO Init config: cleanenv
	/*
		Что в себе он содержит ? // настройки сервера...
		Для чего ? // для чтения параметров настроек из файла
	*/
	cfg := config.MustLoadConfig()

	//TODO Init logger: slog
	/*
		логировать ошибки ? // да, ошибки и дебаг сообщения
		Почему тут ? // после запуска конфига мы можем узнать какой нам нужен логгер что бы потом отлавливать все ошибки
	*/
	log := setupLogger(cfg.Env)

	log.Info("Server started", slog.String("ENV:", cfg.Env))
	log.Debug("Debug message activated")

	/*
		TODO init storage: sqlite
		что он делает ?
		почему тут ?
	*/
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("Error initializing storage", sl.Err(err))
		os.Exit(1)
	}
	fmt.Println(storage)

	/*
	   TODO init controller: handler
	   что он делает?
	   почему тут?
	*/

	/*
	   TODO init router: chi, 'che render'
	   что делает?
	   почему тут?
	*/

	/*
		TODO run server
		что значит запускаем ?

	*/

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	default:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
