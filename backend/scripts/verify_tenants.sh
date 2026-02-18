#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BACKEND_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"
ENV_FILE="$BACKEND_DIR/.env"

if [[ -f "$ENV_FILE" ]]; then
  set -a
  source "$ENV_FILE"
  set +a
fi

: "${MONGO_URI:?MONGO_URI is required (set in backend/.env or env var)}"
: "${MONGO_DB:?MONGO_DB is required (set in backend/.env or env var)}"

BACKEND_URL="${BACKEND_URL:-http://localhost:${PORT:-8080}}"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "ERROR: '$1' is required but not installed." >&2
    exit 1
  fi
}

require_cmd curl
require_cmd mongosh

pass() { echo "PASS: $1"; }
fail() { echo "FAIL: $1"; exit 1; }
case_header() { echo; echo "CASE: $1"; }

mongo_eval() {
  local js="$1"
  mongosh "$MONGO_URI" --quiet --eval "db = db.getSiblingDB('$MONGO_DB'); $js"
}

echo "== Tenant Verification =="
echo "Backend URL: $BACKEND_URL"
echo "Mongo DB: $MONGO_DB"

case_header "Health check"
echo "Expected: /ping returns pong"
ping_response="$(curl -fsS "$BACKEND_URL/ping" || true)"
echo "Actual: $ping_response"
[[ "$ping_response" == "pong" ]] || fail "Backend health check failed at $BACKEND_URL/ping. Start backend first."
pass "Backend is running"

case_header "Tenant index check"
echo "Expected: unique index named unique_company_code on companyCode"
index_json="$(mongo_eval "const idx = db.tenants.getIndexes().find(i => i.name === 'unique_company_code'); print(JSON.stringify(idx || {}));")"
echo "Actual: $index_json"
echo "$index_json" | grep -q '"name":"unique_company_code"' || fail "Missing index unique_company_code"
echo "$index_json" | grep -q '"unique":true' || fail "Index unique_company_code is not unique"
pass "Unique index exists on companyCode"

case_header "Seed tenants check"
echo "Expected: at least 3 tenants with companyCode in [acme, nova, prime]"
seed_count="$(mongo_eval "const codes = ['acme','nova','prime']; print(db.tenants.countDocuments({ companyCode: { \$in: codes } }));")"
echo "Actual count: $seed_count"
[[ "$seed_count" -ge 3 ]] || fail "Seed tenants missing. Found $seed_count of expected 3"
pass "Seed tenants are present"

case_header "Tenant schema required fields check"
echo "Expected: 0 tenant docs missing name/companyCode"
invalid_count="$(mongo_eval "print(db.tenants.countDocuments({ \$or: [ { name: { \$exists: false } }, { companyCode: { \$exists: false } } ] }));")"
echo "Actual missing count: $invalid_count"
[[ "$invalid_count" -eq 0 ]] || fail "Found tenants missing required fields (name/companyCode)"
pass "Tenant documents contain required fields"

random_email="tenant_check_$(date +%s)@example.com"
register_body="{\"name\":\"Tenant Check\",\"email\":\"$random_email\",\"password\":\"secret123\"}"

case_header "API valid tenant resolution"
echo "Expected: POST /register?company=acme returns HTTP 201"
http_code_known="$(curl -s -o /tmp/cardflex_tenant_known.json -w "%{http_code}" -X POST "$BACKEND_URL/register?company=acme" -H "Content-Type: application/json" -d "$register_body")"
echo "Actual HTTP: $http_code_known"
echo "Actual body: $(cat /tmp/cardflex_tenant_known.json)"
[[ "$http_code_known" == "201" ]] || fail "Register for valid tenant 'acme' failed (HTTP $http_code_known): $(cat /tmp/cardflex_tenant_known.json)"
pass "Valid tenant resolves and registration works"

case_header "API unknown tenant rejection"
echo "Expected: POST /register?company=unknownco returns HTTP 404 with tenant not found"
http_code_unknown="$(curl -s -o /tmp/cardflex_tenant_unknown.json -w "%{http_code}" -X POST "$BACKEND_URL/register?company=unknownco" -H "Content-Type: application/json" -d "$register_body")"
echo "Actual HTTP: $http_code_unknown"
echo "Actual body: $(cat /tmp/cardflex_tenant_unknown.json)"
[[ "$http_code_unknown" == "404" ]] || fail "Unknown tenant did not return 404 (HTTP $http_code_unknown): $(cat /tmp/cardflex_tenant_unknown.json)"
grep -q 'tenant not found' /tmp/cardflex_tenant_unknown.json || fail "Unknown tenant response did not contain 'tenant not found'"
pass "Unknown tenant is rejected correctly"

echo
echo "All tenant checks passed."
