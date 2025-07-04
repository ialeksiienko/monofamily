package handlers

import (
	"errors"
	"fmt"
	"main-service/internal/sessions"
	"main-service/internal/usecases"
	"time"
	"unicode/utf8"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) CreateFamily(c tb.Context) error {
	h.bot.Send(c.Sender(), "Введи назву нової сім'ї:")

	sessions.SetTextState(c.Sender().ID, sessions.StateWaitingFamilyName)

	return nil
}

func (h *Handler) processFamilyCreation(c tb.Context, familyName string) error {
	userID := c.Sender().ID

	if utf8.RuneCountInString(familyName) > 20 {
		return c.Send("Назва сім'ї не має містити більше 20 символів.")
	}

	code, expiresAt, err := h.usecases.FamilyService.Create(familyName, userID)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeFailedToGenerateInviteCode {
				return c.Send("Не вдалося створити новий код запрошення. Спробуйте пізніше.")
			}
		}
		return c.Send(ErrInternalServerForUser.Error)
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnEnterMyFamily},
	}

	return c.Send(fmt.Sprintf("Сім'я створена. Код запрошення:\n\n`%s`\n\nДійсний до — %s (час за Гринвічем, GMT)", code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	}, &tb.ReplyMarkup{InlineKeyboard: inlineKeys})
}

func (h *Handler) JoinFamily(c tb.Context) error {
	h.bot.Send(c.Sender(), "Введи код запрошення:")

	sessions.SetTextState(c.Sender().ID, sessions.StateWaitingFamilyCode)

	return nil
}

func (h *Handler) processFamilyJoin(c tb.Context, code string) error {
	userID := c.Sender().ID

	if utf8.RuneCountInString(code) != 6 {
		return c.Send("Код запрошення має містити 6 символів.")
	}

	familyName, err := h.usecases.FamilyService.Join(code, userID)
	if err != nil {
		switch e := err.(type) {
		case *usecases.CustomError[time.Time]:
			if e.Code == usecases.ErrCodeFamilyCodeExpired {
				return c.Send(fmt.Sprintf("Код запрошення не дійсний, закінчився - %s (час за Гринвічем, GMT)", e.Data.Format("02.01.2006 о 15:04")))
			}
		case *usecases.CustomError[struct{}]:
			if e.Code == usecases.ErrCodeFamilyNotFound {
				return c.Send("Сім'ю з цим кодом запрошення не знайдено.")
			}
		}
		return c.Send(ErrInternalServerForUser.Error)
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnEnterMyFamily},
	}

	return c.Send(fmt.Sprintf("Ви успішно приєдналися до сім'ї! Назва - %s", familyName), &tb.ReplyMarkup{
		InlineKeyboard: inlineKeys,
	})
}

func (h *Handler) EnterMyFamily(c tb.Context) error {
	userID := c.Sender().ID

	families, err := h.usecases.FamilyService.GetFamilies(userID)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeUserHasNoFamily {
				inlineKeys := [][]tb.InlineButton{
					{BtnCreateFamily}, {BtnJoinFamily},
				}

				return c.Send("Привіт! У вас поки немає жодної сім'ї. Створіть або приєднайтеся.", &tb.ReplyMarkup{
					InlineKeyboard: inlineKeys,
				})
			}
		}
		return c.Send(ErrInternalServerForUser.Error)
	}

	sessions.SetUserPageState(userID, &sessions.UserPageState{
		Page:     0,
		Families: families,
	})

	return showFamilyListPage(c, families, 0)
}
