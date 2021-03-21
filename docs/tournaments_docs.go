package docs

import "sport4all/app/models"

// ------------------------------------------------ СОЗДАНИЕ ТУРНИРА ---------------------------------------------------

// swagger:route POST /api/tournaments Tournaments PostApiTournamentsRequest
// Создаем турнир
// responses:
//   200: PostApiTournaments200Response
//   401: General401Response
//   500: General500Response

// 200, успешно создали турнир
// swagger:response PostApiTournaments200Response
type PostApiTournaments200Response struct {
	// Описание ответа
	// in:body
	Body models.Tournaments
}

// swagger:parameters PostApiTournamentsRequest
type PostApiTournamentsRequest struct {
	// Описание запроса
	// in:body
	Body TournamentsRequestBody
}

type TournamentsRequestBody struct {
	// example: Чемпионат мира
	Name string `json:"name" gorm:"index" faker:"name"`

	// example: Moscow
	Location string `json:"location" gorm:"index"`

	// example: olympic
	System string `json:"system"`

	// example: турнир по игре с котиками
	About string `json:"about"`
}
