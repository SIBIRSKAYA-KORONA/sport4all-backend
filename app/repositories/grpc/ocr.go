package grpc

import (
	"context"
	"sport4all/app/models"
	"sport4all/app/repositories"
	"sport4all/app/repositories/grpc/proto"
	"sport4all/pkg/logger"

	"google.golang.org/grpc"
)

type OcrStore struct {
	client  proto.OcrServiceClient
	context context.Context
}

func CreateOcrRepository(conn *grpc.ClientConn) repositories.OcrRepository {
	return &OcrStore{client: proto.NewOcrServiceClient(conn), context: context.Background()}
}

func (ocrStore *OcrStore) GetTextByImage(protocolImage *models.ProtocolImage) (*[]models.PlayerStat, error) {
	resp, err := ocrStore.client.GetStatsByImage(ocrStore.context,
		&proto.Image{
			Link:         protocolImage.Link,
			PlayerColumn: protocolImage.PlayerColumn,
			ScoreColumn:  protocolImage.ScoreColumn,
		})

	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Info(resp.Stats[0]) // TODO remove it
	stats := make([]models.PlayerStat, len(resp.Stats))
	for idx, elem := range resp.Stats {
		stats[idx] = models.PlayerStat{Name: elem.Name, Surname: elem.Surname, Score: elem.Score}
	}

	return &stats, nil
}