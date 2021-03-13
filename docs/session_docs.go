package docs

// swagger:route POST /api/sessions Sessions PostApiSessionsRequest
// Операция логина
// responses:
//   200: PostApiSessions200Response

// 200, успешно авторизовались
// swagger:response PostApiSessions200Response
type PostApiSessions200Response struct {
}

// swagger:parameters PostApiSessionsRequest
type PostApiSessionsRequest struct {
	// Описание реквеста
	// in:body
	Body SessionsRequestBody
}

type SessionsRequestBody struct {
	// example: dendi
	Nickname string `json:"nickname"`

	// example: qwerty
	Password string `json:"password"`
}
