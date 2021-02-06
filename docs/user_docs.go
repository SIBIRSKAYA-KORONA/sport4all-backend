package docs

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

// swagger:route POST /api/settings/{uid} Settings PostApiSettingsRequest
// Регистрируем пользователя
// responses:
//   200: PostApiSettings200Response
//   404: PostApiSettings404Response

// 200, всё ок
// swagger:response PostApiSettings200Response
type PostApiSettings200Response struct {
	// in:body
	Body models.Users
}

// 404, такого нет, у тебя проблемы
// swagger:response PostApiSettings404Response
type PostApiSettings404Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Ппц ты лох
		Message string

		// Example: Думать надо было
		FieldName string
	}
}

// swagger:parameters PostApiSettingsRequest
type PostApiSettingsRequest struct {
	// Описание реквеста
	// in:body
	Body models.User

	// Количество чего-нибудь
	// in: query
	// example: 50/100/150
	ExampleQueryParameter string `json:"count"`

	// Идентификатор пользователя
	// in: path
	// example: 5445
	ExamplePathParameter string `json:"uid"`
}
