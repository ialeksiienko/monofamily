package handler

import (
	"context"
	"errors"
	"fmt"
	"monofamily/internal/errorsx"
	"monofamily/internal/session"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) CreateNewInviteCode(c tb.Context) error {
	userID := c.Sender().ID
	ctx := context.Background()

	us, exists := session.GetUserState(userID)
	if !exists || us.Family == nil {
		h.bot.Send(c.Sender(), "Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	code, expiresAt, err := h.usecase.CreateNewInviteCode(ctx, us.Family, userID)
	if err != nil {
		var custErr *errorsx.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == errorsx.ErrCodeNoPermission {
				return c.Send("У вас немає прав на створення нового коду запрошення.")
			}
			if custErr.Code == errorsx.ErrCodeFailedToGenerateInviteCode {
				return c.Send("Не вдалося створити новий код запрошення. Спробуйте пізніше.")
			}
		}
		return c.Send("Не вдалося створити код запрошення. Спробуйте ще раз пізніше.")
	}

	return c.Send(fmt.Sprintf("Новий код запрошення: `%s`\n\nДійсний до — %s (час за Гринвічем, GMT)", code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
}
