package repositories

import "sport4all/app/models"

type OcrRepository interface {
	GetStatsByImage(protocolImage *models.ProtocolImage) (*[]models.PlayerStat, error)
}
