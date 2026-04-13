#!/usr/bin/env bash
set -euo pipefail

IMAGE="${IMAGE:-test:local}"
CONTAINER_NAME="${CONTAINER_NAME:-http-probe-test-app-test}"

cleanup() {
  docker rm -f "${CONTAINER_NAME}" >/dev/null 2>&1 || true
}
trap cleanup EXIT

fail() {
  echo "FAIL: $*"
  exit 1
}

pass() {
  echo "PASS: $*"
}

require_cmd() {
  command -v "$1" >/dev/null 2>&1 || fail "missing required command: $1"
}

require_cmd docker
require_cmd curl

start_container() {
  cleanup

  local cid
  if (( $# > 0 )); then
    cid="$(docker run -d --rm --name "${CONTAINER_NAME}" -p 0:8080 "$@" "${IMAGE}")"
  else
    cid="$(docker run -d --rm --name "${CONTAINER_NAME}" -p 0:8080 "${IMAGE}")"
  fi

  local port
  port="$(docker port "${cid}" 8080/tcp | awk -F: '{print $NF}' | tail -n1 | tr -d '\r')"
  if [[ -z "${port}" ]]; then
    fail "unable to determine mapped port"
  fi

  echo "http://127.0.0.1:${port}"
}

http_code() {
  local url="$1"
  curl -sS -o /dev/null -w "%{http_code}" "${url}"
}

http_body() {
  local url="$1"
  curl -sS "${url}"
}

post_code() {
  local url="$1"
  curl -sS -X POST -o /dev/null -w "%{http_code}" "${url}"
}

wait_ready() {
  local url="$1"
  local deadline=$((SECONDS + 10))
  while (( SECONDS < deadline )); do
    if [[ "$(http_code "${url}")" == "200" ]]; then
      return 0
    fi
    sleep 0.2
  done
  return 1
}

base_url="$(start_container)"
if ! wait_ready "${base_url}/healthz"; then
  fail "container did not become reachable at ${base_url}"
fi

echo "--- Test 1: basic endpoints"
[[ "$(http_code "${base_url}/")" == "200" ]] || fail "GET / expected 200"
[[ "$(http_code "${base_url}/healthz")" == "200" ]] || fail "GET /healthz expected 200"
[[ "$(http_code "${base_url}/readyz")" == "200" ]] || fail "GET /readyz expected 200"
[[ "$(http_code "${base_url}/info")" == "200" ]] || fail "GET /info expected 200"
[[ "$(http_code "${base_url}/metrics")" == "200" ]] || fail "GET /metrics expected 200"

info="$(http_body "${base_url}/info")"
echo "${info}" | grep -q '"version"' || fail "/info missing version"
echo "${info}" | grep -q '"git_commit"' || fail "/info missing git_commit"
echo "${info}" | grep -q '"uptime"' || fail "/info missing uptime"
echo "${info}" | grep -q '"readiness"' || fail "/info missing readiness"

metrics="$(http_body "${base_url}/metrics")"
echo "${metrics}" | grep -q '^dummy_test_requests_total' || fail "/metrics missing dummy_test_requests_total"
pass "basic endpoints"

echo "--- Test 2: failure override"
[[ "$(http_code "${base_url}/healthz?fail=1")" == "500" ]] || fail "GET /healthz?fail=1 expected 500"
pass "failure override"

echo "--- Test 3: readiness toggle"
[[ "$(post_code "${base_url}/toggle-ready")" == "200" ]] || fail "POST /toggle-ready expected 200"
[[ "$(http_code "${base_url}/readyz")" == "503" ]] || fail "GET /readyz expected 503 after toggle"
[[ "$(post_code "${base_url}/toggle-ready")" == "200" ]] || fail "POST /toggle-ready expected 200 (toggle back)"
[[ "$(http_code "${base_url}/readyz")" == "200" ]] || fail "GET /readyz expected 200 after toggle back"
pass "readiness toggle"

echo "--- Test 4: liveness threshold"
base_url="$(start_container -e FAIL_LIVENESS_AFTER_N_REQUESTS=3)"
wait_ready "${base_url}/healthz" || fail "container not reachable for liveness threshold test"
curl -sS -o /dev/null "${base_url}/"
curl -sS -o /dev/null "${base_url}/"
curl -sS -o /dev/null "${base_url}/"
[[ "$(http_code "${base_url}/healthz")" == "500" ]] || fail "GET /healthz expected 500 after threshold"
pass "liveness threshold"

echo "--- Test 5: readiness delay"
base_url="$(start_container -e READY_DELAY_SECONDS=3)"
wait_ready "${base_url}/healthz" || fail "container not reachable for readiness delay test"
[[ "$(http_code "${base_url}/readyz")" == "503" ]] || fail "GET /readyz expected 503 immediately with delay"
sleep 4
[[ "$(http_code "${base_url}/readyz")" == "200" ]] || fail "GET /readyz expected 200 after delay elapsed"
pass "readiness delay"

echo "ALL TESTS PASSED"
