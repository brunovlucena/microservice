module github.com/brunovlucena/microservice/cmd/repository

go 1.13

replace (
	github.com/brunovlucena/microservice/cmd/data => ../data
	github.com/brunovlucena/microservice/cmd/utils => ../utils
)

require (
	github.com/brunovlucena/microservice/cmd/data v0.0.0-20191216092540-ad52db4f9356
	github.com/brunovlucena/microservice/cmd/messaging v0.0.0-20191216195322-87902ceb4bcd
	github.com/brunovlucena/microservice/cmd/utils v0.0.0-20191216092540-ad52db4f9356
	github.com/lib/pq v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.6.1
	github.com/stretchr/testify v1.4.0
)
