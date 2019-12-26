module github.com/brunovlucena/microservice/cmd/api

go 1.13

replace (
	github.com/brunovlucena/microservice/cmd/data => ../data
	github.com/brunovlucena/microservice/cmd/messaging => ../messaging
	github.com/brunovlucena/microservice/cmd/utils => ../utils
)

require (
	github.com/brunovlucena/microservice/cmd/data v0.0.0-20191216092540-ad52db4f9356
	github.com/brunovlucena/microservice/cmd/messaging v0.0.0-20191216195322-87902ceb4bcd
	github.com/brunovlucena/microservice/cmd/utils v0.0.0-20191216092540-ad52db4f9356
	github.com/go-chi/chi v4.0.2+incompatible
	github.com/go-chi/render v1.0.1
	github.com/prometheus/client_golang v1.2.1
	github.com/rs/cors v1.7.0
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/goconvey v1.6.4
	github.com/spf13/viper v1.6.1
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/stretchr/testify v1.4.0 // indirect
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980
)
