module github.com/brunovlucena/microservice/cmd/messaging

go 1.13

require (
	github.com/brunovlucena/microservice/cmd/data v0.0.0-20191216092540-ad52db4f9356
	github.com/sirupsen/logrus v1.4.2
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
)

replace github.com/brunovlucena/microservice/cmd/data => ../data
