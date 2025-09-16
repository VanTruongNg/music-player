#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(git rev-parse --show-toplevel 2>/dev/null || pwd)"
PRIV_DIR="${ROOT_DIR}/infra/jwt/private"
PUB_DIR="${ROOT_DIR}/infra/jwt/public"
JWKS_DIR="${ROOT_DIR}/infra/jwt/jwks"
UPDATE_JWKS="${ROOT_DIR}/infra/scripts/update_jwks.sh"

mkdir -p "$PRIV_DIR" "$PUB_DIR" "$JWKS_DIR"

#Timestamp kid + random suffix to avoid collisions
timestamp=$(date -u +%Y%m%dT%H%M%SZ)
suffix=$(head -c 8 /dev/urandom | base32 | tr -d '=' | tr 'A-Z' 'a-z' | cut -c1-5)
kid="${timestamp}-${suffix}"

priv="${PRIV_DIR}/ed25519-${kid}.pem"
pub="${PUB_DIR}/ed25519-${kid}.pub.pem"

echo ">>> Generating Ed25519 keypair kid=$kid"
openssl genpkey -algorithm ED25519 -out "$priv"
openssl pkey -in "$priv" -pubout -out "$pub"

if command -v chmod >/dev/null 2>&1; then
    chmod 600 "$priv" || true
    chmod 644 "$pub" || true
fi

#Fingerprint
fp_pub="$(openssl pkey -pubin -in "$pub" -outform DER | openssl sha256 | awk '{print $2}')"
sha_priv="$( (command -v sha256sum >/dev/null 2>&1 && sha256sum "$priv" | awk '{print $1}') || (shasum -a 256 "$priv" | awk '{print $1}') )"
sha_pub="$( (command -v sha256sum >/dev/null 2>&1 && sha256sum "$pub"  | awk '{print $1}') || (shasum -a 256 "$pub"  | awk '{print $1}') )"

echo "fingerprint_pub(SPKI_SHA256)=${fp_pub}"
echo "sha256_private=${sha_priv}"
echo "sha256_public=${sha_pub}"

# --- Update JWKS from all public keys ---
if [[ -x "$UPDATE_JWKS" ]]; then
    "$UPDATE_JWKS"
else
    echo "WARN: $UPDATE_JWKS not found or not executable. Skipping JWKS update."
fi

echo
echo ">>> Done."
echo "Private: $priv"
echo "Public : $pub"
echo "JWKS   : ${JWKS_DIR}/jwks.json (if updated)"
echo
echo "Set ENV for Auth:"
echo "  JWT_ACCESS_PRIVATE_KEY_FILE=/etc/keys/jwt/private/ed25519-${kid}.pem"
echo "  JWT_ACCESS_KID=${kid}"