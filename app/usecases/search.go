package usecases

import (
	"net"
	"sport4all/app/models"
)

type SearchUseCase interface {
	GetResult(uid *uint, input *models.SearchInput) (*models.SearchOutput, error)
	GetGeo(ip net.IP) (*models.Location, error)
}
