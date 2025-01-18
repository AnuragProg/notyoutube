package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/anuragprog/notyoutube/dag-scheduler-service/types"
	mqType "github.com/anuragprog/notyoutube/dag-scheduler-service/types/mq"
	"github.com/anuragprog/notyoutube/dag-scheduler-service/utils"
	"github.com/google/uuid"

	dbRepo "github.com/anuragprog/notyoutube/dag-scheduler-service/repository/database"
	mqRepo "github.com/anuragprog/notyoutube/dag-scheduler-service/repository/mq"
	dbTypes "github.com/anuragprog/notyoutube/dag-scheduler-service/types/database"

	rawVideoServiceRepoImpl "github.com/anuragprog/notyoutube/dag-scheduler-service/repository_impl/raw_video_service"
)

// What it does: Downloads raw video file from file service.
// How: Requests file service for presigned url for the file.
// Makes get request on the presigned url and stores locally.
// Returns filename of the locally stored file.
func downloadRawVideoFileFromFileService(
	ctx context.Context,
	rawVideoServiceClient rawVideoServiceRepoImpl.RawVideoServiceClient,
	metadata *mqType.RawVideoMetadata,
) (string, error) {
	// Generating a unique file for local store
	filename, err := filepath.Abs(fmt.Sprintf("%s%s", uuid.NewString(), filepath.Ext(metadata.GetFilename())))
	if err != nil {
		return "", err
	}
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return "", err
	}
	defer file.Close()

	presignedUrl, err := rawVideoServiceClient.GetRawVideoDownloadPresignedUrl(
		ctx,
		&rawVideoServiceRepoImpl.GetRawVideoDownloadPresignedUrlRequest{
			Id: metadata.GetId(),
		},
	)
	if err != nil {
		return "", err
	}

	presignResponse, err := http.Get(presignedUrl.GetPresignedUrl())
	if err != nil {
		return "", err
	}
	defer presignResponse.Body.Close()

	if _, err = io.Copy(file, presignResponse.Body); err != nil {
		return "", err
	}
	return filename, nil
}

func DAGWorker(
	ctx context.Context,
	db dbRepo.Database,
	mq *mqRepo.MessageQueueManager,
	rawVideoServiceClient rawVideoServiceRepoImpl.RawVideoServiceClient,
	metadata *mqType.RawVideoMetadata,
) error {
	filename, err := downloadRawVideoFileFromFileService(ctx, rawVideoServiceClient, metadata)
	if err != nil {
		return err
	}
	defer os.Remove(filename)

	dag, err := utils.CreateDAG(ctx, filename)
	if err != nil {
		return err
	}

	// save everything in database
	err = db.WithTransaction(ctx, func(repo dbRepo.Database) error {

		dagId, err := uuid.Parse(dag.Id)
		if err != nil {
			return err
		}
		traceId, err := uuid.Parse(metadata.TraceId)
		if err != nil {
			return err
		}
		createdAt, err := time.Parse(time.RFC3339, dag.CreatedAt)
		if err != nil {
			return err
		}

		if err = repo.CreateDAG(ctx, dbTypes.Dag{
			ID:        dagId,
			TraceId:   traceId,
			CreatedAt: createdAt,
		}); err != nil {
			return err
		}

		workers := make([]dbTypes.Worker, 0, len(dag.GetWorkers()))
		for _, worker := range dag.GetWorkers() {
			var workerConfig []byte
			switch config := worker.GetWorkerConfig().(type){
			case *mqType.Worker_AsciiEncoderConfig:
				workerConfig, err = json.Marshal(map[string]any{
					"ascii_encoder_config": map[string]any{
						"encoding": mqType.AsciiEncoding_name[int32(config.AsciiEncoderConfig.GetEncoding())],
						"width": config.AsciiEncoderConfig.GetWidth(),
						"height": config.AsciiEncoderConfig.GetHeight(),
						"fps": config.AsciiEncoderConfig.GetFps(),
					},
				})
			case *mqType.Worker_VideoEncoderConfig:
				workerConfig, err = json.Marshal(map[string]any{
					"video_encoder_config": map[string]any{
						"encoding": mqType.VideoEncoding_name[int32(config.VideoEncoderConfig.GetEncoding())],
						"width": config.VideoEncoderConfig.GetWidth(),
						"height": config.VideoEncoderConfig.GetHeight(),
						"bitrate": config.VideoEncoderConfig.GetBitrate(),
					},
				})
			}
			if err != nil {
				return err
			}

			workers = append(workers, dbTypes.Worker{
				ID: uuid.MustParse(worker.GetId()),
				DagID: dagId,
				Name: worker.GetName(),
				Description: worker.GetDescription(),
				WorkerType: types.ProtoWorkerTypeToWorkerType[worker.GetWorkerType()],
				WorkerConfig: workerConfig,
			})
		}
		if err = repo.CreateWorkers(ctx, workers); err != nil {
			return err
		}

		dependencies := make([]dbTypes.Dependency, 0, len(dag.GetDependencies()))
		dependencySources := make([]dbTypes.DependencySource, 0, len(dag.GetDependencies())) // I know multiple sources hence will be >=; if u know how to make it exact number please do so
		dependencyTargets := make([]dbTypes.DependencyTarget, 0, len(dag.GetDependencies()))
		for _, dependency := range dag.GetDependencies() {
			dbDep := dbTypes.Dependency{
				ID: uuid.New(),
				DagID: dagId,
			}
			dependencies = append(dependencies, dbDep)
			for _, depSrc := range dependency.GetSourceIds() {
				dependencySources = append(dependencySources, dbTypes.DependencySource{
					ID: uuid.New(),
					DagID: dagId,
					DependencyID: dbDep.ID,
					SourceID: uuid.MustParse(depSrc),
				})
			}
			dependencyTargets = append(dependencyTargets, dbTypes.DependencyTarget{
				ID: uuid.New(),
				DagID: dagId,
				DependencyID: dbDep.ID,
				TargetID: uuid.MustParse(dependency.GetTargetId()),
			})
		}
		if err = repo.CreateDependencies(ctx, dependencies); err != nil {
			return err
		}

		if err = repo.CreateDependencySources(ctx, dependencySources); err != nil {
			return err
		}
		if err = repo.CreateDependencyTargets(ctx, dependencyTargets); err != nil {
			return err
		}

		// push to message queue
		return mq.PublishToDAGTopic(dag)
	})
	if err != nil {
		return err
	}
	return nil
}
