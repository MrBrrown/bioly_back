#!/usr/bin/env bash
set -e

CERT_DIR="./nginx/certs"
CERT_KEY="${CERT_DIR}/privkey.pem"
CERT_FULL="${CERT_DIR}/fullchain.pem"
DAYS=365
DOMAIN="bioly.localhost"
ALT_DOMAINS="DNS:bioly.localhost,DNS:localhost,IP:127.0.0.1"

mkdir -p "${CERT_DIR}"

if [[ -f "${CERT_KEY}" || -f "${CERT_FULL}" ]]; then
  echo "Certificates already exist in ${CERT_DIR}"
  read -p "Overwrite existing certificates? [y/N]: " confirm
  if [[ "${confirm}" != "y" && "${confirm}" != "Y" ]]; then
    echo "Aborted."
    exit 0
  fi
  rm -f "${CERT_KEY}" "${CERT_FULL}"
fi

echo "🔧 Generating self-signed certificate for ${DOMAIN}..."
openssl req -x509 -nodes -newkey rsa:2048 \
  -subj "/C=RU/ST=Dev/L=Local/O=Bioly/CN=${DOMAIN}" \
  -addext "subjectAltName=${ALT_DOMAINS}" \
  -keyout "${CERT_KEY}" \
  -out "${CERT_FULL}" \
  -days "${DAYS}" >/dev/null 2>&1

echo "Certificates generated successfully:"
echo "   - ${CERT_FULL}"
echo "   - ${CERT_KEY}"
echo
echo "You can now run Nginx with HTTPS support."