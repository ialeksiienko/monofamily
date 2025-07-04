package routes

import (
	"main-service/internal/handlers"

	tb "gopkg.in/telebot.v3"
)

func SetupRoutes(bot *tb.Bot, h *handlers.Handler) {

	bot.Handle("/start", h.Start)

	bot.Handle(tb.OnText, h.HandleText)

	// first buttons
	{
		bot.Handle(&handlers.BtnCreateFamily, h.CreateFamily)

		bot.Handle(&handlers.BtnJoinFamily, h.JoinFamily)

		bot.Handle(&handlers.BtnEnterMyFamily, h.EnterMyFamily)
	}

	// enter my family
	{
		bot.Handle(&tb.InlineButton{Unique: "select_family"}, h.SelectMyFamily)

		bot.Handle(&handlers.BtnNextPage, h.NextPage)

		bot.Handle(&handlers.BtnPrevPage, h.PrevPage)

		bot.Handle(&tb.InlineButton{Unique: "go_home"}, h.GoHome)
	}

	// family menu
	{
		bot.Handle(&handlers.MenuViewMembers, h.GetMembers)

		{
			bot.Handle(&handlers.MenuLeaveFamily, h.LeaveFamily)

			bot.Handle(&handlers.BtnLeaveFamilyNo, h.CancelLeaveFamily)
			bot.Handle(&handlers.BtnLeaveFamilyYes, h.ProcessLeaveFamily)
		}

		// admin menu
		{
			bot.Handle(&tb.InlineButton{Unique: "delete_member"}, h.DeleteMember)

			bot.Handle(&handlers.BtnMemberDeleteNo, h.CancelMemberDeletion)
			bot.Handle(&tb.InlineButton{Unique: "delete_member_yes"}, h.ProcessMemberDeletion)
		}

		{
			bot.Handle(&handlers.MenuDeleteFamily, h.DeleteFamily)

			bot.Handle(&handlers.BtnFamilyDeleteNo, h.CancelFamilyDeletion)
			bot.Handle(&handlers.BtnFamilyDeleteYes, h.ProcessFamilyDeletion)
		}

		bot.Handle(&handlers.MenuCreateNewCode, h.CreateNewInviteCode)

		bot.Handle(&handlers.MenuGoHome, h.GoHome)
	}

	//bot.Handle("/family", func(c tb.Context) error {
	//	userID := c.Sender().ID
	//
	//	families, err := repo.GetFamiliesByUserID(userID)
	//	if err != nil {
	//		sl.Error("failed to get family", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
	//		return c.Send("Не вдалося отримати дані про сім'ю.")
	//	}
	//
	//	if len(families) == 0 {
	//		return c.Send("Ти ще не приєднаний до жодної сім’ї.")
	//	}
	//
	//	msg := "Твої сім’ї:\n"
	//	for _, f := range families {
	//		users, err := repo.GetAllUsersInFamily(&f)
	//		if err != nil {
	//			sl.Error("failed to get users", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
	//			return c.Send("Не вдалося отримати всіх користувачів сім'ї. Спробуй пізніше.")
	//		}
	//
	//		usersList := make([]string, len(users))
	//		for _, u := range users {
	//			usersList = append(usersList, u.Username)
	//		}
	//
	//		msg += fmt.Sprintf("🔸 %s (ID: %d). \nУчасники: %s\n", f.Name, f.ID, strings.Join(usersList, ", "))
	//	}
	//
	//	return c.Send(msg)
	//})
	//
	//button := tb.InlineButton{
	//			Unique: "mono_link",
	//			Text:   "Силка",
	//			URL:    "https://api.monobank.ua/",
	//		}
	//
	//		inlineKeys := [][]tb.InlineButton{
	//			{button},
	//		}
	//
	//		bot.Send(c.Sender(), "Привіт, цей бот допоможе дізнатися рахунок на карті monobank.\n\nПерейди по силці внизу та відправ свій токен в цей чат.", &tb.ReplyMarkup{
	//			InlineKeyboard: inlineKeys,
	//		})
	//
	//		bot.Handle(tb.OnText, func(c tb.Context) error {
	//			return c.Send("Данные успешно сохранены в базе данных!")
	//		}, middlewares.CheckTokenValid)
	//		return nil

	//h.bot.Handle("/balance", func(c tb.Context) error {
	//	buttonBlack := tb.InlineButton{Unique: "black", Text: "Черная"}
	//	buttonWhite := tb.InlineButton{Unique: "white", Text: "Белая"}
	//
	//	inlineKeysCardType := [][]tb.InlineButton{
	//		{buttonBlack},
	//		{buttonWhite},
	//	}
	//
	//	h.bot.Send(c.Sender(), "Напиши какого типа карточки ты хотел бы узнать баланс.", &tb.ReplyMarkup{InlineKeyboard: inlineKeysCardType})
	//
	//	h.bot.Handle(&buttonBlack, func(c tb.Context) error {
	//		buttonHryvnia := tb.InlineButton{Unique: "hryvnia", Text: "Гривны"}
	//		buttonZloty := tb.InlineButton{Unique: "zloty", Text: "Злотые"}
	//		buttonDollars := tb.InlineButton{Unique: "dollars", Text: "Доллары"}
	//
	//		inlineKeysCurrency := [][]tb.InlineButton{
	//			{buttonHryvnia},
	//			{buttonZloty},
	//			{buttonDollars},
	//		}
	//
	//		h.bot.Send(c.Sender(), "Теперь нажми на кнопку в какой валюте ты хочешь узнать баланс.", &tb.ReplyMarkup{InlineKeyboard: inlineKeysCurrency})
	//
	//		h.bot.Handle(&buttonZloty, func(c tb.Context) error {
	//			rabbitDataSelect := &messaging.RabbitMQ{
	//				Operation: "select",
	//				User: &entities.User{
	//					ID: c.Sender().ID,
	//				},
	//			}
	//
	//			userFromBank, err := h.rabbitMQConn.SetupAndConsume(messaging.BankService, rabbitDataSelect)
	//			if err != nil || userFromBank == nil {
	//				h.log.Error("failed to get user from db", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя из базы данных.")
	//			}
	//
	//			rabbitDataUser := &messaging.RabbitMQ{
	//				Operation: "user",
	//				User: &entities.User{
	//					Token:    userFromBank.Token,
	//					CardType: buttonBlack.Unique,
	//					Currency: buttonZloty.Unique,
	//				},
	//			}
	//
	//			userFromApi, err := h.rabbitMQConn.SetupAndConsume(messaging.ApiService, rabbitDataUser)
	//			if err != nil || userFromApi == nil {
	//				h.log.Error("failed to get user data from api service", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя.")
	//			}
	//
	//			return c.Send("Баланс: " + userFromApi.Balance + " Тип карточки: " + userFromApi.CardType + " Валюта: " + userFromApi.Currency)
	//		})
	//		h.bot.Handle(&buttonDollars, func(c tb.Context) error {
	//			rabbitDataSelect := &messaging.RabbitMQ{
	//				Operation: "select",
	//				User: &entities.User{
	//					ID: c.Sender().ID,
	//				},
	//			}
	//
	//			userFromBank, err := h.rabbitMQConn.SetupAndConsume(messaging.BankService, rabbitDataSelect)
	//			if err != nil || userFromBank == nil {
	//				h.log.Error("failed to get user from db", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя из базы данных.")
	//			}
	//
	//			rabbitDataUser := &messaging.RabbitMQ{
	//				Operation: "user",
	//				User: &entities.User{
	//					Token:    userFromBank.Token,
	//					CardType: buttonBlack.Unique,
	//					Currency: buttonDollars.Unique,
	//				},
	//			}
	//
	//			userFromApi, err := h.rabbitMQConn.SetupAndConsume(messaging.ApiService, rabbitDataUser)
	//			if err != nil || userFromApi == nil {
	//				h.log.Error("failed to get user data from api service", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя.")
	//			}
	//
	//			return c.Send("Баланс: " + userFromApi.Balance + " Тип карточки: " + userFromApi.CardType + " Валюта: " + userFromApi.Currency)
	//		})
	//		h.bot.Handle(&buttonHryvnia, func(c tb.Context) error {
	//			rabbitDataSelect := &messaging.RabbitMQ{
	//				Operation: "select",
	//				User: &entities.User{
	//					ID: c.Sender().ID,
	//				},
	//			}
	//
	//			userFromBank, err := h.rabbitMQConn.SetupAndConsume(messaging.BankService, rabbitDataSelect)
	//			if err != nil || userFromBank == nil {
	//				h.log.Error("failed to get user from db", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя из базы данных.")
	//			}
	//
	//			rabbitDataUser := &messaging.RabbitMQ{
	//				Operation: "user",
	//				User: &entities.User{
	//					Token:    userFromBank.Token,
	//					CardType: buttonBlack.Unique,
	//					Currency: buttonHryvnia.Unique,
	//				},
	//			}
	//
	//			userFromApi, err := h.rabbitMQConn.SetupAndConsume(messaging.ApiService, rabbitDataUser)
	//			if err != nil || userFromApi == nil {
	//				h.log.Error("failed to get user data from api service", slog.String("error", err.Error()))
	//				return c.Send("Не удалось получить данные пользователя.")
	//			}
	//
	//			return c.Send("Баланс: " + userFromApi.Balance + " Тип карточки: " + userFromApi.CardType + " Валюта: " + userFromApi.Currency)
	//		})
	//		return nil
	//	})
	//	return nil
	//})
}
