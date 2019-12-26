package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	. "github.com/brunovlucena/microservice/cmd/data"
	. "github.com/brunovlucena/microservice/cmd/messaging"
	. "github.com/brunovlucena/microservice/cmd/repository/postgres"
	"github.com/getsentry/sentry-go"
	"github.com/streadway/amqp"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	appName = "repositoryApp"
)

var (
	client           IMessagingClient
	connectionString string
	repo             Repository
	dType            string
	dHost            string
	dPort            string
	dUser            string
	dPass            string
	dbName           string
	amqpAddr         string
)

func init() {
	// Setup Config path
	viper.SetConfigName("config.yaml")
	viper.AddConfigPath("/etc/appname/")
	viper.AddConfigPath("$HOME/.appname")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// Read Config

	// Broker
	if connectionString = os.Getenv("AMQP_ADDR"); connectionString == "" {
		connectionString = viper.Get("amqpAddr").(string)
	}
	client = &MessagingClient{}

	// Logrus
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	switch viper.Get("log") {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Postgres
	if dType = os.Getenv("DATABASE_TYPE"); dType == "" {
		dType = viper.Get("dType").(string)
	}
	if dHost = os.Getenv("DATABASE_HOST"); dHost == "" {
		dHost = viper.Get("dHost").(string)
	}
	if dPort = os.Getenv("DATABASE_PORT"); dPort == "" {
		dPort = viper.Get("dPort").(string)
	}
	if dbName = os.Getenv("DATABASE_NAME"); dbName == "" {
		dbName = viper.Get("dbName").(string)
	}

	// Credentials
	if vault := viper.Get("vault").(bool); vault {
		// Vault
		viper.SetConfigName("configs")
		viper.AddConfigPath("/vault/secrets")
		err := viper.ReadInConfig()
		if err != nil {
			panic(fmt.Errorf("Fatal error config file: %s \n", err))
		}
		dUser = viper.Get("username").(string)
		dPass = viper.Get("password").(string)
	} else {
		if dUser = os.Getenv("DATABASE_USER"); dUser == "" {
			dUser = viper.Get("dUser").(string)
		}
		if dPass = os.Getenv("DATABASE_PASS"); dPass == "" {
			dPass = viper.Get("dPass").(string)
		}
	}

	logrus.WithFields(logrus.Fields{
		"dType": dType,
		"dHost": dHost,
		"dPort": dPort,
		"dUser": dUser,
		"dPass": dPass,
	}).Debug("Parameters Loaded")

}

func main() {
	// Sentry
	sentry.Init(sentry.ClientOptions{
		Dsn: "https://a089a970262f463992edbe6e2008b243@sentry.io/1865784",
	})
	// TODO: TEST
	sentry.CaptureException(errors.New("my error"))
	sentry.Flush(time.Second * 5)

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

	// RabbitMQ
	exchangeName := "api"
	exchangeType := "fanout"
	queueName := "requests"

	logrus.Info("Starting " + appName)
	client.ConnectToBroker(connectionString)

	forever := make(chan bool)

	client.Subscribe(exchangeName, exchangeType, appName, queueName, handler)

	<-forever
}

func handler(d amqp.Delivery) {
	var cr *ConfigRequest
	err := json.Unmarshal(d.Body, &cr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd": "handler",
		}).Debug(err.Error())
	}
	path := cr.Path
	method := cr.Method
	result := &Config{}
	switch strings.ToLower(method) {
	case "get":
		if path.Path == "/search" {
			//params, _ := url.ParseQuery(cr.Path.URL.String())
			//configs := repo.Search()
		} else if path.Path == "/configs" {
			//configs := repo.FindAll()
		} else {
			result, err = repo.Find(cr.Config.Data["name"].(string))
		}
	case "post":
		result, _ = repo.Create(&cr.Config)
	case "put":
	case "patch":
		result, _ = repo.Update(&cr.Config)
	case "delete":
		result, _ = repo.Remove(cr.Config.Data["name"].(string))
	}

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd":    "handler",
			"result": result,
		}).Debug(err.Error())
	}

	jBytes, err := json.Marshal(&result)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"cmd": "handler",
		}).Debug(err.Error())
	}

	exchangeName := "api"
	exchangeType := "fanout"
	queueName := "responses"
	err = client.Publish(jBytes, exchangeName, exchangeType, queueName)
	if err != nil {

	}
}
