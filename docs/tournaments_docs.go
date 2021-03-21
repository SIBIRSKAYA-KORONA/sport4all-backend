package docs

import "sport4all/app/models"

// ------------------------------------------------ СОЗДАНИЕ ТУРНИРА ---------------------------------------------------

// swagger:route POST /api/tournaments Tournaments CreateTournamentRequest
//
// Создаем турнир
//
// Передаем в теле json с нужными полями
//
// responses:
//   200: CreateTournament200Response
//   400: General400Response
//   401: General401Response
//   500: General500Response

// 200, успешно создали турнир
// swagger:response CreateTournament200Response
type CreateTournament200Response struct {
	// Описание ответа
	// in:body
	Body GetTournamentResponseBody
}

// swagger:parameters CreateTournamentRequest
type CreateTournamentRequest struct {
	// Авторизационная кука
	// in: header
	// example: session_id=215c5a74-efa3-41f9-8c27-55d8e13ecf64
	Cookie string `json:"Cookie"`

	// Описание запроса
	// in:body
	Body CreateTournamentRequestBody
}

type CreateTournamentRequestBody struct {
	// example: Чемпионат мира
	Name string `json:"name"`

	// example: Moscow
	Location string `json:"location"`

	// example: olympic
	System string `json:"system"`

	// example: турнир по игре с котиками
	About string `json:"about"`
}

// ------------------------------------------------ ПОЛУЧЕНИЕ ТУРНИРА --------------------------------------------------

// swagger:route GET /api/teams/{tournamentId} Tournaments GetTournamentByID
//
// Получаем турнир по его ID
//
// Передаем ID турнира в урле
//
// responses:
//   200: GetTournamentByID200Response
//   404: General404Response
//   500: General500Response

// 200, успешно получили турнир
// swagger:response GetTournamentByID200Response
type GetTournamentByID200Response struct {
	// Описание ответа
	// in:body
	Body GetTournamentResponseBody
}

// swagger:parameters GetTournamentByID
type CreateGetTournamentByIDRequest struct {
	// ID турнира
	// in: path
	// example: 1
	Tid string `json:"tournamentId"`
}

type GetTournamentResponseBody struct {
	// example: 10
	ID uint `json:"id"`

	// example: 4
	OwnerId uint `json:"ownerId"`

	// example: Чемпионат мира
	Name string `json:"name"`

	// example: Moscow
	Location string `json:"location"`

	// example: olympic
	System string `json:"system"`

	// example: 1
	Status models.EventStatus `json:"status"`

	// example: турнир по игре с котиками
	About string `json:"about"`

	// example: 1234
	Created int64 `json:"created"`
}

// ------------------------------------------------ ОБНОВЛЕНИЕ ТУРНИРА -------------------------------------------------

// swagger:route PUT /api/teams/{tournamentId} Tournaments UpdateTournamentByID
//
// Обновляем турнир по его ID
//
// Передаем ID турнира в урле и json в теле с полями, которые нужно обновить
//
// responses:
//   200: UpdateTournamentByID200Response
//   400: General400Response
//   401: General401Response
//   404: General404Response
//   406: General406Response
//   500: General500Response

// 200, успешно обновили турнир
// swagger:response UpdateTournamentByID200Response
type UpdateTournamentByID200Response struct {
}

// swagger:parameters UpdateTournamentByID
type UpdateTournamentByIDRequest struct {
	// Авторизационная кука
	// in: header
	// required: true
	// example: session_id=215c5a74-efa3-41f9-8c27-55d8e13ecf64
	Cookie string `json:"Cookie"`

	// ID турнира
	// in: path
	// example: 1
	Tid string `json:"tournamentId"`

	// Описание запроса
	// in:body
	Body UpdateTournamentRequestBody
}

type UpdateTournamentRequestBody struct {
	// example: Чемпионат мира
	Name string `json:"name"`

	// example: Moscow
	Location string `json:"location"`

	// example: olympic
	System string `json:"system"`

	// example: 1
	Status uint `json:"status"`

	// example: турнир по игре с котиками
	About string `json:"about"`
}

// ------------------------------------------------ ПОЛУЧЕНИЕ ВСЕХ КОМАНД ТУРНИРА --------------------------------------

// swagger:route GET /{tournamentId}/teams Tournaments GetAllTeams
//
// Получаем все команды турнира по его ID
//
// Передаем ID турнира в урле
//
// responses:
//   200: GetAllTeams200Response
//   400: General400Response
//   404: General404Response
//   500: General500Response

// 200, успешно обновили турнир
// swagger:response GetAllTeams200Response
type GetAllTeams200Response struct {
	// in:body
	Body models.Teams
}

// swagger:parameters GetAllTeams
type GetAllTeamsRequest struct {
	// ID турнира
	// in: path
	// example: 1
	Tid string `json:"tournamentId"`
}
