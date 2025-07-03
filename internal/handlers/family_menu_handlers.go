package handlers

import (
	"errors"
	"fmt"
	"main-service/internal/sessions"
	"main-service/internal/usecases"
	"strconv"

	tb "gopkg.in/telebot.v3"
)

func (h *Handler) GetMembers(c tb.Context) error {
	userID := c.Sender().ID

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		h.bot.Send(c.Sender(), "Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	members, err := h.usecases.UserService.GetMembersInfo(us.Family, userID)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeFamilyHasNoMembers {
				return c.Send("У вашій сім'ї поки немає учасників.")
			}
		}
		return c.Send("Не вдалося отримати інформацію про учасників сім'ї.")
	}

	c.Send("📋 Список учасників сім'ї:\n")

	for _, member := range members {
		role := "Учасник"
		if member.IsAdmin {
			role = "Адміністратор"
		}

		userLabel := ""
		if member.IsCurrent {
			userLabel = " (це ви)"
		}

		text := fmt.Sprintf(
			"👤 %s @%s %s\n- Роль: %s\n- ID: %d",
			member.Firstname,
			member.Username,
			userLabel,
			role,
			member.ID,
		)

		isAdmin := userID == us.Family.CreatedBy

		if !member.IsCurrent && isAdmin {
			btn := tb.InlineButton{
				Unique: "delete_member",
				Text:   "🗑 Видалити",
				Data:   strconv.FormatInt(member.ID, 10),
			}

			markup := &tb.ReplyMarkup{}
			markup.InlineKeyboard = [][]tb.InlineButton{
				{btn},
			}

			c.Send(text, markup)
		} else {
			c.Send(text)
		}
	}

	return c.Send(fmt.Sprintf("Всього учасників: %d", len(members)))
}

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

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		h.bot.Send(c.Sender(), "Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	err := h.usecases.UserService.LeaveFamily(us.Family, userID)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeCannotRemoveSelf {
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

// admin handlers

func (h *Handler) DeleteMember(c tb.Context) error {
	data := c.Callback().Data

	memberID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID.")
	}

	member, err := h.usecases.UserService.GetUserByID(memberID)
	if err != nil {
		return c.Send(ErrInternalServerForUser.Error())
	}

	inlineKeys := [][]tb.InlineButton{
		{BtnMemberDeleteNo}, {tb.InlineButton{Unique: "delete_member_yes", Text: "✅ Так", Data: strconv.FormatInt(member.ID, 10)}},
	}

	return c.Send(
		fmt.Sprintf("Ви дійсно хочете видалити учасника `%s`?", member.Firstname),
		&tb.SendOptions{
			ParseMode:   tb.ModeMarkdown,
			ReplyMarkup: &tb.ReplyMarkup{InlineKeyboard: inlineKeys},
		},
	)
}

func (h *Handler) ProcessMemberDeletion(c tb.Context) error {
	userID := c.Sender().ID
	data := c.Callback().Data

	memberID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return c.Send("Некоректний ID.")
	}

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		h.bot.Send(c.Sender(), "Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	removeErr := h.usecases.AdminService.RemoveMember(us.Family, userID, memberID)
	if removeErr != nil {
		switch e := err.(type) {
		case *usecases.CustomError[struct{}]:
			if e.Code == usecases.ErrCodeNoPermission {
				return c.Send("У вас немає прав на видалення.")
			}
			if e.Code == usecases.ErrCodeCannotRemoveSelf {
				return c.Send("Ви не можете видалити себе.")
			}
		}
		return c.Send("Не вдалося видалити користувача з сім'ї. Спробуйте ще раз пізніше.")
	}

	h.bot.Edit(c.Message(), "Учасника успішно видалено. Оновлюю список...")

	h.bot.Send(c.Sender(), "── 🔹 Оновлення списку 🔹 ──")

	return h.GetMembers(c)
}

func (h *Handler) CancelMemberDeletion(c tb.Context) error {
	h.bot.Delete(c.Message())

	return c.Send("Скасовано. Учасника не було видалено.")
}

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

	h.bot.Delete(c.Message())

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		h.bot.Send(c.Sender(), "Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	err := h.usecases.AdminService.DeleteFamily(us.Family, userID)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeNoPermission {
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

func (h *Handler) CreateNewInviteCode(c tb.Context) error {
	userID := c.Sender().ID

	us, exists := sessions.GetUserState(userID)
	if !exists || us.Family == nil {
		h.bot.Send(c.Sender(), "Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
		return h.GoHome(c)
	}

	code, expiresAt, err := h.usecases.AdminService.CreateNewFamilyCode(us.Family, userID)
	if err != nil {
		var custErr *usecases.CustomError[struct{}]
		if errors.As(err, &custErr) {
			if custErr.Code == usecases.ErrCodeNoPermission {
				return c.Send("У вас немає прав на створення нового коду запрошення.")
			}
			if custErr.Code == usecases.ErrCodeFailedToGenerateInviteCode {
				return c.Send("Не вдалося створити новий код запрошення. Спробуйте пізніше.")
			}
		}
		return c.Send("Не вдалося створити код запрошення. Спробуйте ще раз пізніше.")
	}

	return c.Send(fmt.Sprintf("Новий код запрошення: `%s`\n\nДійсний до — %s (час за Гринвічем, GMT)", code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
}
