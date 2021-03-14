package docs

import (
	"sport4all/app/models"
)

// swagger:route POST /api/teams Teams PostApiTeamsRequest
// Создаем команду
// responses:
//   200: PostApiTeams200Response
//   401: PostApiTeams401Response

// 200, успешно создали команду
// swagger:response PostApiTeams200Response
type PostApiTeams200Response struct {
}

// 401, проблемы с авторизацией
// swagger:response PostApiTeams401Response
type PostApiTeams401Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}

// swagger:parameters PostApiTeamsRequest
type PostApiTeamsRequest struct {
	// Описание реквеста
	// in:body
	Body TeamsRequestBody
}

type TeamsRequestBody struct {
	// example: ЦСКА
	Name string `json:"name"`

	// example: Moscow
	Location string `json:"location" `

	// example: Один из ведущих футбольных клубов Москвы
	About string `json:"about"`
}

// swagger:route GET /api/teams Teams GetApiTeamsRequest
// Получаем список команд
// responses:
//   200: GetApiTeams200Response
//   401: GetApiTeams401Response

// 200, успешно получили список команд
// swagger:response GetApiTeams200Response
type GetApiTeams200Response struct {
	// in:body
	Body models.Teams
}

// 401, проблемы с авторизацией
// swagger:response GetApiTeams401Response
type GetApiTeams401Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}

// swagger:parameters GetApiTeamsRequest
type GetApiTeamsRequest struct {
	// Роль
	// in: query
	// example: player/owner
	ExampleQueryParameter string `json:"role"`
}

// swagger:route GET /api/teams/{tid}/members/search Teams GetTeamsInviteListRequest
// Получаем список пользователей для приглашения в команду
// responses:
//   200: GetTeamsInviteList200Response
//   401: GetTeamsInviteList401Response

// 200, успешно получили список пользователей для приглашения
// swagger:response GetTeamsInviteList200Response
type GetTeamsInviteList200Response struct {
	// in:body
	Body models.Users
}

// 401, проблемы с авторизацией
// swagger:response GetTeamsInviteList401Response
type GetTeamsInviteList401Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}

// swagger:parameters GetTeamsInviteListRequest
type GetTeamsInviteListRequest struct {
	// Часть ника
	// in: query
	// example: nickna
	NicknamePart string `json:"nicknamePart"`

	// Лимит
	// in: query
	// example: 5
	Limit string `json:"limit"`
}

// swagger:route POST /api/teams/{tid}/members/{uid} Teams PostTeamsInviteRequest
// Приглашаем пользователя в команду в качестве игрока
// responses:
//   200: PostTeamsInvite200Response
//   401: PostTeamsInvite401Response

// 200, успешно пригласили пользователя
// swagger:response PostTeamsInvite200Response
type PostTeamsInvite200Response struct {
}

// 401, проблемы с авторизацией
// swagger:response PostTeamsInvite401Response
type PostTeamsInvite401Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}

// swagger:parameters PostTeamsInviteRequest
type PostTeamsInviteRequest struct {
	// ID команды
	// in: path
	// example: 1
	Tid string `json:"tid"`

	// ID игрока
	// in: path
	// example: 10
	Uid string `json:"uid"`
}
