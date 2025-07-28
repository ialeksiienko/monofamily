package handler

import (
	"context"
	"errors"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) LeaveFamily(c tb.Context) error {
	inlineKeys := [][]tb.InlineButton{
		{BtnLeaveFamilyNo}, {BtnLeaveFamilyYes},
	}

	return c.Send("Ви дійсно хочете вийти з сім'ї?", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h *Handler) ProcessLeaveFamily(c tb.Context) error {
	userID := c.Sender().ID
	ctx := context.Background()

	us, exists := session.GetUserState(userID)
	if !exists || us.Family == nil {
		h.bot.Send(c.Sender(), "Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	err := h.usecase.LeaveFamily(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeCannotRemoveSelf {
				return c.Send("Адміністратор не може вийти з сім'ї.")
			}
		}
		return c.Send("Не вдалося вийти з сім'ї. Спробуйте ще раз пізніше.")
	}

	h.bot.Send(c.Sender(), "Ви успішно вийшли з сім'ї.")

	return h.GoHome(c)
}

func (h *Handler) CancelLeaveFamily(c tb.Context) error {
	h.bot.Delete(c.Message())

	return c.Send("Скасовано. Ви не вийшли з сім'ї.")
}
