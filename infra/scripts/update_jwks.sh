#!/usr/bin/env bash
set -euo pipefail

# --- Paths ---
ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
PUB_DIR="${ROOT_DIR}/infra/jwt/public"
JWKS_DIR="${ROOT_DIR}/infra/jwt/jwks"
OUT="${JWKS_DIR}/jwks.json"

mkdir -p "$JWKS_DIR"

# --- Build JWKS from all public PEMs ---
tmp="$(mktemp)"
echo '{"keys":[' > "$tmp"
first=1

shopt -s nullglob
for pem in "$PUB_DIR"/ed25519-*.pub.pem; do
  base="$(basename "$pem")"
  # kid from filename: ed25519-<kid>.pub.pem
  kid="${base#ed25519-}"
  kid="${kid%.pub.pem}"

  # Convert PKIX public key to DER, then take last 32 bytes (raw Ed25519 'x'), base64url
  der_tmp="$(mktemp)"
  openssl pkey -pubin -in "$pem" -outform DER -out "$der_tmp" >/dev/null 2>&1
  x_b64url="$(tail -c 32 "$der_tmp" | base64 | tr '+/' '-_' | tr -d '=')"
  rm -f "$der_tmp"

  jwk="{\"kty\":\"OKP\",\"crv\":\"Ed25519\",\"use\":\"sig\",\"kid\":\"${kid}\",\"x\":\"${x_b64url}\"}"
  if [[ $first -eq 1 ]]; then
    first=0
    echo -n "$jwk" >> "$tmp"
  else
    echo -n ",$jwk" >> "$tmp"
  fi
done
shopt -u nullglob

echo ']}' >> "$tmp"

# Pretty with jq if present
if command -v jq >/dev/null 2>&1; then
  jq . "$tmp" > "$OUT"
else
  cp "$tmp" "$OUT"
fi
rm -f "$tmp"

echo ">>> JWKS updated: $OUT"