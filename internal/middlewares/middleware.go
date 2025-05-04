package middlewares

import (
	tb "gopkg.in/telebot.v3"
	"regexp"
)

func CheckTokenValid(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		message := c.Message().Text

		var tokenPattern = `^[a-zA-Z0-9_-]{44}$`
		matched, err := regexp.MatchString(tokenPattern, message)
		if err != nil {
			return c.Send("Ошибка при проверке токена.")
		}

		if !matched {
			return c.Send("Неверный формат токена.")
		}

		return next(c)
	}
}
