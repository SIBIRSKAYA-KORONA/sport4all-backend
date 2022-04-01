package docs

import "sport4all/app/models"

// ---------------------------------------------- Генерация матчей -----------------------------------------------------

// swagger:route PUT /api/tournaments/{tournamentId}/meetings Meetings GenerateMeetingsRequest
//
// Генерация встреч
//
// Передаем идентификатор турнира
//
// responses:
//   200: GenerateMeetings200
//   401: General401Response
//   500: General500Response

// 200, успешно сгенерировали встречи
// swagger:response GenerateMeetings200
type GenerateMeetings200 struct{}

// swagger:parameters GenerateMeetingsRequest
type GenerateMeetingsRequest struct{}

// ---------------------------------------------- Генерация матчей -----------------------------------------------------

// ----------------------------------------- Получение статистика встречи ----------------------------------------------

// swagger:route GET /api/meetings/{meetingId} Meetings GetMeetingStatsRequest
//
// Запрос статистики встречи
//
// Передаем идентификатор встречи
//
// responses:
//   200: GetMeetingStats200Response
//   401: General401Response
//   500: General500Response

// 200, успешно получили статистику
// swagger:response GetMeetingStats200Response
type GetMeetingStats200Response struct {
	// in:body
	Body models.StatsSet
}

// swagger:parameters GetMeetingStatsRequest
type GetMeetingStatsRequest struct {
	// ID встречи
	// in: path
	// example: 12
	Mid string `json:"meetingId"`
}

// ----------------------------------------- Получение статистика встречи ----------------------------------------------

// --------------------------------------- Выставление результатов команды ---------------------------------------------

// swagger:route PUT /api/meetings/{meetingId}/teams/{teamId}/stat Meetings UpdateTeamStatsRequest
//
// Выставляем результат команды
//
// Передаем идентификатор встречи и команды
//
// responses:
//   200: UpdateTeamStats200Response
//   401: General401Response
//   500: General500Response

// 200, успешно выставили результат команды
// swagger:response UpdateTeamStats200Response
type UpdateTeamStats200Response struct{}

// swagger:parameters UpdateTeamStatsRequest
type UpdateTeamStatsRequest struct {
	// ID встречи
	// in: path
	// example: 12
	Mid string `json:"meetingId"`

	// ID команды
	// in: path
	// example: 4
	Tid string `json:"teamId"`

	// in:body
	Body TeamsRequestBody
}

type UpdateStatsBody struct {
	// example: 8
	Score string `json:"score"`
}

// --------------------------------------- Выставление результатов команды ---------------------------------------------

// ---------------------------------------- Выставление результатов игрока ---------------------------------------------

// swagger:route PUT /api/meetings/{meetingId}/teams/{teamId}/players/{playerId}/stat Meetings UpdatePlayerStatsRequest
//
// Выставляем результат игрока
//
// Передаем идентификатор встречи, команды и игрока
//
// responses:
//   200: UpdatePlayerStats200Response
//   401: General401Response
//   500: General500Response

// 200, успешно выставили результат команды
// swagger:response UpdatePlayerStats200Response
type UpdatePlayerStats200Response struct{}

// swagger:parameters UpdatePlayerStatsRequest
type UpdatePlayerStatsRequest struct {
	// ID встречи
	// in: path
	// example: 12
	Mid string `json:"meetingId"`

	// ID команды
	// in: path
	// example: 4
	Tid string `json:"teamId"`

	// ID игрока
	// in: path
	// example: 10
	Pid string `json:"playerId"`

	// in:body
	Body TeamsRequestBody
}

// ---------------------------------------- Выставление результатов игрока ---------------------------------------------
