package usecases

import (
	"sport4all/app/models"
)

type SearchUseCase interface {
	GetResult(uid *uint, input *models.SearchInput) (*models.SearchOutput, error)
}
