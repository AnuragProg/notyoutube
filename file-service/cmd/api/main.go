package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/anuragprog/notyoutube/file-service/configs"
	"github.com/anuragprog/notyoutube/file-service/handlers"
	"github.com/anuragprog/notyoutube/file-service/middlewares"
	databaseRepo "github.com/anuragprog/notyoutube/file-service/repository/database"
	loggerRepo "github.com/anuragprog/notyoutube/file-service/repository/logger"
	mqRepo "github.com/anuragprog/notyoutube/file-service/repository/mq"
	storeRepo "github.com/anuragprog/notyoutube/file-service/repository/store"
	"github.com/anuragprog/notyoutube/file-service/utils"
	"google.golang.org/grpc"

	databaseRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/database"
	loggerRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/logger"
	mqRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/mq"
	storeRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/store"
	rawVideoServiceRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/raw_video_service"
	"github.com/labstack/echo/v4"
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

	var store storeRepo.Store
	if configs.USE_NOOP_STORE {
		store = storeRepoImpl.NewNoopStore()
	}else {
		store = utils.Must(storeRepoImpl.NewMinioStore(configs.MINIO_URI, configs.MINIO_SERVER_ACCESS_KEY, configs.MINIO_SERVER_SECRET_KEY))
	}
	defer store.Close()
	var storeManager = storeRepo.NewStoreManager(configs.STORE_BUCKET, store)

	
	var db databaseRepo.Database
	if configs.USE_NOOP_DB {
		db = databaseRepoImpl.NewNoopDatabse()
	}else {
		db = utils.Must(databaseRepoImpl.NewMongoDatabase(configs.MONGO_URI, configs.MONGO_DB_NAME, configs.MONGO_RAW_VIDEO_COL))
	}
	defer db.Close()


	var mq mqRepo.MessageQueue 
	if configs.USE_NOOP_MQ {
		mq = mqRepoImpl.NewNoopQueue()
	}else {
		mq = utils.Must(mqRepoImpl.NewKafkaQueue(configs.KAFKA_BROKERS))
	}
	var mqManager = mqRepo.NewMessageQueueManager(mq)

	doneChan := make(chan bool)

	go func(){
		fmt.Printf("GRPC Server listening on %v\n", configs.GRPC_PORT)
		defer fmt.Println("GRPC Server stopped")
		grpcListener, err := net.Listen("tcp", fmt.Sprintf(":%v", configs.GRPC_PORT)) 
		if err != nil {
			panic(err.Error())
		}
		grpcApp := grpc.NewServer()
		if configs.USE_NOOP_RAW_VIDEO_SERVICE {
			rawVideoServiceRepoImpl.RegisterRawVideoServiceServer(
				grpcApp,
				rawVideoServiceRepoImpl.NewNoopRawVideoService(),
			)
		}else {
			rawVideoServiceRepoImpl.RegisterRawVideoServiceServer(
				grpcApp,
				rawVideoServiceRepoImpl.NewRawVideoService(storeManager),
			)
		}
		if err := grpcApp.Serve(grpcListener); err != nil {
			fmt.Println(err.Error())
		}
		doneChan<- true
	}()

	go func() {
		fmt.Printf("API Server listening on %v\n", configs.API_PORT)
		defer fmt.Println("API Server stopped")
		app := SetupRouter(db, storeManager, appLogger, mqManager)
		if err := app.Start(fmt.Sprintf(":%v", configs.API_PORT)); err != nil {
			fmt.Println(err.Error())
		}
		doneChan<- true
	}()

	<-doneChan
}

func SetupRouter(
	db databaseRepo.Database,
	storeManager *storeRepo.StoreManager,
	appLogger loggerRepo.Logger,
	mqManager *mqRepo.MessageQueueManager,
) *echo.Echo {

	app := echo.New()

	// setup middlewares
	// 1. setting up x-request-id and x-trace-id header
	app.Use(middlewares.GetRequestIdMiddleware())
	app.Use(middlewares.GetTraceIdMiddleware())
	// 2. setting up logger (requires x-request-id, hence after 1st), also handles panics and reports as critical errors
	app.Use(middlewares.GetLoggerMiddleware(appLogger))
	// 3. setting error response handler (make sure it is called before logger middleware to handle custom api errors)
	app.Use(middlewares.GetErrorResponseHandlerMiddleware())

	// setup api
	api := app.Group("/v1")

	api.GET("/health", handlers.GetHealthHandler())

	rawVideoGrp := api.Group("/raw-videos")
	{
		rawVideoGrp.POST("", handlers.PostRawVideoHandler(db, storeManager, mqManager))
		rawVideoGrp.GET("", handlers.GetRawVideoMetadatasHandler(db))
		rawVideoGrp.GET("/:videoId", handlers.GetRawVideoHandler(db, storeManager))
	}

	return app
}
