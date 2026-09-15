#!/bin/bash
set -e

CONFIG_FILE="/etc/postfix/opengovmail.yaml"

# Fail fast if the config file is missing or unreadable.
# Never fall back to a placeholder domain — that would leave Postfix
# running for the wrong domain.
if [ ! -f "$CONFIG_FILE" ]; then
    echo "ERROR: Config file not found: $CONFIG_FILE" >&2
    echo "ERROR: Please ensure opengovmail.yaml is mounted correctly. Refusing to start." >&2
    exit 1
fi

if [ ! -r "$CONFIG_FILE" ]; then
    echo "ERROR: Config file not readable: $CONFIG_FILE" >&2
    exit 1
fi

# Extract primary (first) domain from the domains list using awk
MAIL_DOMAIN=$(awk '/^[[:space:]]*-[[:space:]]*domain:/ {sub(/^[[:space:]]*-[[:space:]]*domain:[[:space:]]*/, ""); print; exit}' "$CONFIG_FILE")
# Trim surrounding whitespace (awk leaves empty/whitespace-only values behind)
MAIL_DOMAIN="$(echo "$MAIL_DOMAIN" | sed 's/^[[:space:]]*//; s/[[:space:]]*$//')"

# Fail fast if extraction failed — never use a placeholder domain
if [ -z "$MAIL_DOMAIN" ] || [ "$MAIL_DOMAIN" = "null" ]; then
    echo "ERROR: Could not extract domain from $CONFIG_FILE" >&2
    echo "ERROR: Please ensure opengovmail.yaml has domains configured. Refusing to start." >&2
    exit 1
fi

# -------------------------------
# Environment variables
# -------------------------------
export MAIL_DOMAIN
MAIL_HOSTNAME=${MAIL_HOSTNAME:-mail.$MAIL_DOMAIN}
RELAYHOST=${RELAYHOST:-}

echo "Using domain: $MAIL_DOMAIN"
echo "$MAIL_DOMAIN" > /etc/mailname

# -------------------------------
# Check SQLite database
# -------------------------------
DB_PATH="/app/data/databases/shared.db"

echo "=== Checking SQLite database ==="
if [ -f "$DB_PATH" ]; then
    echo "✓ SQLite database found at $DB_PATH"

    # Ensure domain exists in database
    sqlite3 "$DB_PATH" "INSERT OR IGNORE INTO domains (domain, enabled) VALUES ('${MAIL_DOMAIN}', 1);" 2>/dev/null || echo "Note: Could not insert domain (may already exist)"

    # Set proper permissions
    chmod 644 "$DB_PATH"
else
    echo "⚠ Warning: SQLite database not found at $DB_PATH"
    echo "  Database should be created by raven"
    echo "  Postfix will start but mail delivery may fail until database is available"
fi

echo "=== Database setup completed ==="

# -------------------------------
# Fix for DNS resolution in chroot
# -------------------------------
mkdir -p /var/spool/postfix/etc
cp /etc/host.conf /etc/resolv.conf /etc/services /var/spool/postfix/etc/
chmod 644 /var/spool/postfix/etc/*

# -------------------------------
# Verify configuration
# -------------------------------
echo "=== Verifying Postfix configuration ==="
postconf virtual_mailbox_domains
postconf virtual_mailbox_maps
postconf virtual_mailbox_base
postconf virtual_transport

# -------------------------------
# Start Postfix
# -------------------------------
echo "=== Starting Postfix ==="
service postfix start

# Keep container running
sleep infinity