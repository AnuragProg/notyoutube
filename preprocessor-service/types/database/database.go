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
