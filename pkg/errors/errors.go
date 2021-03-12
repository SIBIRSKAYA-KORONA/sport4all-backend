package errors

import (
	"errors"
	"net/http"
)

const (
	Internal     = "internal error"
	Conflict     = "conflict with exists data"
	NoPermission = "no permission for current operation"

	UserNotFound  = "user not exist"
	WrongPassword = "wrong password"

	SessionNotFound = "cookie invalid, session not exist"

	TeamNotFound = "team not exist"
)

var (
	// общие
	ErrInternal     = errors.New(Internal)
	ErrConflict     = errors.New(Conflict)
	ErrNoPermission = errors.New(NoPermission)

	// ошибки, связанные с юзером
	ErrUserNotFound  = errors.New(UserNotFound)
	ErrWrongPassword = errors.New(WrongPassword)

	// ошибки, связанные с сессией
	ErrSessionNotFound = errors.New(SessionNotFound)

	// ошибки, связанные с командой
	ErrTeamNotFound = errors.New(TeamNotFound)
)

var errorToCodeMap = map[error]int{
	// общие
	ErrInternal:     http.StatusInternalServerError,
	ErrConflict:     http.StatusConflict,
	ErrNoPermission: http.StatusForbidden,

	// ошибки, связанные с юзером
	ErrUserNotFound:  http.StatusNotFound,
	ErrWrongPassword: http.StatusPreconditionFailed,

	// ошибки, связанные с сессией
	ErrSessionNotFound: http.StatusForbidden,

	// ошибки, связанные с командой
	ErrTeamNotFound: http.StatusNotFound,
}

func ResolveErrorToCode(err error) (code int) {
	code, has := errorToCodeMap[err]
	if !has {
		return http.StatusInternalServerError
	}
	return code
}
