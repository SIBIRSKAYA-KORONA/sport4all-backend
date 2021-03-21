package docs

import (
	"sport4all/app/models"
)

// ------------------------------------------------ СОЗДАНИЕ ПОЛЬЗОВАТЕЛЯ ----------------------------------------------

// swagger:route POST /api/settings Settings CreateSettingsRequest
//
// Регистрируем пользователя
//
// Передаем в теле json с нужными полями
//
// responses:
//   200: CreateSettings200Response
//   400: General400Response
//   404: General404Response
//   500: General500Response

// 200, создали пользователя
// swagger:response CreateSettings200Response
type CreateSettings200Response struct {
	// Авторизационная кука
	// in: header
	// required: true
	// example: session_id=215c5a74-efa3-41f9-8c27-55d8e13ecf64
	Cookie string `json:"Cookie"`
}

// swagger:parameters CreateSettingsRequest
type CreateSettingsRequest struct {
	// Описание запроса
	// in:body
	Body models.User
}

// ------------------------------------------------ ПОЛУЧЕНИЕ ПОЛЬЗОВАТЕЛЯ ---------------------------------------------

// swagger:route GET /api/settings Settings GetSettingsRequest
//
// Получаем пользовательские данные
//
// Передаем Cookie в заголовке запроса
//
// responses:
//   200: GetSettings200Response
//   401: General401Response
//   404: General404Response
//   500: General500Response

// 200, получили пользователя
// swagger:response GetSettings200Response
type GetSettings200Response struct {
	// Описание ответа
	// in:body
	Body models.User
}

// swagger:parameters GetSettingsRequest
type GetSettingsRequest CreateSettings200Response
