package app

import (
	"context"
	"monofamily/internal/adapter/database"
	"monofamily/internal/adapter/database/familyrepo"
	"monofamily/internal/adapter/database/userrepo"
	"monofamily/internal/delivery/telegram"
	"monofamily/internal/delivery/telegram/handler"
	"monofamily/internal/pkg/sl"
	"monofamily/internal/service/familyservice"
	"monofamily/internal/service/userservice"
	"monofamily/internal/usecase"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	tb "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	bot *tb.Bot

	pgsqlxpool *pgxpool.Pool

	sl sl.Logger
}

type TBConfig struct {
	BotToken   string
	LongPoller int

	Pgsqlxpool *pgxpool.Pool

	Logger sl.Logger
}

func NewBot(cfg TBConfig) (*TelegramBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: time.Duration(cfg.LongPoller) * time.Second},
	})
	if err != nil {
		return nil, err
	}

	tgBot := &TelegramBot{
		bot:        b,
		sl:         cfg.Logger,
		pgsqlxpool: cfg.Pgsqlxpool,
	}

	return tgBot, nil
}

func (tgbot *TelegramBot) RunBot() {
	logger := tgbot.sl

	db := database.New(tgbot.pgsqlxpool)
	familyrepo := familyrepo.New(db.DB, logger)
	userrepo := userrepo.New(db.DB, logger)

	familyservice := familyservice.New(familyrepo, logger)
	userservice := userservice.New(userrepo, logger)

	usecase := usecase.New(userservice, userservice, familyservice)

	handler := handler.New(usecase, tgbot.bot, logger)

	go func() {
		for {
			err := familyservice.ClearInviteCodes(context.Background())
			if err != nil {
				logger.Error(err.Error())
			} else {
				logger.Debug("invite codes cleared successfully")
			}
			time.Sleep(24 * time.Hour)
		}
	}()

	telegram.SetupRoutes(tgbot.bot, handler)

	tgbot.bot.Start()
}
