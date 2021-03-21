package docs

// ------------------------------------------------ СОЗДАНИЕ СЕССИИ ----------------------------------------------------

// swagger:route POST /api/sessions Sessions CreateSessionRequest
// Операция логина
// responses:
//   200: CreateSession200Response
//   400: General400Response
//   500: General500Response

// 200, успешно авторизовались
// swagger:response CreateSession200Response
type CreateSession200Response struct {
}

// swagger:parameters CreateSessionRequest
type CreateSessionRequest struct {
	// Описание запроса
	// in:body
	Body SessionRequestBody
}

type SessionRequestBody struct {
	// example: dendi
	Nickname string `json:"nickname"`

	// example: qwerty
	Password string `json:"password"`
}

// ------------------------------------------------ УДАЛЕНИЕ СЕССИИ ----------------------------------------------------

// swagger:route DELETE /api/sessions Sessions DeleteSessionRequest
//
// Операция логаута
//
// обязательно наличие куки
//
// responses:
//   200: DeleteSession200Response
//   401: General401Response
//   500: General500Response

// 200, успешно авторизовались
// swagger:response DeleteSession200Response
type DeleteSession200Response struct {
}

// swagger:parameters DeleteSessionRequest
type DeleteSessionRequest struct {
	// Авторизационная кука
	// in: header
	// required: true
	// example: session_id=215c5a74-efa3-41f9-8c27-55d8e13ecf64
	Cookie string `json:"Cookie"`
}
