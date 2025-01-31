package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/anuragprog/notyoutube/preprocessor-service/configs"
	"github.com/anuragprog/notyoutube/preprocessor-service/handlers"
	"github.com/anuragprog/notyoutube/preprocessor-service/middlewares"
	databaseRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/database"
	loggerRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/logger"
	mqRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/mq"
	"github.com/anuragprog/notyoutube/preprocessor-service/utils"
	"github.com/anuragprog/notyoutube/preprocessor-service/workers"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	databaseRepoImpl "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/database"
	loggerRepoImpl "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/logger"
	mqRepoImpl "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/mq"

	mqType "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"
	rawVideoService "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/raw_video_service"
)

func main() {
	appLogger := loggerRepoImpl.NewZeroLogger(
		loggerRepoImpl.NewBatchLogger(
			os.Stdout,
			10<<10, // 10kb
			time.Second*10,
		),
		configs.SERVICE_NAME,
		string(configs.ENVIRONMENT),
	)
	defer appLogger.Close()

	var db databaseRepo.Database
	if configs.USE_NOOP_DB {
		db = databaseRepoImpl.NewNoopDatabase()
	} else {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		_db, err := databaseRepoImpl.NewPostgresDatabase(
			ctx,
			configs.POSTGRES_HOST,
			configs.POSTGRES_PORT,
			configs.POSTGRES_USER,
			configs.POSTGRES_PASSWORD,
			configs.POSTGRES_DBNAME,
		)
		if err != nil {
			panic(err)
		}
		db = _db
	}
	defer db.Close()

	var mq mqRepo.MessageQueue
	if configs.USE_NOOP_MQ {
		mq = mqRepoImpl.NewNoopQueue()
	} else {
		mq = utils.Must(mqRepoImpl.NewKafkaQueue(configs.KAFKA_BROKERS))
	}
	var mqManager = mqRepo.NewMessageQueueManager(mq)

	// setup worker listeners
	rawVideoServiceConn, err := grpc.NewClient(configs.RAW_VIDEO_SERVICE_URL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer rawVideoServiceConn.Close()
	rawVideoServiceClient := rawVideoService.NewRawVideoServiceClient(rawVideoServiceConn)
	rawVideoTopicCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := mqManager.SubscribeToRawVideoTopic(rawVideoTopicCtx, func(rvm *mqType.RawVideoMetadata) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		return workers.DAGWorker(ctx, db, mqManager, rawVideoServiceClient, rvm)
	})
	go func() {
		for err := range errChan {
			panic(err)
		}
	}()

	app := SetupRouter(db, appLogger, mqManager)
	doneChan := make(chan bool)

	go func() {
		if err := app.Start(fmt.Sprintf(":%v", configs.API_PORT)); err != nil {
			fmt.Println(err.Error())
		}
		doneChan <- true
	}()

	fmt.Printf("Server listening on %v\n", configs.API_PORT)
	<-doneChan
	fmt.Println("Server stopped")
}

// TODO: Not yet decided which apis to give for preprocessor but need to find out
func SetupRouter(
	db databaseRepo.Database,
	appLogger loggerRepo.Logger,
	mqManager *mqRepo.MessageQueueManager,
) *echo.Echo {

	app := echo.New()

	// setup middlewares
	// 1. setting up x-request-id header
	app.Use(middlewares.GetRequestIdMiddleware())
	app.Use(middlewares.GetTraceIdMiddleware())
	// 2. setting up logger (requires x-request-id, hence after 1st), also handles panics and reports as critical errors
	app.Use(middlewares.GetLoggerMiddleware(appLogger))
	// 3. setting error response handler (make sure it is called before logger middleware to handle custom api errors)
	app.Use(middlewares.GetErrorResponseHandlerMiddleware())

	// setup api
	api := app.Group("/v1")

	api.GET("/health", handlers.GetHealthHandler())

	return app
}
