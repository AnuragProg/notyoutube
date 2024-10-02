package main


import (
	"fmt"
	"os"
	"time"

	"github.com/anuragprog/notyoutube/file-service/configs"
	"github.com/anuragprog/notyoutube/file-service/handlers"
	"github.com/anuragprog/notyoutube/file-service/middlewares"
	mqRepo "github.com/anuragprog/notyoutube/file-service/repository/mq"
	storeRepo "github.com/anuragprog/notyoutube/file-service/repository/store"
	loggerRepo "github.com/anuragprog/notyoutube/file-service/repository/logger"
	databaseRepo "github.com/anuragprog/notyoutube/file-service/repository/database"

	mqRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/mq"
	storeRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/store"
	loggerRepoImpl "github.com/anuragprog/notyoutube/file-service/repository_impl/logger"
	databaseRepoImpl"github.com/anuragprog/notyoutube/file-service/repository_impl/database"
	"github.com/gofiber/fiber/v2"
)

var (
	appLogger = loggerRepoImpl.NewZeroLogger(
		loggerRepoImpl.NewBatchLogger(
			os.Stdout,
			10<<10, // 10kb
			time.Second*10,
		), 
		configs.SERVICE_NAME,
		string(configs.ENVIRONMENT),
	)
)

func main(){
	defer appLogger.Close()

	minio, err := storeRepoImpl.NewMinioStore(configs.MINIO_URI, configs.MINIO_SERVER_ACCESS_KEY, configs.MINIO_SERVER_SECRET_KEY)
	if err != nil {
		panic(err)
	}
	defer minio.Close()
	storeManager := storeRepo.NewStoreManager(configs.STORE_BUCKET, minio)

	db, err := databaseRepoImpl.NewMongoDatabase(configs.MONGO_URI, configs.MONGO_DB_NAME, configs.MONGO_RAW_VIDEO_COL)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	mq, err := mqRepoImpl.NewKafkaClient(configs.KAFKA_BROKERS)
	if err != nil {
		panic(err)
	}

	app := SetupRouter(db, storeManager, appLogger, mq)
	doneChan := make(chan bool)

	go func(){
		if err := app.Listen(fmt.Sprintf(":%v", configs.API_PORT)); err != nil {
			fmt.Println(err.Error())
		}
		doneChan<- true
	}()

	fmt.Printf("Server listening on %v\n", configs.API_PORT)
	<-doneChan
	fmt.Println("Server stopped")
}

func SetupRouter(
	db databaseRepo.Database,
	storeManager *storeRepo.StoreManager,
	appLogger loggerRepo.Logger,
	mq mqRepo.MessageQueue,
) *fiber.App {

	app := fiber.New(fiber.Config{
		ServerHeader: "not-youtube",
		BodyLimit: 50<<20, // 50 mb
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

	rawVideoGrp := api.Group("/raw-videos")
	{
		rawVideoGrp.Post("", handlers.PostRawVideoHandler(db, storeManager))
		rawVideoGrp.Get("", handlers.GetRawVideoMetadatasHandler(db))
		rawVideoGrp.Get("/:video_id", handlers.GetRawVideoHandler(db, storeManager))
	}

	return app
}
