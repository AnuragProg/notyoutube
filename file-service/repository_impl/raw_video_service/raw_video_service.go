package raw_video_service

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	storeRepo "github.com/anuragprog/notyoutube/file-service/repository/store"
)

type rawVideoService struct {
	UnimplementedRawVideoServiceServer

	store *storeRepo.StoreManager
}

func NewRawVideoService(store *storeRepo.StoreManager) rawVideoService {
	return rawVideoService{
		store: store,
	}
}

func (rvs rawVideoService) GetRawVideoDownloadPresignedUrl(
	ctx context.Context,
	request *GetRawVideoDownloadPresignedUrlRequest,
) (*GetRawVideoDownloadPresignedUrlResponse, error) {
	presignResult, err := rvs.store.GetPresignedUrl(ctx, storeRepo.RAW_VIDEO, request.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &GetRawVideoDownloadPresignedUrlResponse{
		PresignedUrl: presignResult.Url,
		Method: presignResult.Method,
	}, nil
}
