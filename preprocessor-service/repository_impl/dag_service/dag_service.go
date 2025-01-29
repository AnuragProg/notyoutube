package dag_service

import (
	context "context"
	"encoding/json"
	"errors"
	"sync"

	dbRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/database"
	"github.com/anuragprog/notyoutube/preprocessor-service/types"
	dbType "github.com/anuragprog/notyoutube/preprocessor-service/types/database"
	mqType "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"
	"github.com/anuragprog/notyoutube/preprocessor-service/utils"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type dagService struct {
	UnimplementedDAGServiceServer

	db dbRepo.Database
}

func (ds dagService) ListWorkersOfDAG(ctx context.Context, req *ListWorkersOfDAGRequest) (*ListWorkersOfDAGResponse, error) {
	dagId, err := uuid.Parse(req.GetDagId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dagid, uuid expected!")
	}

	dbWorkers, err := ds.db.ListWorkersOfDAG(ctx, dagId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	workers := utils.Map(
		dbWorkers,
		func(w dbType.Worker) *mqType.Worker {
			mqWorker := &mqType.Worker{
				Id: w.ID.String(),
				Name: w.Name,
				Description: w.Description,
				WorkerType: types.WorkerTypeToProtoWorkerType[w.WorkerType],
			}

			var videoEncoderConfig dbType.VideoEncoderConfig
			if err = json.Unmarshal(w.WorkerConfig, &videoEncoderConfig); err == nil {
				mqWorker.WorkerConfig = &mqType.Worker_VideoEncoderConfig{
					VideoEncoderConfig: &mqType.VideoEncoderWorkerConfig{
						Encoding: mqType.VideoEncoding(mqType.VideoEncoding_value[videoEncoderConfig.Config.Encoding]),
						Width: videoEncoderConfig.Config.Width,
						Height: videoEncoderConfig.Config.Height,
						Bitrate: videoEncoderConfig.Config.Bitrate,
					},
				}
			}

			var asciiEncoderConfig dbType.AsciiEncoderConfig
			if err = json.Unmarshal(w.WorkerConfig, &asciiEncoderConfig); err == nil {
				mqWorker.WorkerConfig = &mqType.Worker_AsciiEncoderConfig{
					AsciiEncoderConfig: &mqType.AsciiEncoderWorkerConfig{
						Encoding: mqType.AsciiEncoding(mqType.AsciiEncoding_value[asciiEncoderConfig.Config.Encoding]),
						Width: asciiEncoderConfig.Config.Width,
						Height: asciiEncoderConfig.Config.Height,
						Fps: asciiEncoderConfig.Config.Fps,
					},
				}
			}

			return mqWorker
		},
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &ListWorkersOfDAGResponse{
		Workers: workers,
	}, nil
}

func (ds dagService) ListDependenciesOfDAG(ctx context.Context, req *ListDependenciesOfDAGRequest) (*ListDependenciesOfDAGResponse, error) {
	dagId, err := uuid.Parse(req.GetDagId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dagid, uuid expected!")
	}

	var deps []dbType.Dependency
	var sources []dbType.DependencySource
	var targets []dbType.DependencyTarget

	var wg sync.WaitGroup
	errChan := make(chan error)
	wg.Add(1)
	go func(){
		defer wg.Done()
		deps, err = ds.db.ListDependenciesOfDAG(ctx, dagId)
		if err != nil {
			errChan<- err
		}
	}()
	wg.Add(1)
	go func(){
		defer wg.Done()
		sources, err = ds.db.ListDependencySourcesOfDAG(ctx, dagId)
		if err != nil {
			errChan<- err
		}
	}()
	wg.Add(1)
	go func(){
		defer wg.Done()
		targets, err = ds.db.ListDependencyTargetsOfDAG(ctx, dagId)
		if err != nil {
			errChan<- err
		}
	}()

	// close errchan on completion of all goroutines
	go func(){ 
		wg.Wait()
		close(errChan)
	}()

	// handle errors
	var errorsOccurred []error
	for err := range errChan {
		errorsOccurred = append(errorsOccurred, err)
	}
	if len(errorsOccurred) > 0 {
		return nil, status.Error(codes.Internal, errors.Join(errorsOccurred...).Error())
	}

	depMap := make(map[string]*mqType.Dependency)
	for _, dep := range deps {
		depMap[dep.ID.String()] = &mqType.Dependency{}
	}
	for _, src := range sources {
		dep := depMap[src.DependencyID.String()]
		dep.SourceIds = append(dep.SourceIds, src.SourceID.String())
	}
	for _, tgt := range targets {
		dep := depMap[tgt.DependencyID.String()]
		dep.TargetId = tgt.TargetID.String()
	}

	mqDeps := make([]*mqType.Dependency, 0, len(depMap))
	for _, dep := range depMap {
		mqDeps = append(mqDeps, dep)
	}

	return &ListDependenciesOfDAGResponse{
		Dependencies: mqDeps,
	}, nil
}

func (ds dagService) GetWorker(ctx context.Context, req *GetWorkerRequest) (*GetWorkerResponse, error) {
	workerId, err := uuid.Parse(req.GetWorkerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid workerid, uuid expected!")
	}

	worker, err := ds.db.GetWorkerById(ctx, workerId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	mqWorker := &mqType.Worker{
		Id: worker.ID.String(),
		Name: worker.Name,
		Description: worker.Description,
		WorkerType: types.WorkerTypeToProtoWorkerType[worker.WorkerType],
	}

	var videoEncoderConfig dbType.VideoEncoderConfig
	if err = json.Unmarshal(worker.WorkerConfig, &videoEncoderConfig); err == nil {
		mqWorker.WorkerConfig = &mqType.Worker_VideoEncoderConfig{
			VideoEncoderConfig: &mqType.VideoEncoderWorkerConfig{
				Encoding: mqType.VideoEncoding(mqType.VideoEncoding_value[videoEncoderConfig.Config.Encoding]),
				Width: videoEncoderConfig.Config.Width,
				Height: videoEncoderConfig.Config.Height,
				Bitrate: videoEncoderConfig.Config.Bitrate,
			},
		}
	}

	var asciiEncoderConfig dbType.AsciiEncoderConfig
	if err = json.Unmarshal(worker.WorkerConfig, &asciiEncoderConfig); err == nil {
		mqWorker.WorkerConfig = &mqType.Worker_AsciiEncoderConfig{
			AsciiEncoderConfig: &mqType.AsciiEncoderWorkerConfig{
				Encoding: mqType.AsciiEncoding(mqType.AsciiEncoding_value[asciiEncoderConfig.Config.Encoding]),
				Width: asciiEncoderConfig.Config.Width,
				Height: asciiEncoderConfig.Config.Height,
				Fps: asciiEncoderConfig.Config.Fps,
			},
		}
	}
	if err != nil { // checking for config parsing errors
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &GetWorkerResponse{
		Worker: mqWorker,
	}, nil
}

func (ds dagService) GetDependenciesWhereWorkerIsSource(ctx context.Context, req *GetDependenciesWhereWorkerIsSourceRequest) (*GetDependenciesWhereWorkerIsSourceResponse, error) {
	workerId, err := uuid.Parse(req.GetWorkerId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid dagid, uuid expected!")
	}

	depSrcsWhereWorkerIsSrc, err := ds.db.ListDependencySourcesWhereWorkerIsSource(ctx, workerId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	relatedDepIds := utils.Map(depSrcsWhereWorkerIsSrc, func(src dbType.DependencySource) uuid.UUID  {
		return src.DependencyID
	})

	// depSrcs, err 
	var sources []dbType.DependencySource
	var targets []dbType.DependencyTarget

	sources, err = ds.db.BatchListDependencySourcesOfDependency(ctx, relatedDepIds)
	if err != nil {
		return nil, err
	}
	targets, err = ds.db.BatchListDependencyTargetsOfDependency(ctx, relatedDepIds)
	if err != nil {
		return nil, err
	}
	
	depMap := make(map[string]*mqType.Dependency)
	for _, depId := range relatedDepIds {
		depMap[depId.String()] = new(mqType.Dependency)
	}
	for _, src := range sources {
		dep := depMap[src.DependencyID.String()]
		dep.SourceIds = append(dep.SourceIds, src.SourceID.String())
	}
	for _, tgt := range targets {
		dep := depMap[tgt.DependencyID.String()]
		dep.TargetId = tgt.TargetID.String()
	}

	deps := make([]*mqType.Dependency, 0, len(depMap))
	for _, dep := range depMap {
		deps = append(deps, dep)
	}

	return &GetDependenciesWhereWorkerIsSourceResponse{
		Dependencies: deps,
	}, nil
}
