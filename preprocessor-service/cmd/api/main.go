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
	storeRepo "github.com/anuragprog/notyoutube/preprocessor-service/repository/store"

	databaseRepoImpl "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/database"
	loggerRepoImpl "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/logger"
	mqRepoImpl "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/mq"
	storeRepoImpl "github.com/anuragprog/notyoutube/preprocessor-service/repository_impl/store"

	mqType "github.com/anuragprog/notyoutube/preprocessor-service/types/mq"
	"github.com/gofiber/fiber/v2"
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
	} else {
		store = storeRepoImpl.MustNewMinioStore(configs.MINIO_URI, configs.MINIO_SERVER_ACCESS_KEY, configs.MINIO_SERVER_SECRET_KEY)
	}
	defer store.Close()
	var storeManager = storeRepo.NewStoreManager(configs.STORE_BUCKET, store)

	var db databaseRepo.Database
	if configs.USE_NOOP_DB {
		db = databaseRepoImpl.NewNoopDatabse()
	} else {
		db = databaseRepoImpl.MustNewMongoDatabase(configs.MONGO_URI, configs.MONGO_DB_NAME)
	}
	defer db.Close()

	var mq mqRepo.MessageQueue
	if configs.USE_NOOP_MQ {
		mq = mqRepoImpl.NewNoopQueue()
	} else {
		mq = mqRepoImpl.MustNewKafkaQueue(configs.KAFKA_BROKERS)
	}
	var mqManager = mqRepo.NewMessageQueueManager(mq)


	// setup worker listeners
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	errChan := mqManager.SubscribeToRawVideoTopic(ctx, func(rvm *mqType.RawVideoMetadata) error {
		return nil
	})
	go func(){
		for err := range errChan {
			panic(err)
		}
	}()


	app := SetupRouter(db, storeManager, appLogger, mqManager)
	doneChan := make(chan bool)

	go func() {
		if err := app.Listen(fmt.Sprintf(":%v", configs.API_PORT)); err != nil {
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
	storeManager *storeRepo.StoreManager,
	appLogger loggerRepo.Logger,
	mqManager *mqRepo.MessageQueueManager,
) *fiber.App {

	app := fiber.New(fiber.Config{
		ServerHeader: "not-youtube",
		BodyLimit:    50 << 20, // 50 mb
	})

	// setup middlewares
	// 1. setting up x-request-id header
	app.Use(middlewares.GetRequestIdMiddleware())
	// 2. setting up logger (requires x-request-id, hence after 1st), also handles panics and reports as critical errors
	app.Use(middlewares.GetLoggerMiddleware(appLogger))
	// 3. setting error response handler (make sure it is called before logger middleware to handle custom api errors)
	app.Use(middlewares.GetErrorResponseHandlerMiddleware())

	// setup api
	api := app.Group("/v1")

	api.Get("/health", handlers.GetHealthHandler())

	return app
}
