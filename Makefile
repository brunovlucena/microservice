.PHONY: help test build deploy api repository skaffold

help: ## Help. 
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

# App
build: ## Builds a binary for service (make build service=[api]).
	@./helper.sh build ${service}

checks: ## Checks for erros in helper.sh
	@shellcheck helper.sh || true

crud: ## Perform simple crud operations (On Cluster).
	@./crud.sh

connect-postgres: ## Runs dlv (make debug service=[api]).
	@./helper.sh connect-postgres

debug: ## Runs dlv (make debug service=[api]).
	@./helper.sh debug ${service}

debug-tests: ## Runs dlv test to debug Tests (make debug service=[api]).
	@./helper.sh debug-tests ${service}

deploy: ## Deploys Api and Repository. (make deploy service=[api] version=v0.0.1 namespace=[dev])
	@./helper.sh deploy ${service} ${version} ${namespace}

deploy-infra-local: ## Runs infra on localhost.
	 @./helper.sh run-infra-local

load-test: ## Run Load Tests (make load-test service=[api])
	@./helper.sh load-test ${service}

load-examples: ## Run Load Tests (make load-examples)
	@./helper.sh load-examples

run: ## Runs service on localhost (make run service=[api]).
	@./helper.sh run ${service}

skaffold: ## Uses skaffold during the development (make scaffold).
	@./helper.sh skaffold

test: ## Runs Tests (make test service=[api]).
	@./helper.sh test ${service}

test-gui: ## Runs Tests (Browser) (make test-gui service=[api]).
	@./helper.sh test-gui ${service}

tunnel: ## Creates a tunnel to minikube's registry.
	@./helper.sh tunnel ${service}

# Security
check-pod-security: ## outputs infomation about the cluster(make check-pod-security label=[api] namespace=[dev])
	@./security.sh check-pod-security ${label} ${namespace}

sniff: ## Sniffs comunication (make sniff label=[api] namespace=[dev])
	@./security.sh sniff ${label} ${namespace}
