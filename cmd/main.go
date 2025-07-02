package main

import (
	"io"
	"log/slog"
	"main-service/internal/bot"
	"main-service/internal/config"
	"main-service/internal/database"
	"main-service/internal/handlers"
	"main-service/internal/repository"
	"main-service/internal/routes"
	"main-service/internal/sl"
	"main-service/internal/usecases"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	logger := setupLogger(cfg.Env)

	dbPool, _, err := database.NewDBPool(database.DatabaseConfig{
		Username: cfg.DB.User,
		Password: cfg.DB.Pass,
		Hostname: cfg.DB.Host,
		Port:     cfg.DB.Port,
		DBName:   cfg.DB.Name,

		Logger: logger,
	})
	if err != nil {
		logger.Fatal("unexpected error while tried to connect to database", slog.String("error", err.Error()))
	}

	defer dbPool.Close()

	tgBot, err := bot.New(bot.TBConfig{
		BotToken:   cfg.Bot.Token,
		LongPoller: cfg.Bot.LongPoller,
	})
	if err != nil {
		logger.Fatal("failed to build tg bot:", slog.String("error", err.Error()))
	}

	b := tgBot.Bot

	repo := 	repository.New(dbPool, logger)
	service := usecases.New(repo,repo, repo, logger)
	handler := handlers.New(b, logger, service)

	routes.SetupRoutes(b, handler)

	go func() {
		for {
			err := service.FamilyInviteCodeService.ClearInviteCodes()
			if err != nil {
				logger.Error("failed to clear invite codes", slog.String("error", err.Error()))
			} else {
				logger.Debug("invite codes cleared successfully")
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	logger.Info("bot is running")

	tgBot.Start()
}

func setupLogger(env string) *sl.MyLogger {
	var level zerolog.Level
	var writer io.Writer

	var slogLevel slog.Level

	switch strings.ToLower(env) {
	case envLocal:
		level = zerolog.DebugLevel
		slogLevel = slog.LevelDebug

		writer = zerolog.ConsoleWriter{Out: os.Stdout}
	case envDev:
		level = zerolog.DebugLevel
		slogLevel = slog.LevelDebug

		writer = os.Stdout
	case envProd:
		level = zerolog.InfoLevel
		slogLevel = slog.LevelInfo

		writer = os.Stdout
	default:
		level = zerolog.InfoLevel
		slogLevel = slog.LevelInfo

		writer = os.Stdout
	}

	zlogger := zerolog.New(writer).Level(level).With().Timestamp().Logger()

	handler := slogzerolog.Option{
		Level:  slogLevel,
		Logger: &zlogger,
	}.NewZerologHandler()

	return sl.New(slog.New(handler), func(msg string, attrs ...slog.Attr) {
		zlogger.Fatal().Fields(attrsToMap(attrs)).Msg(msg)
	})
}

func attrsToMap(attrs []slog.Attr) map[string]any {
	m := make(map[string]any)
	for _, attr := range attrs {
		m[attr.Key] = attr.Value.Any()
	}
	return m
}
