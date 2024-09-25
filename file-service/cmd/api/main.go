package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/anuragprog/notyoutube/file-service/configs"
	"github.com/anuragprog/notyoutube/file-service/handlers"
	"github.com/anuragprog/notyoutube/file-service/middlewares"
	loggerRepo "github.com/anuragprog/notyoutube/file-service/repository/logger"
	loggerImplRepo "github.com/anuragprog/notyoutube/file-service/repository_impl/logger"
)

var (
	appLogger = loggerImplRepo.NewZeroLogger(
		loggerImplRepo.NewBatchLogger(
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

	app := SetupRouter(appLogger)
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

func SetupRouter(appLogger loggerRepo.Logger) *fiber.App {

	app := fiber.New(fiber.Config{
		ServerHeader: "not-youtube",
		ErrorHandler: handlers.GetErrorHandler(),
	})

	// setup middlewares
	// 1. setting up x-request-id header
	app.Use(middlewares.GetRequestIdMiddleware())
	// 2. setting up recover middleware (setting it before logger, so that when we catch a panic, will report as critical issue)
	app.Use(middlewares.GetRecoverMiddleware(appLogger))
	// 3. setting up logger (requires x-request-id, hence after 1st)
	app.Use(middlewares.GetLoggerMiddleware(appLogger))

	// setup api
	api := app.Group("/v1")

	api.Get("/health", handlers.GetHealthHandler())

	rawVideoGrp := api.Group("/raw-video")
	{
		rawVideoGrp.Get("", )
		rawVideoGrp.Post("", )
		rawVideoGrp.Patch("", )
	}

	return app
}
