package validate

import "regexp"

var tokenRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{44}$`)

func IsValidBankToken(token string) bool {
	return tokenRegex.MatchString(token)
}
