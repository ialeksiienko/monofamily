package bot

import (
	"time"

	tb "gopkg.in/telebot.v3"
)

type TelegramBot struct {
	Bot *tb.Bot
}

type TBConfig struct {
	BotToken   string
	LongPoller int
}

func New(cfg TBConfig) (*TelegramBot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  cfg.BotToken,
		Poller: &tb.LongPoller{Timeout: time.Duration(cfg.LongPoller) * time.Second},
	})
	if err != nil {
		return nil, err
	}

	tgBot := &TelegramBot{
		Bot: b,
	}

	return tgBot, nil
}

func (tgbot TelegramBot) Start() {
	tgbot.Bot.Start()
}
