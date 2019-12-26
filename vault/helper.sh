#!/usr/bin/env bash
[[ "$DEBUG" ]] && set -x # Print commands and their arguments as they are executed.

APP=microservice
NS=dev
TBL=configs
DB=configsdb
SVC=postgres.storage
ROLE=configs


kubernetes() {
    vault auth enable kubernetes

    vault write auth/kubernetes/config \
        token_reviewer_jwt="$(cat /var/run/secrets/kubernetes.io/serviceaccount/token)" \
        kubernetes_host=https://${KUBERNETES_PORT_443_TCP_ADDR}:443 \
        kubernetes_ca_cert=@/var/run/secrets/kubernetes.io/serviceaccount/ca.crt

    vault policy write "$APP" microservice.hcl

    vault write auth/kubernetes/role/"$APP"\
        bound_service_account_names="$APP" \
        bound_service_account_namespaces="$NS" \
        policies="$APP" \
        ttl=1h

    vault policy write kubernetes kubernetes.hcl
}

postgres_init() {
    local CONN="user=postgres password=postgres host=localhost port=5432 dbname=configsdb sslmode=disable"

    psql "$CONN" -c "CREATE ROLE "$APP" WITH LOGIN PASSWORD 'foo';" || true
    psql "$CONN" -c "GRANT ALL PRIVILEGES ON TABLE "$TBL" TO "$APP";" || true
}

postgres() {
    vault secrets enable database

    vault write database/config/postgresql \
        plugin_name=postgresql-database-plugin \
        allowed_roles="$ROLE" \
        connection_url="postgresql://{{username}}:{{password}}@"$SVC":5432/"$DB"?sslmode=disable" \
        username="postgres" \
        password="postgres"

    vault write database/static-roles/"$ROLE" \
        db_name=postgresql \
        rotation_statements=@rotation.sql \
        username="$APP" \
        rotation_period=86400
}

main() {
    TOKEN = vault token create -policy="$APP"
    VAULT_TOKEN="$TOKEN" vault read database/static-creds/"$ROLE"
}

kubernetes
postgres_init
postgres
