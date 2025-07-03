package middlewares

import (
	"regexp"

	tb "gopkg.in/telebot.v3"
)

func CheckTokenValid(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		message := c.Message().Text

		var tokenPattern = `^[a-zA-Z0-9_-]{44}$`
		matched, err := regexp.MatchString(tokenPattern, message)
		if err != nil {
			return c.Send("Не вдалося перевірити токен. Спробуйте пізніше.")
		}

		if !matched {
			return c.Send("Неправильний формат токена.")
		}

		return next(c)
	}
}
