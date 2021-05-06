package repositories

import "sport4all/app/models"

type OcrRepository interface {
	GetTextByImage(protocolImage *models.ProtocolImage) (*[]models.PlayerStat, error)
}
