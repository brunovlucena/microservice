package main

import (
	"fmt"
	"os"

	. "github.com/brunovlucena/microservice/cmd/messaging"
	. "github.com/brunovlucena/microservice/cmd/repository/postgres"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	appName = "repositoryApp"
)

var (
	client   IMessagingClient
	repo     Repository
	dType    string
	dHost    string
	dPort    string
	dUser    string
	dPass    string
	dbName   string
	amqpAddr string
)

func init() {
	// Setup Config path
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Read Config
	viper.SetConfigType("yaml")
	if dType = os.Getenv("DATABASE_TYPE"); dType == "" {
		dType = viper.Get("dType").(string)
	}
	if dHost = os.Getenv("DATABASE_HOST"); dHost == "" {
		dHost = viper.Get("dHost").(string)
	}
	if dPort = os.Getenv("DATABASE_PORT"); dPort == "" {
		dPort = viper.Get("dPort").(string)
	}
	if dUser = os.Getenv("DATABASE_USER"); dUser == "" {
		dUser = viper.Get("dUser").(string)
	}
	if dPass = os.Getenv("DATABASE_PASS"); dPass == "" {
		dPass = viper.Get("dPass").(string)
	}
	if dbName = os.Getenv("DATABASE_NAME"); dbName == "" {
		dbName = viper.Get("dbName").(string)
	}

	logrus.WithFields(logrus.Fields{
		"dType": dType,
		"dHost": dHost,
		"dPort": dPort,
		"dUser": dUser,
		"dPass": dPass,
	}).Debug("Parameters Loaded")

	// Broker
	var connectionString string
	if connectionString = os.Getenv("AMQP_ADDR"); connectionString == "" {
		connectionString = viper.Get("amqpAddr").(string)
	}
	client = &MessagingClient{}
	client.ConnectToBroker(connectionString)

	// Setup Log
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	switch viper.Get("log") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func main() {
	// PostgresSQL
	var err error
	repo, err = NewRepository(dType, dHost, dPort, dUser, dPass, dbName)

	// log error
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd": "NewRepository",
		}).Error(err.Error())
		panic(err.Error())
	}

	logrus.Info("Starting " + appName)
	forever := make(chan bool)

	<-forever
}
