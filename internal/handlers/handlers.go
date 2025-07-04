package handlers

import (
	"errors"
	"log/slog"
	"main-service/internal/entities"
	"main-service/internal/sessions"
	"main-service/internal/sl"
	"main-service/internal/usecases"
	"strings"

	tb "gopkg.in/telebot.v3"
)

var (
	ErrInternalServerForUser = errors.New("Сталася помилка на боці серверу, спробуйте пізніше.")
)

type Handler struct {
	bot      *tb.Bot
	sl       *sl.MyLogger
	usecases *usecases.Services
}

func New(bot *tb.Bot, sl *sl.MyLogger, usecases *usecases.Services) *Handler {
	return &Handler{
		bot:      bot,
		sl:       sl,
		usecases: usecases,
	}
}

func (h *Handler) Start(c tb.Context) error {
	userID := c.Sender().ID

	_, err := h.usecases.UserService.Register(&entities.User{
		ID:        userID,
		Username:  c.Sender().Username,
		Firstname: c.Sender().FirstName,
	})
	if err != nil {
		h.sl.Error("failed to save user", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
		return c.Send("Сталася помилка при зберіганні данних користувача. Спробуй пізніше.")
	}

	if _, exists := sessions.GetUserState(userID); !exists {
		msg, _ := h.bot.Send(c.Sender(), ".", &tb.SendOptions{
			ReplyMarkup: &tb.ReplyMarkup{
				RemoveKeyboard: true,
			},
		})

		h.bot.Delete(msg)
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnCreateFamily}, {BtnJoinFamily}, {BtnEnterMyFamily},
	}

	return c.Send("Привіт! Цей бот допоможе дізнатися рахунок на карті Monobank.\n\nВибери один з варіантів на клавіатурі.", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h *Handler) HandleText(c tb.Context) error {
	userID := c.Sender().ID
	state := sessions.GetTextState(userID)

	if state == sessions.StateNone {
		return h.handleRegularText(c)
	}

	text := c.Text()

	switch state {
	case sessions.StateWaitingFamilyName:
		sessions.ClearTextState(userID)
		return h.processFamilyCreation(c, text)

	case sessions.StateWaitingFamilyCode:
		sessions.ClearTextState(userID)
		return h.processFamilyJoin(c, strings.ToUpper(text))

	default:
		sessions.ClearTextState(userID)
		return h.handleRegularText(c)
	}
}

func (h *Handler) handleRegularText(c tb.Context) error {
	return c.Send("Будь ласка, скористайтеся кнопками для взаємодії з ботом.")
}

func (h *Handler) GoHome(c tb.Context) error {
	userID := c.Sender().ID

	sessions.DeleteUserState(userID)

	{
		msg, _ := h.bot.Send(c.Sender(), ".", &tb.SendOptions{
			ReplyMarkup: &tb.ReplyMarkup{
				RemoveKeyboard: true,
			},
		})

		h.bot.Delete(msg)
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnCreateFamily}, {BtnJoinFamily}, {BtnEnterMyFamily},
	}

	c.Delete()

	return c.Send("Виберіть один з варіантів на клавіатурі.", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}
