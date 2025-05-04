package main

import (
	"log"
	"log/slog"
	"main-service/internal/adapters/database"
	"main-service/internal/bot"
	"main-service/internal/config"
	"main-service/internal/repository"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log0 := setupLogger(cfg.Env)

	dbPool, _, err := database.NewDBPool(database.DatabaseConfig{
		Username: "golang",
		Password: "golang",
		Hostname: "localhost",
		Port:     "5432",
		DBName:   "golangtest",
	})

	defer dbPool.Close()

	if err != nil {
		log.Fatalf("unexpected error while tried to connect to database: %v\n", err)
	}

	db := repository.New(dbPool, log0)

	tgBot := bot.NewTelegramBot(bot.TBConfig{
		Log:        log0,
		BotToken:   cfg.Bot.Token,
		LongPoller: cfg.Bot.LongPoller,
		DB:         db,
	})

	log0.Info("bot is running")

	tgBot.Start()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	default:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
