package database

import (
	"time"

	"github.com/google/uuid"
)

type WorkerType string

const (
	WorkerTypeVideoEncoder       WorkerType = "video-encoder"
	WorkerTypeAsciiEncoder       WorkerType = "ascii-encoder"
	WorkerTypeThumbnailGenerator WorkerType = "thumbnail-generator"
	WorkerTypeAssembler          WorkerType = "assembler"
	WorkerTypeVideoExtractor     WorkerType = "video-extractor"
	WorkerTypeAudioExtractor     WorkerType = "audio-extractor"
	WorkerTypeMetadataExtractor  WorkerType = "metadata-extractor"
)

type Dag struct {
	ID        uuid.UUID `json:"id"`
	TraceId   uuid.UUID `json:"trace_id"`
	CreatedAt time.Time `json:"created_at"`
}

type Dependency struct {
	ID    uuid.UUID `json:"id"`
	DagID uuid.UUID `json:"dag_id"`
}

type DependencySource struct {
	ID           uuid.UUID `json:"id"`
	DagID        uuid.UUID `json:"dag_id"`
	DependencyID uuid.UUID `json:"dependency_id"`
	SourceID     uuid.UUID `json:"source_id"`
}

type DependencyTarget struct {
	ID           uuid.UUID `json:"id"`
	DagID        uuid.UUID `json:"dag_id"`
	DependencyID uuid.UUID `json:"dependency_id"`
	TargetID     uuid.UUID `json:"target_id"`
}

type Worker struct {
	ID           uuid.UUID  `json:"id"`
	DagID        uuid.UUID  `json:"dag_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	WorkerType   WorkerType `json:"worker_type"`
	WorkerConfig []byte     `json:"worker_config"`
}

type AsciiEncoderConfig struct {
	Config struct {
		Encoding string `json:"encoding"`
		Width    uint32 `json:"width"`
		Height   uint32 `json:"height"`
		Fps      uint32 `json:"fps"`
	} `json:"ascii_encoder_config"`
}
func NewAsciiEncoderConfig(
	encoding string,
	width, height, fps uint32,
) AsciiEncoderConfig {
	config := AsciiEncoderConfig{}
	config.Config.Encoding = encoding
	config.Config.Width = width
	config.Config.Height = height
	config.Config.Fps = fps

	return config
}

type VideoEncoderConfig struct {
	Config struct {
		Encoding string `json:"encoding"`
		Width    uint32 `json:"width"`
		Height   uint32 `json:"height"`
		Bitrate  uint32 `json:"bitrate"`
	} `json:"video_encoder_config"`
}
func NewVideoEncoderConfig(
	encoding string,
	width, height, bitrate uint32,
) VideoEncoderConfig {
	config := VideoEncoderConfig{}
	config.Config.Encoding = encoding
	config.Config.Width = width
	config.Config.Height = height
	config.Config.Bitrate = bitrate

	return config
}
