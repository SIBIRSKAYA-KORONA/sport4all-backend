package docs

// swagger:route POST /api/tournaments Tournaments PostApiTournamentsRequest
// Создаем турнир
// responses:
//   200: PostApiTournaments200Response
//   401: PostApiTournaments401Response

// 200, успешно создали турнир
// swagger:response PostApiTournaments200Response
type PostApiTournaments200Response struct {
}

// 401, проблемы с авторизацией
// swagger:response PostApiTournaments401Response
type PostApiTournaments401Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}

// swagger:parameters PostApiTournamentsRequest
type PostApiTournamentsRequest struct {
	// Описание реквеста
	// in:body
	Body TournamentsRequestBody
}

type TournamentsRequestBody struct {
	// example: Чемпионат мира
	Name string `json:"name"`

	// example: Moscow
	Location string `json:"location" `
}
