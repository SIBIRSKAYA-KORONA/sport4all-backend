package errors

import (
	"errors"
	"net/http"
)

var (
	// общие
	ErrInternal     = errors.New("internal error")
	ErrConflict     = errors.New("conflict with exists data")
	ErrNoPermission = errors.New("no permission for current operation")

	// ошибки, связанные с юзером
	ErrUserNotFound  = errors.New("user not found")
	ErrWrongPassword = errors.New("wrong password")

	// ошибки, связанные с сессией
	ErrSessionNotFound = errors.New("cookie invalid or session not found")

	// ошибки, связанные с командой
	ErrTeamNotFound = errors.New("team not found")

	// ошибки, связанные со встречей
	ErrMeetingNotFound = errors.New("meeting not found")

	// ошибки, связанные с турнирами
	ErrTournamentNotFound            = errors.New("tournament not found")
	ErrTournamentStatusNotAcceptable = errors.New("tournament status not acceptable")
	ErrTournamentSystemNotAcceptable = errors.New("tournament system not acceptable")
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

	// ошибки, связанные со встречей
	ErrMeetingNotFound: http.StatusNotFound,

	// ошибки, связанные с турнирами
	ErrTournamentNotFound:            http.StatusNotFound,
	ErrTournamentStatusNotAcceptable: http.StatusNotAcceptable,
	ErrTournamentSystemNotAcceptable: http.StatusNotAcceptable,
}

func ResolveErrorToCode(err error) (code int) {
	code, has := errorToCodeMap[err]
	if !has {
		return http.StatusInternalServerError
	}
	return code
}
