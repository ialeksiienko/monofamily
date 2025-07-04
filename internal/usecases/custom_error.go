package usecases

type ErrorCode string

const (
	ErrCodeFamilyCodeExpired  ErrorCode = "family_code_expired"
	ErrCodeFamilyNotFound     ErrorCode = "family_not_found"
	ErrCodeUserHasNoFamily    ErrorCode = "user_has_no_family"
	ErrCodeFamilyHasNoMembers ErrorCode = "family_has_no_members"
	ErrCodeUserNotInFamily    ErrorCode = "user_not_in_family"

	ErrCodeNoPermission     ErrorCode = "no_permission"
	ErrCodeCannotRemoveSelf ErrorCode = "cannot_remove_self"

	ErrCodeFailedToGenerateInviteCode ErrorCode = "failed_to_generate_invite_code"
)

type CustomError[T any] struct {
	Data T
	Msg  string
	Code ErrorCode
}

func (e *CustomError[T]) Error() string {
	return e.Msg
}
