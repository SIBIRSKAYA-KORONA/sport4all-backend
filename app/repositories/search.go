package repositories

import (
	"net"
	"sport4all/app/models"
)

type SearchRepository interface {
	GetGeo(ip net.IP) (*models.Location, error)
}
