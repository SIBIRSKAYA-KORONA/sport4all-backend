package docs

// swagger:route PUT /api/tournaments/{tournamentId}/meetings Meetings GenerateMeetingsRequest
// Генерация встреч
// responses:
//   200: GenerateMeetings200
//   401: GenerateMeetings401

// 200, успешно сгенерировали встречи
// swagger:response GenerateMeetings200
type GenerateMeetings200 struct {}

// 401, проблемы с авторизацией
// swagger:response GenerateMeetings401
type GenerateMeetings401 struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}


// swagger:parameters GenerateMeetingsRequest
type GenerateMeetingsRequest struct {}

