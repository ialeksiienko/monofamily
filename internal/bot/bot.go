package bot

import (
	"log/slog"
	"monofinances/internal/repository"
	"monofinances/internal/routes"
	"time"

	tb "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	bot *tb.Bot
	log *slog.Logger
}

type TBConfig struct {
	Log        *slog.Logger
	BotToken   string
	LongPoller int
	UserToken  string
	DB         *repository.Database
}

func NewTelegramBot(cfg TBConfig) *TelegramBot {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: time.Duration(cfg.LongPoller) * time.Second},
	})
	if err != nil {
		panic("failed to build tg bot: " + err.Error())
	}

	tgBot := &TelegramBot{
		bot: b,
		log: cfg.Log,
	}

	routes.SetupRoutes(b, cfg.DB)

	return tgBot
}

func (tgbot TelegramBot) Start() {
	tgbot.bot.Start()
}
