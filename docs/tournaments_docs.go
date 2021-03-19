package docs

// swagger:route POST /api/tournaments Tournaments PostApiTournamentsRequest
// Создаем турнир
// responses:
//   200: PostApiTournaments200Response
//   401: General401Response
//   500: General500Response

// 200, успешно создали турнир
// swagger:response PostApiTournaments200Response
type PostApiTournaments200Response struct {
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
