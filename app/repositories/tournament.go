package repositories

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type TournamentRepository interface {
	Create(user *models.Tournament) error
}
