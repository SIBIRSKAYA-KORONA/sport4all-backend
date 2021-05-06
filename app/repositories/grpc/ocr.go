package grpc

import (
	"context"
	"sport4all/app/repositories"
	"sport4all/app/repositories/grpc/proto"
	"sport4all/pkg/logger"

	"google.golang.org/grpc"
)

type OcrStore struct {
	client proto.OcrServiceClient
	context context.Context
}

func CreateOcrRepository(conn *grpc.ClientConn) repositories.OcrRepository {
	return &OcrStore{client: proto.NewOcrServiceClient(conn), context: context.Background()}
}

func (ocrStore *OcrStore) GetTextByImage()  {
	resp, err := ocrStore.client.GetStatsByImage(ocrStore.context, &proto.Image{Link: "/www/xxx/yyy/"})
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Info(resp.Stats[0])
}
