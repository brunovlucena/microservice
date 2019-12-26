#!/usr/bin/env bash
[[ "$DEBUG" ]] && set -x # Print commands and their arguments as they are executed.

set -e # Exit immediately if a command exits with a non-zero status.

# Build a binary.
#
# Usage:
#  $ ./helper.sh param1 [param2]
# * param1: build
# * param2: [api|repository]
build() {
    local APP="$1"

    cd "cmd/$APP"

    gochecknoglobals ./... || true

    go mod tidy
    go build -gcflags -m -o "../../build/_output/bin/$APP"
}

# Connect to Postgres.
#
# Usage:
#  $ ./helper.sh param1 [param2]
# * param1: load-examples
connect_postgres() {
    local CONN="user=postgres password=postgres host=localhost port=5432 dbname=configsdb sslmode=disable"

    psql "$CONN" -c "select * from configs;"
}

# Builds and Deploys to K8s.
#
# Usage:
#  $ ./helper.sh param1 [param2]
# * param1: build
# * param2: [api|repository]
deploy() {
    local APP="$1"
    local VERSION="$2"
    local NAMESPACE="$3"

    cd "cmd/$APP"

    local IMAGE="localhost:5000/$APP:$VERSION"

	docker rmi "$IMAGE" || true

	docker build -f Dockerfile -t "$IMAGE" .
    docker push "$IMAGE"
}

# x.
#
# Usage:
#  $ ./helper.sh param1
# * param1: debug
#
# b postgres/postgres_test.go:32
# c
debug() {
    local APP="$1"

    cd "cmd/$APP"

    dlv debug
}

# x.
#
# Usage:
#  $ ./helper.sh param1
# * param1: debug-test
#
# b postgres/postgres_test.go:32
# c
debug_tests() {
    local APP="$1"

    cd "cmd/$APP/postgres"

    dlv test
}

# Perform Load test.
#
# Usage:
#  $ ./helper.sh param1 [param2]
# * param1: load_test
# * param2: [api|repository]
load_test() {
    cd test
    k6 run -d 1s load.js
}

# Load tests on Postgres.
#
# Usage:
#  $ ./helper.sh param1 [param2]
# * param1: load-examples
load_examples() {
    CONN1="user=postgres password=postgres host=localhost port=5432 sslmode=disable"
    CONN2="user=postgres password=postgres host=localhost port=5432 dbname=configsdb sslmode=disable"
    DATABASE=configsdb

    psql "$CONN1" -c "create database $DATABASE" || true
    psql "$CONN2"< test/examples.sql || true
}

# Runs App.
#
# Usage:
#  $ ./helper.sh param1
# * param1: [api|repository]
run() {
    local APP="$1"

    cd "cmd/$APP"

    go mod tidy
    go run main.go
}

# Creates a tunnel to registry.
#
# Usage:
#  $ ./helper.sh skaffold param1
# * param1: [api|repository]
run_skaffold() {
    ENV=dev skaffold dev --cache-artifacts=false --watch-poll-interval=2000
}

# Run Tests.
#
# Usage:
#  $ ./helper.sh param1 [param2]
# * param1: test
# * param2: [api|repository]
test() {
    local APP="$1"

    cd "cmd/$APP"
    go test ./...
}

# Run Tests.
#
# Usage:
#  $ ./helper.sh param1 [param2]
# * param1: test-gui
# * param2: [api|repository]
test_gui() {
    true
}

# Creates a tunnel to registry.
#
# Usage:
#  $ ./helper.sh param1
# * param1: tunnel-registry
tunnel() {
    local ARG0="$1"
    case "$ARG0" in
        api)
	        kubectl port-forward "$(kubectl get pod -l component=api -o jsonpath='{.items[0].metadata.name}' -n dev)" 8000:8000 -n dev
        ;;
        registry)
	        kubectl port-forward "$(kubectl get pod -l actual-registry=true -o jsonpath='{.items[0].metadata.name}' -n kube-system)" 5000:5000 -n kube-system
        ;;
        postgres)
	        kubectl port-forward "$(kubectl get pod -l app.kubernetes.io/name=postgres -o jsonpath='{.items[0].metadata.name}' -n storage)" 5432:5432 -n storage
        ;;
        rabbitmq)
	        kubectl port-forward "$(kubectl get pod -l app.kubernetes.io/name=rabbitmq -o jsonpath='{.items[0].metadata.name}' -n storage)" 5672:5672 -n storage
        ;;
        rabbitmq-gui)
	        kubectl port-forward "$(kubectl get pod -l app.kubernetes.io/name=rabbitmq -o jsonpath='{.items[0].metadata.name}' -n storage)" 15672:15672 -n storage
        ;;
    esac
}

main() {
  local ARG0="$1"
  local ARG1="$2"
  local ARG2="$3"
  local ARG3="$4"
  case "$ARG0" in
    build)
        build "$ARG1"
    ;;
    connect-postgres)
        connect_postgres "$ARG1" "$ARG2"
    ;;
    debug)
        debug "$ARG1"
    ;;
    debug-tests)
        debug_tests "$ARG1"
    ;;
    deploy)
        deploy "$ARG1" "$ARG2" "$ARG3"
    ;;
    load-test)
        load_test "$ARG1"
    ;;
    load-examples)
        load_examples "$ARG1" "$ARG2"
    ;;
    run)
        run "$ARG1"
    ;;
    skaffold)
        run_skaffold
    ;;
    test)
        test "$ARG1"
    ;;
    test-gui)
        test_gui "$ARG1"
    ;;
    tunnel)
        tunnel "$ARG1"
    ;;
  esac
}

main "$@"
