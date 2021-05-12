package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/logger"
)

type SearchUseCaseImpl struct {
	teamRepo       repositories.TeamRepository
	tournamentRepo repositories.TournamentRepository
	userRepo       repositories.UserRepository
}

func CreateSearchUseCase(teamRepo repositories.TeamRepository,
	tournamentRepo repositories.TournamentRepository,
	userRepo repositories.UserRepository) usecases.SearchUseCase {
	return &SearchUseCaseImpl{teamRepo: teamRepo, tournamentRepo: tournamentRepo, userRepo: userRepo}
}

func (searchUseCase *SearchUseCaseImpl) GetResult(uid *uint, input *models.SearchInput) (*models.SearchOutput, error) {
	output := new(models.SearchOutput)
	if input.TeamQuery != nil {
		teams, err := searchUseCase.teamRepo.GetTeamsByNamePart(input.TeamQuery.Base.Text, 3)
		if err != nil {
			logger.Error(err)
		} else {
			output.Teams = teams
		}
	}
	if input.TournamentQuery != nil {
		tournaments, err := searchUseCase.tournamentRepo.GetTournamentsByNamePart(input.TournamentQuery.Base.Text, 3)
		if err != nil {
			logger.Error(err)
		} else {
			output.Tournaments = tournaments
		}
	}
	if input.UserQuery != nil {
		users, err := searchUseCase.userRepo.SearchUsers(uid, input.UserQuery.Base.Text, 3)
		if err != nil {
			logger.Error(err)
		} else {
			output.Users = (*models.Users)(users)
		}
	}
	return output, nil
}
