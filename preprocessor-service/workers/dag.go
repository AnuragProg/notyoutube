package workers

import (
	"time"
	"context"

	mqType "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"
	"github.com/anuragprog/notyoutube/preprocessor-service/utils"
	"github.com/google/uuid"

	dbRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/database"
	mqRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/mq"
	dbTypes "github.com/anuragprog/notyoutube/preprocessor-service/types/database"
)

func DAGWorker(ctx context.Context, mq *mqRepo.MessageQueueManager, database dbRepo.Database, metadata *mqType.RawVideoMetadata) error {
	// TODO: Download file and store it locally somewhere
	var filename string = "input.mp4"
	//////////////////////////////////

	dag, err := utils.CreateDAG(ctx, filename)
	if err != nil {
		return err
	}

	// save everything in database
	err = database.WithTransaction(ctx, func(repo dbRepo.Database) error {

		dagId, err := uuid.Parse(dag.Id)
		if err != nil {
			return err
		}
		createdAt, err := time.Parse(time.RFC3339, dag.CreatedAt) 
		if err != nil {
			return err
		}

		if err = repo.CreateDAG(ctx, dbTypes.Dag{
			ID: dagId,
			CreatedAt: createdAt,
		}); err != nil {
			return err
		}

		workers := make([]dbTypes.Worker, 0, len(dag.GetWorkers()))
		if err = repo.CreateWorkers(ctx, workers); err != nil { return err }

		dependencies := make([]dbTypes.Dependency, 0, len(dag.GetDependencies()))
		if err = repo.CreateDependencies(ctx, dependencies); err != nil { return err }

		dependencySources := make([]dbTypes.DependencySource, 0, len(dag.GetDependencies()))
		dependencyTargets := make([]dbTypes.DependencyTarget, 0, len(dag.GetDependencies()))
		if err = repo.CreateDependencySources(ctx, dependencySources); err != nil { return err }
		if err = repo.CreateDependencyTargets(ctx, dependencyTargets); err != nil { return err }

		// push to message queue
		return mq.PublishToDAGTopic(dag)
	})
	if err != nil {
		return err
	}

	return nil
}

