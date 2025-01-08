package raw_video_service

import (
	"context"
)

type noopRawVideoService struct {
	UnimplementedRawVideoServiceServer
}

func NewNoopRawVideoService() noopRawVideoService {
	return noopRawVideoService{}
}

func (nrvs noopRawVideoService) GetRawVideoDownloadPresignedUrl(
	ctx context.Context,
	request *GetRawVideoDownloadPresignedUrlRequest,
) (*GetRawVideoDownloadPresignedUrlResponse, error) {
	return &GetRawVideoDownloadPresignedUrlResponse{}, nil
}
