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

	NoCookie        = "not found cookie header"
	SessionNotFound = "cookie invalid, session not exist"
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
	ErrNoCookie        = errors.New(NoCookie)
	ErrSessionNotFound = errors.New(SessionNotFound)
)

//var messToError = map[string]error{
//	Internal:     ErrInternal,
//	Conflict:     ErrConflict,
//	NoPermission: ErrNoPermission,
//
//	UserNotFound:  ErrUserNotFound,
//	WrongPassword: ErrWrongPassword,
//
//	NoCookie:        ErrNoCookie,
//	SessionNotFound: ErrSessionNotFound,
//}

var errorToCodeMap = map[error]int{
	// общие
	ErrInternal:     http.StatusInternalServerError,
	ErrConflict:     http.StatusConflict,
	ErrNoPermission: http.StatusForbidden,

	// ошибки, связанные с юзером
	ErrUserNotFound:  http.StatusNotFound,
	ErrWrongPassword: http.StatusPreconditionFailed,

	// ошибки, связанные с сессией
	ErrNoCookie:        http.StatusForbidden,
	ErrSessionNotFound: http.StatusNotFound,
}

func ResolveErrorToCode(err error) (code int) {
	code, has := errorToCodeMap[err]
	if !has {
		return http.StatusInternalServerError
	}
	return code
}

//func ResolveFromRPC(err error) error {
//	err, has := messToError[status.Convert(err).Message()]
//	if !has {
//		return ErrInternal
//	}
//	return err
//}
