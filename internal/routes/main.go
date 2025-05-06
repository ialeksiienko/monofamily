package routes

import (
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	tb "gopkg.in/telebot.v3"
	"log/slog"
	"main-service/internal/models"
	"main-service/internal/repository"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

var (
	InternalServerErrorForUser = fmt.Errorf("Сталася помилка на боці серверу, спробуйте пізніше.")
)

var (
	btnCreateFamily = tb.InlineButton{
		Unique: "create_family_button",
		Text:   "Створити сім'ю",
	}
	btnJoinFamily = tb.InlineButton{
		Unique: "join_family_button",
		Text:   "Приєднатися",
	}
	btnRemoveFamily = tb.InlineButton{
		Unique: "remove_family_button",
		Text:   "Видалити сім'ю",
	}
	btnLeaveFamily = tb.InlineButton{
		Unique: "leave_family_button",
		Text:   "Вийти з сім'ї",
	}
	btnCreateNewCode = tb.InlineButton{
		Unique: "create_new_code_button",
		Text:   "Створити код запрошення",
	}
)

func SetupRoutes(bot *tb.Bot, db *repository.Database, sl *slog.Logger) {

	bot.Handle("/start", func(c tb.Context) error {
		userID := c.Sender().ID

		_, err := db.SaveUser(&models.User{
			Username:  c.Sender().Username,
			Firstname: c.Sender().FirstName,
		})
		if err != nil {
			sl.Error("failed to save user", slog.Int("userID", int(c.Sender().ID)), slog.String("err", err.Error()))
			return c.Send(InternalServerErrorForUser.Error())
		}

		families, err := db.GetFamiliesByUserID(userID)
		if err != nil {
			sl.Error("failed to get family by userID", slog.Int("userID", int(c.Sender().ID)), slog.String("err", err.Error()))
			return c.Send(InternalServerErrorForUser.Error())
		}

		inlineKeys := [][]tb.InlineButton{
			{btnCreateFamily}, {btnJoinFamily},
		}

		if len(families) == 0 {
			return c.Send("Привіт! У тебе поки немає жодної сім'ї. Створи або приєднайся.", &tb.ReplyMarkup{
				InlineKeyboard: inlineKeys,
			})
		}

		var familyList string
		for i, f := range families {
			familyList += fmt.Sprintf("%d. %s\n", i+1, f.Name)
		}

		msg := fmt.Sprintf("Привіт! Цей бот допоможе дізнатися рахунок на карті Monobank.\n\n"+
			"Твої сім'ї (%d):\n%s", len(families), familyList)

		bot.Send(c.Sender(), msg, &tb.ReplyMarkup{
			InlineKeyboard: inlineKeys,
		})

		return nil
	})

	bot.Handle(&btnCreateFamily, func(c tb.Context) error {
		userID := c.Sender().ID

		bot.Send(c.Sender(), "Введи назву нової сім'ї:")

		bot.Handle(tb.OnText, func(c tb.Context) error {
			familyName := c.Text()

			if utf8.RuneCountInString(familyName) > 20 {
				return c.Send("Назва сім'ї не має містити більше 20 символів.")
			}

			f, err := db.CreateFamily(&models.Family{
				Name:      familyName,
				CreatedBy: userID,
			})
			if err != nil {
				sl.Error("failed to create family", slog.Int("familyID", int(userID)), slog.String("err", err.Error()))
				return c.Send(InternalServerErrorForUser.Error())
			}

			sl.Debug("family created", slog.Int("familyID", int(userID)))

			saveErr := db.SaveUserToFamily(f)
			if saveErr != nil {
				sl.Error("unable to save user to family", slog.Int("userID", int(c.Sender().ID)), slog.String("err", err.Error()))
				return c.Send(InternalServerErrorForUser.Error())
			}

			code := generateInviteCode()

			expiresAt, err := db.SaveFamilyInviteCode(userID, f.ID, code)
			if err != nil {
				sl.Error("failed to save family invite code", slog.Int("familyID", int(c.Sender().ID)), slog.String("err", err.Error()))
				return c.Send(InternalServerErrorForUser.Error())
			}

			return c.Send(fmt.Sprintf("Сім'я створена. Код запрошення:\n\n`%s`\n\nДійсний до — %s", code, expiresAt.Format("02.01.2006 15:04")), &tb.SendOptions{
				ParseMode: tb.ModeMarkdown,
			})
		})
		return nil
	})

	bot.Handle(&btnJoinFamily, func(c tb.Context) error {
		bot.Send(c.Sender(), "Введи код запрошення.")

		bot.Handle(tb.OnText, func(c tb.Context) error {
			code := strings.ToUpper(c.Text())

			if len(code) != 6 {
				return c.Send("Код запрошення має містити 6 символів.")
			}

			f, expiresAt, err := db.GetFamilyByCode(code)
			if err != nil {
				sl.Error("failed to get family by code", slog.String("err", err.Error()))
				if errors.Is(err, pgx.ErrNoRows) {
					sl.Error("family not found with code", slog.String("code", code))
					return c.Send("Сім'ю з цим кодом запрошення не знайдено.")
				}
				return c.Send(InternalServerErrorForUser.Error())
			}

			if time.Now().After(expiresAt) {
				sl.Error("expired family by code", slog.String("err", err.Error()))
				return c.Send(fmt.Sprintf("Код запрошення не дійсний, закінчився - %s", expiresAt.Format("02.01.2006 о 15:04")))
			}

			saveErr := db.SaveUserToFamily(f)
			if saveErr != nil {
				sl.Error("unable to save user to family", slog.Int("userID", int(c.Sender().ID)), slog.String("err", err.Error()))
				return c.Send(InternalServerErrorForUser.Error())
			}

			return c.Send(fmt.Sprintf("Ви успішно приєдналися до сім'ї! Назва - %s", f.Name))
		})
		return nil
	})

	bot.Handle("/family", func(c tb.Context) error {
		userID := c.Sender().ID

		families, err := db.GetFamiliesByUserID(userID)
		if err != nil {
			sl.Error("failed to get family", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
			return c.Send("Не вдалося отримати дані про сім'ю.")
		}

		if len(families) == 0 {
			return c.Send("Ти ще не приєднаний до жодної сім’ї.")
		}

		msg := "Твої сім’ї:\n"
		for _, f := range families {
			users, err := db.GetAllUsersInFamily(&f)
			if err != nil {
				sl.Error("failed to get users", slog.Int("userID", int(userID)), slog.String("err", err.Error()))
				return c.Send("Не вдалося отримати всіх користувачів сім'ї. Спробуй пізніше.")
			}

			usersList := make([]string, len(users))
			for _, u := range users {
				usersList = append(usersList, u.Username)
			}

			msg += fmt.Sprintf("🔸 %s (ID: %d). \nУчасники: %s\n", f.Name, f.ID, strings.Join(usersList, ", "))
		}

		return c.Send(msg)
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

const codeLength = 6

var generateInviteCode = func() string {
	const letters = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b := make([]byte, codeLength)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
