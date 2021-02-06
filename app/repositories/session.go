package repositories

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type SessionRepository interface {
	Create(session *models.Session) error
	Get(sid string) (uint, error)
	Delete(sid string) error
}
