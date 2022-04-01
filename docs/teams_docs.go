package docs

import (
	"sport4all/app/models"
)

// ------------------------------------------------ СОЗДАНИЕ КОМАНДЫ ---------------------------------------------------

// swagger:route POST /api/teams Teams CreateTeamRequest
//
// Создаем команду
//
// Передаем в теле json с нужными полями
//
// responses:
//   200: CreateTeam200Response
//   401: General401Response
//   500: General500Response

// 200, успешно создали команду
// swagger:response CreateTeam200Response
type CreateTeam200Response struct{}

// swagger:parameters CreateTeamRequest
type CreateTeamRequest struct {
	// Описание реквеста
	// in:body
	Body TeamsRequestBody
}

// ------------------------------------------------ СОЗДАНИЕ КОМАНДЫ ---------------------------------------------------

type TeamsRequestBody struct {
	// example: ЦСКА
	Name string `json:"name"`

	// example: Moscow
	Location string `json:"location" `

	// example: Один из ведущих футбольных клубов Москвы
	About string `json:"about"`
}

// -------------------------------------------- ПОИСК КОМАНД ПО НАЗВАНИЮ -----------------------------------------------

// swagger:route GET /api/teams/search Teams SearchTeamsByNameRequest
//
// Ищем команды по части названия
//
// Передаем строчку и лимит
//
// responses:
//   200: SearchTeamsByName200Response
//   401: General401Response
//   500: General500Response

// 200, возвращаем список подходящих команд
// swagger:response SearchTeamsByName200Response
type SearchTeamsByName200Response struct {
	// in:body
	Body models.Teams
}

// swagger:parameters SearchTeamsByNameRequest
type SearchTeamsByNameRequest struct {
	// Часть названия команды
	// in: query
	// example: ЦСК
	NamePart string `json:"namePart"`

	// Лимит
	// in: query
	// example: 5
	Limit string `json:"limit"`
}

// -------------------------------------------- ПОИСК КОМАНД ПО НАЗВАНИЮ -----------------------------------------------

// --------------------------------------- ПОЛУЧЕНИЕ СПИСКА КОМАНД ПОЛЬЗОВАТЕЛЯ ----------------------------------------

// swagger:route GET /api/teams Teams GetTeamsRequest
//
// Получаем список команд пользователя
//
// Нужно передать параметр роли
//
// responses:
//   200: GetTeams200Response
//   401: General401Response
//   500: General500Response

// 200, успешно получили список команд
// swagger:response GetTeams200Response
type GetTeams200Response struct {
	// in:body
	Body models.Teams
}

// swagger:parameters GetTeamsRequest
type GetTeamsRequest struct {
	// Роль
	// in: query
	// example: player/owner
	ExampleQueryParameter string `json:"role"`
}

// --------------------------------------- ПОЛУЧЕНИЕ СПИСКА КОМАНД ПОЛЬЗОВАТЕЛЯ ----------------------------------------

// ------------------------------------ ПОИСК ПОЛЬЗОВАТЕЛЕЙ ДЛЯ ДОБАВЛЕНИЯ В КОМАНДУ -----------------------------------

// swagger:route GET /api/teams/{tid}/members/search Teams FindUsersForInviteRequest
//
// Получаем список пользователей для приглашения в команду
//
// Передаем часть никнейма и лимит
//
// responses:
//   200: FindUsersForInvite200Response
//   401: General401Response
//   500: General500Response

// 200, успешно получили список пользователей для приглашения
// swagger:response FindUsersForInvite200Response
type FindUsersForInvite200Response struct {
	// in:body
	Body models.Users
}

// swagger:parameters FindUsersForInviteRequest
type FindUsersForInviteRequest struct {
	// Часть ника
	// in: query
	// example: nickna
	NicknamePart string `json:"nicknamePart"`

	// Лимит
	// in: query
	// example: 5
	Limit string `json:"limit"`
}

// ------------------------------------ ПОИСК ПОЛЬЗОВАТЕЛЕЙ ДЛЯ ДОБАВЛЕНИЯ В КОМАНДУ -----------------------------------

// -------------------------------- ПРИГЛАШЕНИЕ ПОЛЬЗОВАТЕЛЯ В КОМАНДУ В КАЧЕСТВЕ ИГРОКА -------------------------------

// swagger:route POST /api/teams/{tid}/members/{uid} Teams InvitePlayerRequest
//
// Приглашаем пользователя в команду в качестве игрока
//
// Обязательно нужно быть владельцем команды
//
// responses:
//   200: InvitePlayer200Response
//   401: General401Response
//   403: General403Response
//   500: General500Response

// 200, успешно пригласили пользователя
// swagger:response InvitePlayer200Response
type InvitePlayer200Response struct{}

// swagger:parameters InvitePlayerRequest
type InvitePlayerRequest struct {
	// ID команды
	// in: path
	// example: 1
	Tid string `json:"tid"`

	// ID игрока
	// in: path
	// example: 10
	Uid string `json:"uid"`
}

// -------------------------------- ПРИГЛАШЕНИЕ ПОЛЬЗОВАТЕЛЯ В КОМАНДУ В КАЧЕСТВЕ ИГРОКА -------------------------------

// --------------------------------------- ИСКЛЮЧЕНИЕ ПОЛЬЗОВАТЕЛЯ ИЗ КОМАНДЫ ------------------------------------------

// swagger:route DELETE /api/teams/{tid}/members/{uid} Teams DeletePlayerRequest
//
// Исключаем пользователя из команды
//
// Обязательно нужно быть владельцем команды
//
// responses:
//   200: DeletePlayer200Response
//   401: General401Response
//   403: General403Response
//   500: General500Response

// 200, успешно исключили пользователя
// swagger:response DeletePlayer200Response
type DeletePlayer200Response struct{}

// swagger:parameters DeletePlayerRequest
type DeletePlayerRequest struct {
	// ID команды
	// in: path
	// example: 1
	Tid string `json:"tid"`

	// ID игрока
	// in: path
	// example: 10
	Uid string `json:"uid"`
}

// --------------------------------------- ИСКЛЮЧЕНИЕ ПОЛЬЗОВАТЕЛЯ ИЗ КОМАНДЫ ------------------------------------------
