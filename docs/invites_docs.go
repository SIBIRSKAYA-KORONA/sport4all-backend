package docs

import "sport4all/app/models"

// ---------------------------------------------- Создание приглашения -------------------------------------------------

// swagger:route POST /api/invites Invites CreateInviteRequest
//
// Создание приглашение
//
// Нужно вписать тип приглашения: "direct" - от владельца к игроку, "indirect" - от игрока к владельцу
//
// responses:
//   200: CreateInvite200
//   401: General403Response
//   500: General500Response

// 200, успешно создали приглашение
// swagger:response CreateInvite200
type CreateInvite200 struct{}

// swagger:parameters CreateInviteRequest
type CreateInviteRequest struct {
	// Описание запроса
	// in:body
	Body models.Invite
}

// ----------------------------------------- Отвечаем на приглашение ---------------------------------------------------

// swagger:route POST /api/invites/{inviteId} Invites AcceptInviteRequest
//
// Принимаем/отклоняем приглашение
//
// Передаем статус с решением
//
// responses:
//   200: AcceptInvite200Response
//   401: General401Response
//   500: General500Response

// 200, успешно ответили на приглашение
// swagger:response AcceptInvite200Response
type AcceptInvite200Response struct{}

// swagger:parameters AcceptInviteRequest
type AcceptInviteRequest struct {
	// Ответ на приглашение: 1 - Rejected, 2 - Accepted
	// in:body
	// example: 2
	State uint `json:"state"`
}

// ----------------------------------------- Получение списка приглашений ----------------------------------------------

// swagger:route GET /api/invites Invites GetInvitesRequest
//
// Получаем список приглашений
//
// Должна быть выставлена кука
//
// responses:
//   200: GetInvites200Response
//   401: General401Response
//   500: General500Response

// 200, успешно выставили результат команды
// swagger:response GetInvites200Response
type GetInvites200Response struct {
	// Список приглашений
	// in:body
	Body models.Invite
}

// swagger:parameters GetInvitesRequest
type GetInvitesRequest struct{}
