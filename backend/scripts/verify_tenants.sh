#!/usr/bin/env bash
set -euo pipefail

BACKEND_URL="${BACKEND_URL:-http://localhost:${PORT:-8080}}"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "ERROR: '$1' is required but not installed." >&2
    exit 1
  fi
}

require_cmd curl

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1"; exit 1; }
case_header() { echo; echo "CASE: $1"; }

echo "== Tenant Verification =="
echo "Backend URL: $BACKEND_URL"

case_header "Health check"
echo "Expected: /ping returns pong"
ping_response="$(curl -fsS "$BACKEND_URL/ping" || true)"
echo "Actual: $ping_response"
[[ "$ping_response" == "pong" ]] || fail "Backend health check failed at $BACKEND_URL/ping. Start backend first."
pass "Backend is running"

random_email="tenant_check_$(date +%s)@example.com"
register_body_valid="{\"name\":\"Tenant Check\",\"email\":\"$random_email\",\"password\":\"secret123\",\"tenantId\":\"acme\"}"
register_body_unknown="{\"name\":\"Tenant Check\",\"email\":\"$random_email\",\"password\":\"secret123\",\"tenantId\":\"unknownco\"}"

case_header "API valid tenant resolution"
echo "Expected: POST /register with tenantId=acme returns HTTP 200"
http_code_known="$(curl -s -o /tmp/cardflex_tenant_known.json -w "%{http_code}" -X POST "$BACKEND_URL/register" -H "Content-Type: application/json" -d "$register_body_valid")"
echo "Actual HTTP: $http_code_known"
echo "Actual body: $(cat /tmp/cardflex_tenant_known.json)"
[[ "$http_code_known" == "200" ]] || fail "Register for valid tenant 'acme' failed (HTTP $http_code_known): $(cat /tmp/cardflex_tenant_known.json)"
pass "Valid tenant resolves and registration works"

case_header "API unknown tenant rejection"
echo "Expected: POST /register with tenantId=unknownco returns HTTP 404 with tenant not found"
http_code_unknown="$(curl -s -o /tmp/cardflex_tenant_unknown.json -w "%{http_code}" -X POST "$BACKEND_URL/register" -H "Content-Type: application/json" -d "$register_body_unknown")"
echo "Actual HTTP: $http_code_unknown"
echo "Actual body: $(cat /tmp/cardflex_tenant_unknown.json)"
[[ "$http_code_unknown" == "404" ]] || fail "Unknown tenant did not return 404 (HTTP $http_code_unknown): $(cat /tmp/cardflex_tenant_unknown.json)"
grep -q 'tenant not found' /tmp/cardflex_tenant_unknown.json || fail "Unknown tenant response did not contain 'tenant not found'"
pass "Unknown tenant is rejected correctly"

echo
echo "All tenant checks passed."
