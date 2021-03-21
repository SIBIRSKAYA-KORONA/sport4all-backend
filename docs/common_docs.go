package docs

// 400, прислали гавно
// swagger:response General400Response
type General400Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Прислали гавно
		Message string
	}
}

// 401, проблемы с авторизацией
// swagger:response General401Response
type General401Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}

// 403, отсутствие прав
// swagger:response General403Response
type General403Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Отсутствует кука
		Message string
	}
}

// 500, внутренняя ошибка сервера
// swagger:response General500Response
type General500Response struct {
	// Описание
	// in: body
	Body struct {
		// The validation message

		// Required: true
		// Example: Ошибка на сервере
		Message string
	}
}
