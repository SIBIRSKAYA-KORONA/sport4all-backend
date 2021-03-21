package docs

// swagger:route PUT /api/tournaments/{tournamentId}/meetings Meetings GenerateMeetingsRequest
// Генерация встреч
// responses:
//   200: GenerateMeetings200
//   401: General401Response
//   500: General500Response

// 200, успешно сгенерировали встречи
// swagger:response GenerateMeetings200
type GenerateMeetings200 struct{}

// swagger:parameters GenerateMeetingsRequest
type GenerateMeetingsRequest struct{}
