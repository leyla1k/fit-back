package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/folklinoff/fitness-app/cmd/app/processor"
	docs "github.com/folklinoff/fitness-app/docs"
)

// @title Fitness App API
// @version 1.0
// @description This is a fitness app server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host 158.160.62.249:8000
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	errors := make(chan error)
	go func() {
		errors <- processor.Run()
	}()
	docs.SwaggerInfo.Host = os.Getenv("HOST")
	stop := make(chan os.Signal)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errors:
		log.Println(err)
	case <-stop:
		processor.Shutdown(context.Background())
	}
}
