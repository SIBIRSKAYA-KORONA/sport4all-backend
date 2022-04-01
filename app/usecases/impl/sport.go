package impl

import (
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/usecases"
	"sport4all/pkg/logger"
)

type SportUseCaseImpl struct {
	sportRepo repositories.SportRepository
}

var baseSports = map[string]bool{
	"Basketball":    true,
	"Football":      true,
	"Hockey":        true,
	"Table tennis":  true,
	"Tennis":        true,
	"Billiards":     true,
	"Chess":         true,
	"Checkers":      true,
	"Mini football": true,
	"Curling":       true,
}

func CreateSportUseCase(sportRepo repositories.SportRepository) usecases.SportUseCase {
	sportUseCase := SportUseCaseImpl{sportRepo: sportRepo}

	for sport := range baseSports {
		if err := sportUseCase.Create(&models.Sport{Name: sport}); err != nil {
			logger.Warn("error while create base sport: ", err)
		}
	}

	return &sportUseCase
}

func (sportUseCaseImpl *SportUseCaseImpl) Create(sport *models.Sport) error {
	if err := sportUseCaseImpl.sportRepo.Create(sport); err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (sportUseCaseImpl *SportUseCaseImpl) GetAll(limit uint) (*[]models.Sport, error) {
	sports, err := sportUseCaseImpl.sportRepo.GetAll(limit)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return sports, nil
}
