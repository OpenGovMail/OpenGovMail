#!/bin/bash

# ============================================
#  Print DNS SRV records for autodiscovery
# ============================================

set -euo pipefail

readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
readonly CONFIG_FILE="$(cd "${SCRIPT_DIR}/../../conf" && pwd)/opengovmail.yaml"

if [[ ! -f "${CONFIG_FILE}" ]]; then
	echo "ERROR: Configuration file not found: ${CONFIG_FILE}" >&2
	exit 1
fi

mapfile -t DOMAINS < <(
	grep -E '^[[:space:]]*-[[:space:]]*domain:' "${CONFIG_FILE}" |
		sed -E 's/^[[:space:]]*-[[:space:]]*domain:[[:space:]]*//' |
		xargs -n1
)

if [[ ${#DOMAINS[@]} -eq 0 ]]; then
	echo "ERROR: No domains configured in ${CONFIG_FILE}" >&2
	exit 1
fi

for domain in "${DOMAINS[@]}"; do
	if [[ -z "${domain}" ]]; then
		echo "ERROR: Empty domain entry in ${CONFIG_FILE}" >&2
		exit 1
	fi
	if ! [[ "${domain}" =~ ^[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
		echo "ERROR: '${domain}' does not look like a valid domain name" >&2
		exit 1
	fi
done

readonly PRIMARY_DOMAIN="${DOMAINS[0]}"
readonly MAIL_HOST="mail.${PRIMARY_DOMAIN}"

echo
echo "Email client autodiscovery DNS records"
echo "--------------------------------------"
echo "Add the following SRV records to the DNS zone for each domain."
echo "The target host is ${MAIL_HOST}."
echo

for domain in "${DOMAINS[@]}"; do
	printf '_imaps._tcp.%-32s 3600 IN SRV 0 1 993 %s.\n' "${domain}." "${MAIL_HOST}"
	printf '_submission._tcp.%-27s 3600 IN SRV 0 1 587 %s.\n' "${domain}." "${MAIL_HOST}"
	printf '_autodiscover._tcp.%-25s 3600 IN SRV 0 1 443 %s.\n' "${domain}." "${MAIL_HOST}"
done

echo
echo "These records are optional. Autodiscovery also works through the HTTPS"
echo "endpoints on ${PRIMARY_DOMAIN} and ${MAIL_HOST}."
echo
