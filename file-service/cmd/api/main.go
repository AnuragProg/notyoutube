package main

import (
	"fmt"
	"os"
	"time"

	"github.com/anuragprog/notyoutube/file-service/configs"
	"github.com/anuragprog/notyoutube/file-service/handlers"
	"github.com/anuragprog/notyoutube/file-service/middlewares"
	"github.com/anuragprog/notyoutube/file-service/utils"
	"github.com/anuragprog/notyoutube/file-service/utils/log"
	"github.com/gofiber/fiber/v2"
)

var (
	appLogger = log.NewZeroLogger(
		log.NewBatchLogger(
			os.Stdout,
			10<<10, // 10kb
			time.Second*10,
		), 
		configs.SERVICE_NAME,
		string(utils.DEVELOPMENT_ENV),
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

func SetupRouter(appLogger log.Logger) *fiber.App {

	app := fiber.New(fiber.Config{
		ServerHeader: "Not-Youtube",
	})

	// setup middlewares
	// 1. setting up x-request-id header
	app.Use(middlewares.GetRequestIdMiddleware())
	// 2. setting up logger (requires x-request-id, hence after 1st)
	app.Use(middlewares.GetLoggerMiddleware(appLogger))

	// setup api
	api := app.Group("/v1")
	api.Get("/health", handlers.HealthHandler())
	
	return app
}
