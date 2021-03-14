package repositories

import (
	"sport4all/app/models"
)

type SessionRepository interface {
	Create(session *models.Session) error
	Get(sid string) (uint, error)
	Delete(sid string) error
}
