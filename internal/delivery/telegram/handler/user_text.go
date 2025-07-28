package handler

import (
	"log/slog"
	"monofamily/internal/session"
	"strings"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) HandleText(c tb.Context) error {
	userID := c.Sender().ID
	state := session.GetTextState(userID)

	if state == session.StateNone {
		h.sl.Warn("unexpected state in HandleText", slog.Int64("user_id", userID), slog.String("state", string(state)))
		return h.handleRegularText(c)
	}

	text := strings.TrimSpace(c.Text())

	session.ClearTextState(userID)

	switch state {
	case session.StateWaitingFamilyName:
		return h.processFamilyCreation(c, text)

	case session.StateWaitingFamilyCode:
		return h.processFamilyJoin(c, strings.ToUpper(text))

	default:
		return h.handleRegularText(c)
	}
}

func (h *Handler) handleRegularText(c tb.Context) error {
	return c.Send("Будь ласка, скористайтеся кнопками для взаємодії з ботом.")
}
