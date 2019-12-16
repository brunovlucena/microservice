package main

import (
	"os"

	"github.com/brunovlucena/microservice/cmd/api/router"
	"github.com/sirupsen/logrus"
)

func init() {
	// Log Config
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Default is stderr
	logrus.SetOutput(os.Stdout)

	// Will log anything that is info or above (warn, error, fatal, panic).
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// Create a new Router
	r := router.NewRouter()

	// Start App
	r.StartWebServerHTTP()
}
