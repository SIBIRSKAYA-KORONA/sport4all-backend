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
	ErrTeamBadRole  = errors.New("unprocessable team role")

	// ошибки, связанные со встречей
	ErrMeetingNotFound                = errors.New("meeting not found")
	ErrMeetingStatusNotAcceptable     = errors.New("meeting status not acceptable")
	ErrMeetingTeamsAreAlreadyAssigned = errors.New("teams are already assigned")

	// ошибки, связанные с турнирами
	ErrTournamentNotFound            = errors.New("tournament not found")
	ErrTournamentBadRole             = errors.New("unprocessable tournament role")
	ErrTournamentStatusNotAcceptable = errors.New("tournament status not acceptable")
	ErrTournamentSystemNotAcceptable = errors.New("tournament system not acceptable")

	// ошибки, связанные с аттачами
	ErrBadFileUploadS3 = errors.New("unsuccessful file upload to s3")
	ErrBadFileDeleteS3 = errors.New("unsuccessful file delete on s3")
	ErrFileNotFound    = errors.New("not found file in db")

	// ошибки, связанные с навыками
	ErrSkillNotFound = errors.New("skill not found")
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
	ErrTeamBadRole:  http.StatusUnprocessableEntity,

	// ошибки, связанные со встречей
	ErrMeetingNotFound:                http.StatusNotFound,
	ErrMeetingStatusNotAcceptable:     http.StatusNotAcceptable,
	ErrMeetingTeamsAreAlreadyAssigned: http.StatusUnavailableForLegalReasons,

	// ошибки, связанные с турнирами
	ErrTournamentNotFound:            http.StatusNotFound,
	ErrTournamentBadRole:             http.StatusUnprocessableEntity,
	ErrTournamentStatusNotAcceptable: http.StatusNotAcceptable,
	ErrTournamentSystemNotAcceptable: http.StatusNotAcceptable,

	// ошибки, связанные с аттачами
	ErrBadFileUploadS3: http.StatusUnprocessableEntity,
	ErrBadFileDeleteS3: http.StatusNotFound,
	ErrFileNotFound:    http.StatusNotFound,

	// ошибки, связанные с навыками
	ErrSkillNotFound: http.StatusNotFound,
}

func ResolveErrorToCode(err error) (code int) {
	code, has := errorToCodeMap[err]
	if !has {
		return http.StatusInternalServerError
	}
	return code
}
