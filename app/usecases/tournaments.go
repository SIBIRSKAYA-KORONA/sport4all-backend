package usecases

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
)

type TournamentUseCase interface {
	Create(ownerId uint, tournament *models.Tournament) error
}
