package docs

import "github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"

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
