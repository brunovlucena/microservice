package messaging

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

var (
	log = "debug"
)

// Defines our interface for connecting and consuming messages.
type IMessagingClient interface {
	ConnectToBroker(connectionString string)
	Publish(msg []byte, exchangeName, exchangeType, queueName string) error
	Subscribe(exchangeName, exchangeType, consumerName, queueName string, handlerFunc func(amqp.Delivery)) error
	Close()
}

// Real implementation, encapsulates a pointer to an amqp.Connection
type MessagingClient struct {
	conn *amqp.Connection
}

func init() {
	// Setup Log
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	switch log {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	}
}

func (m *MessagingClient) ConnectToBroker(connectionString string) {
	if connectionString == "" {
		panic("Cannot initialize connection to broker, connectionString not set. Have you initialized?")
	}

	var err error
	m.conn, err = amqp.Dial(connectionString)
	if err != nil {
		panic("Failed to connect to AMQP compatible broker at: " + connectionString)
	} else {
		// Success
		logrus.WithFields(logrus.Fields{
			"cmd":  "ConnectToBroker",
			"amqp": connectionString,
		}).Infoln("Successfully connected to RabbitMQ!")
	}
}

func (m *MessagingClient) Publish(msg []byte, exchangeName, exchangeType, queueName string) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialized. Don't do that.")
	}

	ch, err := m.conn.Channel() // Get a channel from the connection
	defer ch.Close()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	failOnError(err, "Failed to register an Exchange")

	queue, err := ch.QueueDeclare( // Declare a queue that will be created if not exists with some args
		queueName, // our queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	err = ch.QueueBind(
		queue.Name,   // name of the queue
		exchangeName, // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)

	err = ch.Publish( // Publishes a message onto the queue.
		exchangeName, // exchange
		exchangeName, // routing key      q.Name
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        msg, // Our JSON body as []byte
		})

	logrus.Infoln("A message was sent: %v", string(msg))

	return err
}

func (m *MessagingClient) Subscribe(exchangeName, exchangeType, consumerName, queueName string, handlerFunc func(amqp.Delivery)) error {
	if m.conn == nil {
		panic("Tried to send message before connection was initialized. Don't do that.")
	}

	ch, err := m.conn.Channel()
	failOnError(err, "Failed to open a channel")

	err = ch.ExchangeDeclare(
		exchangeName, // name of the exchange
		exchangeType, // type
		true,         // durable
		false,        // delete when complete
		false,        // internal
		false,        // noWait
		nil,          // arguments
	)
	failOnError(err, "Failed to register an Exchange")

	logrus.Info("declared Exchange, declaring Queue (%s)", "")
	queue, err := ch.QueueDeclare(
		queueName, // name of the queue
		false,     // durable
		false,     // delete when usused
		false,     // exclusive
		false,     // noWait
		nil,       // arguments
	)
	failOnError(err, "Failed to register an Queue")

	logrus.Info("declared Queue (%d messages, %d consumers), binding to Exchange (key '%s')",
		queue.Messages, queue.Consumers, exchangeName)

	err = ch.QueueBind(
		queue.Name,   // name of the queue
		exchangeName, // bindingKey
		exchangeName, // sourceExchange
		false,        // noWait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		queue.Name,   // queue
		consumerName, // consumer
		true,         // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	failOnError(err, "Failed to register a consumer")

	go consumeLoop(msgs, handlerFunc)

	return nil
}

func consumeLoop(deliveries <-chan amqp.Delivery, handlerFunc func(d amqp.Delivery)) {
	for d := range deliveries {
		logrus.Printf("Received a message: %s", d.Body)
		handlerFunc(d)
		logrus.Printf("Done")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func (m *MessagingClient) Close() {
	if m.conn != nil {
		m.conn.Close()
	}
}
