package psql

import (
	"github.com/jinzhu/gorm"
	"github.com/oschwald/geoip2-golang"
	"net"

	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/pkg/errors"
	"sport4all/pkg/logger"
)

type SearchStore struct {
	db        *gorm.DB
	geoReader *geoip2.Reader
}

func CreateSearchRepository(db *gorm.DB, geoReader *geoip2.Reader) repositories.SearchRepository {
	return &SearchStore{db: db, geoReader: geoReader}
}

func (searchStore *SearchStore) GetGeo(ip net.IP) (*models.Location, error) {
	record, err := searchStore.geoReader.City(ip)
	if err != nil {
		logger.Error(err)
		return nil, errors.ErrGeoNotFound
	}

	subdivision := ""
	if len(record.Subdivisions) != 0 {
		subdivision = record.Subdivisions[0].Names["ru"]
	}

	location := models.Location{
		City:        record.City.Names["ru"],
		Country:     record.Country.Names["ru"],
		Continent:   record.Continent.Names["ru"],
		Subdivision: subdivision,
	}

	return &location, nil
}
