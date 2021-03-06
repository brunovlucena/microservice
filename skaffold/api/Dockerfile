FROM golang:1.13 AS builder
# Use ADD in case of a URL only.
COPY ./cmd/api /api
COPY ./cmd/data /data
COPY ./cmd/utils /utils
COPY ./cmd/messaging /messaging
WORKDIR /api
# CGO is disabled as it is not compatible with the alpine image. 
# To use CGO write normal Go code that imports a pseudo-package "C". 
# Use go -gcflags to adds more debug information for the debugger.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /main .

# final stage
FROM alpine:3.10
RUN apk --no-cache add ca-certificates=20190108-r0
COPY ./cmd/api/config.yaml .
COPY ./cmd/api/router/welcome.html /router/welcome.html
COPY --from=builder /main ./
RUN chmod +x ./main
ENTRYPOINT ["./main"]
EXPOSE 8000
