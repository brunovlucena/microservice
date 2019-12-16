### API In Golang

(Async) REST API in Go(1.13) using Chi Router

- Repository: Postgres
- Broker: RabbitMQ

**NOTE1**: Tested with minikube v1.5.1 and Kubernetes: v1.16.2

**NOTE2**: You should change `configs.yaml` (amdqAddr and dHost) to localhost or add to `127.0.0.1 postgres.storage` to `/etc/hosts`

**How to Test(Cluster)**

```sh
go get github.com/brunovlucena/microservice
make tunnel service=registry
make skaffold
make tunnel service=api
# Crud
make crud
# Load Test
make load-test (Broken)
```

**How to Test(locally)**

```sh
go get github.com/brunovlucena/microservice
make tunnel service=postgres
# Test
make test service=api
make test service=repository
# Load Test
make run service=api
make load-test (Broken)
# Crud
make run service=api
make crud
```

**Make Commands**

```sh
build                          Builds a binary for service (make build service=[api]).
check-pod-security             outputs infomation about the cluster(make check-pod-security label=[api] namespace=[dev])
checks                         Checks for erros in helper.sh
connect-postgres               Runs dlv (make debug service=[api]).
crud                           Perform simple crud operations (On Cluster).
debug                          Runs dlv (make debug service=[api]).
debug-tests                    Runs dlv test to debug Tests (make debug service=[api]).
deploy                         Deploys Api and Repository. (make deploy service=[api] version=v0.0.1 namespace=[dev])
deploy-infra-local             Runs infra on localhost.
help                           Help. 
load-examples                  Run Load Tests (make load-examples)
load-test                      Run Load Tests (make load-test service=[api])
run                            Runs service on localhost (make run service=[api]).
skaffold                       Uses skaffold during the development (make scaffold).
sniff                          Sniffs comunication (make sniff label=[api] namespace=[dev])
test-gui                       Runs Tests (Browser) (make test-gui service=[api]).
test                           Runs Tests (make test service=[api]).
tunnel-postgres                Creates a tunnel to minikube's registry.
tunnel-registry                Creates a tunnel to minikube's registry.
```


#### Infra Endpoints

**NOTE**: You should edit `/etc/hosts` ([minikube_ip] api.local)

- [Prometheus-Monitoring](http://api.local:31000)
- [API](ttp://api.local:32000)
- (Optional) [Microservice-Operator](https://github.com/brunovlucena/microservice-operator)


#### Non-Functional Requirements

1. Centralized configuration (DONE)
2. Service Discovery (TODO)
3. Logging (DONE)
4. Distributed Tracing (TODO)
5. Circuit Breaking (TODO)
7. Monitoring (DONE)
8. Security (TODO)


#### Endpoints

| Name   | Method      | URL
| ---    | ---         | ---
| List   | `GET`       | `/configs`
| Create | `POST`      | `/configs`
| Get    | `GET`       | `/configs/{name}`
| Update | `PUT/PATCH` | `/configs/{name}`
| Delete | `DELETE`    | `/configs/{name}`
| Query  | `GET`       | `/search?metadata.key=value`


#### Query

The query endpoint **MUST** return all configs that satisfy the query argument.

Query example-1:

```sh
curl http://localhost:8000/configs/search?metadata.monitoring.enabled=true
```

Response example:

```json
[
  {
    "name": "foo",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "false",
          "value": "300m"
        }
      }
    }
  },
  {
    "name": "bar",
    "metadata": {
      "monitoring": {
        "enabled": "true"
      },
      "limits": {
        "cpu": {
          "enabled": "true",
          "value": "250m"
        }
      }
    }
  },
]
```


#### Schema

- **Config**:
  - Name (string)
  - Metadata (nested key:value pairs where both key and value are strings of arbitrary length)


#### Configuration

- **API Variables**:

| Variable | Type | Example | Description |
| -------- | ---- | ------- | ----------- |
|`SERVER_ADDR`| string | ":8000"			| ## Addres to Listen to

- **Database Variables**:

| Variable | Type | Example | Description |
| -------- | ---- | ------- | ----------- |
|`DATABASE_TYPE`| string | "postgres"			| ## Database type
|`DATABASE_HOST`| string | "postgres.storage"	| ## Database location
|`DATABASE_PORT`| string | "5432"				| ## Port
|`DATABASE_USER`| string | "postgres"			| ## User
|`DATABASE_PASS`| string | "postgres"			| ## Pass
|`DATABASE_NAME`| string | "myapp"				| ## Database name
