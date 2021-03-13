package impl

import (
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/models"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/repositories"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/app/usecases"
	"github.com/SIBIRSKAYA-KORONA/sport4all-backend/pkg/logger"
)

type TournamentUseCaseImpl struct {
	userRepo       repositories.UserRepository
	tournamentRepo repositories.TournamentRepository
}

func CreateTournamentUseCase(userRepo repositories.UserRepository, tournamentRepo repositories.TournamentRepository) usecases.TournamentUseCase {
	return &TournamentUseCaseImpl{userRepo: userRepo, tournamentRepo: tournamentRepo}
}

func (tournamentUseCase *TournamentUseCaseImpl) Create(ownerId uint, tournament *models.Tournament) error {
	_, err := tournamentUseCase.userRepo.GetByID(ownerId)
	if err != nil {
		logger.Error(err)
		return err
	}

	tournament.OwnerId = ownerId

	if err := tournamentUseCase.tournamentRepo.Create(tournament); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
