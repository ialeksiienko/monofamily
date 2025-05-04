package routes

import (
	"fmt"
	tb "gopkg.in/telebot.v3"
	"monofinances/internal/models"
	"monofinances/internal/repository"
)

type GetterUserFamilies interface {
	GetFamilyByUserID(userID string) ([]models.Family, error)
}

func SetupRoutes(bot *tb.Bot, db *repository.Database) {

	bot.Handle("/start", func(c tb.Context) error {
		families, err := db.GetFamilyByUserID(c.Sender().ID)
		if err != nil {
			return err
		}

		if len(families) == 0 {
			return c.Send("Привіт! У тебе поки немає жодної сім'ї. Додай або приєднайся.")
		}

		var familyList string
		for i, f := range families {
			familyList += fmt.Sprintf("%d. %s\n", i+1, f.Name)
		}

		msg := fmt.Sprintf("Привіт! Цей бот допоможе дізнатися рахунок на карті Monobank.\n\n"+
			"Твої сім'ї (%d):\n%s", len(families), familyList)

		bot.Send(c.Sender(), msg)

		return nil
	})

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
	//				User: &models.User{
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
	//				User: &models.User{
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
	//				User: &models.User{
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
	//				User: &models.User{
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
	//				User: &models.User{
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
	//				User: &models.User{
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
