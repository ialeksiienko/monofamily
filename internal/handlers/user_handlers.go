package handlers

// import (
// 	"errors"
// 	"fmt"
// 	"log/slog"
// 	"main-service/internal/sessions"
// 	"main-service/internal/usecases"
// 	"strconv"

// 	tb "gopkg.in/telebot.v3"
// )

// type Validator interface {
// 	BankToken(token string) (bool, error)
// }

// func (h *Handler) GetMembers(c tb.Context) error {
// 	userID := c.Sender().ID

// 	us, ok := c.Get("us").(sessions.UserState)
// 	if !ok {
// 		h.sl.Error("unable to get user state", slog.Int("user_id", int(userID)))
// 		return c.Send(ErrInternalServerForUser.Error())
// 	}

// 	members, err := h.usecases.UserService.GetMembersInfo(us.Family, userID)
// 	if err != nil {
// 		var custErr *usecases.CustomError[struct{}]
// 		if errors.As(err, &custErr) {
// 			if custErr.Code == usecases.ErrCodeFamilyHasNoMembers {
// 				return c.Send("У вашій сім'ї поки немає учасників.")
// 			}
// 		}
// 		return c.Send("Не вдалося отримати інформацію про учасників сім'ї.")
// 	}

// 	c.Send("📋 Список учасників сім'ї:\n")

// 	for _, member := range members {
// 		role := "Учасник"
// 		if member.IsAdmin {
// 			role = "Адміністратор"
// 		}

// 		userLabel := ""
// 		if member.IsCurrent {
// 			userLabel = " (це ви)"
// 		}

// 		text := fmt.Sprintf(
// 			"👤 %s @%s %s\n- Роль: %s\n- ID: %d",
// 			member.Firstname,
// 			member.Username,
// 			userLabel,
// 			role,
// 			member.ID,
// 		)

// 		isAdmin := userID == us.Family.CreatedBy

// 		if !member.IsCurrent && isAdmin {
// 			btn := tb.InlineButton{
// 				Unique: "delete_member",
// 				Text:   "🗑 Видалити",
// 				Data:   strconv.FormatInt(member.ID, 10),
// 			}

// 			markup := &tb.ReplyMarkup{}
// 			markup.InlineKeyboard = [][]tb.InlineButton{
// 				{btn},
// 			}

// 			c.Send(text, markup)
// 		} else {
// 			c.Send(text)
// 		}
// 	}

// 	return c.Send(fmt.Sprintf("Всього учасників: %d", len(members)))
// }

// func (h *Handler) LeaveFamily(c tb.Context) error {
// 	inlineKeys := [][]tb.InlineButton{
// 		{BtnLeaveFamilyNo}, {BtnLeaveFamilyYes},
// 	}

// 	return c.Send("Ви дійсно хочете вийти з сім'ї?", &tb.ReplyMarkup{
// 		InlineKeyboard: inlineKeys,
// 	})
// }

// func (h *Handler) ProcessLeaveFamily(c tb.Context) error {
// 	userID := c.Sender().ID

// 	us, ok := c.Get("us").(sessions.UserState)
// 	if !ok {
// 		h.sl.Error("unable to get user state", slog.Int("user_id", int(userID)))
// 		return c.Send(ErrInternalServerForUser.Error())
// 	}

// 	err := h.usecases.UserService.LeaveFamily(us.Family, userID)
// 	if err != nil {
// 		var custErr *usecases.CustomError[struct{}]
// 		if errors.As(err, &custErr) {
// 			if custErr.Code == usecases.ErrCodeCannotRemoveSelf {
// 				return c.Send("Адміністратор не може вийти з сім'ї.")
// 			}
// 		}
// 		return c.Send("Не вдалося вийти з сім'ї. Спробуйте ще раз пізніше.")
// 	}

// 	h.bot.Send(c.Sender(), "Ви успішно вийшли з сім'ї.")

// 	return h.GoHome(c)
// }

// func (h *Handler) CancelLeaveFamily(c tb.Context) error {
// 	h.bot.Delete(c.Message())

// 	return c.Send("Скасовано. Ви не вийшли з сім'ї.")
// }

// func (h *Handler) SaveUserBankToken(c tb.Context) error {
// 	button := tb.InlineButton{
// 		Unique: "mono_link",
// 		Text:   "Посилання",
// 		URL:    "https://api.monobank.ua/",
// 	}

// 	inlineKeys := [][]tb.InlineButton{
// 		{button},
// 	}

// 	h.bot.Send(c.Sender(), "Перейдіть по посиланню знизу та відправте свій токен в цей чат.", &tb.ReplyMarkup{
// 		InlineKeyboard: inlineKeys,
// 	})

// 	sessions.SetTextState(c.Sender().ID, sessions.StateWaitingBankToken)

// 	return nil
// }

// func (h *Handler) processUserBankToken(c tb.Context, token string) error {
// 	userID := c.Sender().ID

// 	us, exists := sessions.GetUserState(userID)
// 	if !exists || us.Family == nil {
// 		c.Send("Ви не увійшли в сім'ю. Спочатку потрібно увійти в сім'ю.")
// 		return h.GoHome(c)
// 	}

// 	valid, err := h.validator.BankToken(token)
// 	if err != nil {
// 		return c.Send("Не вдалося перевірити токен. Спробуйте пізніше.")
// 	}

// 	if !valid {
// 		return c.Send("Неправильний формат токена.")
// 	}

// 	_, saveErr := h.usecases.UserBankTokenService.Save(us.Family.ID, userID, token)
// 	if saveErr != nil {
// 		return c.Send("Не вдалося зберегти токен. Спробуйте пізніше.")
// 	}

// 	isAdmin := us.Family.CreatedBy == userID

// 	rows := generateFamilyMenu(isAdmin, true)

// 	menu.Reply(rows...)

// 	return c.Send("Ви успішно зберегли токен для цієї сім'ї.", menu)
// }
