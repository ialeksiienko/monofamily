package handler

import (
	"context"
	"errors"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) DeleteFamily(c tb.Context) error {
	inlineKeys := [][]tb.InlineButton{
		{BtnFamilyDeleteNo}, {BtnFamilyDeleteYes},
	}

	return c.Send("Ви дійсно хочете видалити сім'ю?", &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h *Handler) ProcessFamilyDeletion(c tb.Context) error {
	userID := c.Sender().ID
	ctx := context.Background()

	h.bot.Delete(c.Message())

	us, ok := c.Get("user_state").(*session.UserState)
	if !ok || us == nil {
		return c.Send(ErrUnableToGetUserState.Error())
	}

	err := h.usecase.DeleteFamily(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeNoPermission {
				return c.Send("У вас немає прав на видалення.")
			}
		}
		return c.Send("Не вдалося видалити сім'ю. Спробуйте ще раз пізніше.")
	}

	h.bot.Send(c.Sender(), "Сім'ю успішно видалено.")

	return h.GoHome(c)
}

func (h *Handler) CancelFamilyDeletion(c tb.Context) error {
	h.bot.Delete(c.Message())

	return c.Send("Скасовано. Сім’ю не було видалено.")
}
