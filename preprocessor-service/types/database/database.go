package database

import (
	"time"

	"github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/database/postgres"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	CreatedAt time.Time `json:"created_at"`
}

func (dag Dag) ToPostgresDag() postgres.Dag {
	return postgres.Dag{
		ID: pgtype.UUID{Bytes: dag.ID, Valid: true},
		CreatedAt: pgtype.Timestamp{Time: dag.CreatedAt, Valid: true},
	}
}

type Dependency struct {
	ID    uuid.UUID `json:"id"`
	DagID uuid.UUID `json:"dag_id"`
}

func (dep Dependency) ToPostgresDependency() postgres.Dependency{
	return postgres.Dependency{
		ID: pgtype.UUID{Bytes: dep.ID, Valid: true},
		DagID: pgtype.UUID{Bytes: dep.DagID, Valid: true},
	}
}

type DependencySource struct {
	ID           uuid.UUID `json:"id"`
	DagID        uuid.UUID `json:"dag_id"`
	DependencyID uuid.UUID `json:"dependency_id"`
	SourceID     uuid.UUID `json:"source_id"`
}

func (depSrc DependencySource) ToPostgresDependencySource() postgres.DependencySource {
	return postgres.DependencySource{
		ID: pgtype.UUID{Bytes: depSrc.ID, Valid: true},
		DagID: pgtype.UUID{Bytes: depSrc.DagID, Valid: true},
		DependencyID: pgtype.UUID{Bytes: depSrc.DependencyID, Valid: true},
		SourceID: pgtype.UUID{Bytes: depSrc.SourceID, Valid: true},
	}
}

type DependencyTarget struct {
	ID           uuid.UUID `json:"id"`
	DagID        uuid.UUID `json:"dag_id"`
	DependencyID uuid.UUID `json:"dependency_id"`
	TargetID     uuid.UUID `json:"target_id"`
}

func (depTgt DependencyTarget) ToPostgresDependencyTarget() postgres.DependencyTarget {
	return postgres.DependencyTarget{
		ID: pgtype.UUID{Bytes: depTgt.ID, Valid: true},
		DagID: pgtype.UUID{Bytes: depTgt.DagID, Valid: true},
		DependencyID: pgtype.UUID{Bytes: depTgt.DependencyID, Valid: true},
		TargetID: pgtype.UUID{Bytes: depTgt.TargetID, Valid: true},
	}
}

type Worker struct {
	ID           uuid.UUID  `json:"id"`
	DagID        uuid.UUID  `json:"dag_id"`
	Name         string     `json:"name"`
	Description  string     `json:"description"`
	WorkerType   WorkerType `json:"worker_type"`
	WorkerConfig []byte     `json:"worker_config"`
}

