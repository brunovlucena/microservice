#!/usr/bin/env bash
[[ "$DEBUG" ]] && set -x # Print commands and their arguments as they are executed.

set -e # Exit immediately if a command exits with a non-zero status.

# checks deployment security
#
# Usage:
#  $ ./security.sh check_pod_security param1
# * param1: it's the pod name
check_pod_security() {
    local LABEL="$1"
    local NAMESPACE="$2"
    local POD_NAME=$(kubectl get pod -l app="$LABEL" -o jsonpath='{.items[0].metadata.name}' -n "$NAMESPACE")
    kubectl kubesec-scan pod $POD_NAME -n "$NAMESPACE"
}

# sniffs pod
#
# Usage:
#  $ ./security.sh sniff
# * param1: it's the pod name
pod_sniff() {
    local LABEL="$1"
    local NAMESPACE="$2"
    local POD_NAME=$(kubectl get pod -l app="$LABEL" -o jsonpath='{.items[0].metadata.name}' -n "$NAMESPACE")
	kubectl sniff ${POD_NAME} -n default -c chart -o ksniff-dump.pcap -p
}

main() {
  local ARG0="$1"
  local ARG1="$2"
  local ARG2="$3"
  case "$ARG0" in
    check-pod-security)
        check_pod_security "$ARG1" "$ARG2"
    ;;
    sniff)
        pod_sniff "$ARG1" "$ARG2"
    ;;
  esac
}

main "$@"
