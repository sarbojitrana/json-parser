#!/bin/bash
set -euo pipefail

cat > tern.conf <<EOF
[database]
host = ${PARSER_DATABASE_HOST}
port = ${PARSER_DATABASE_PORT}
database = ${PARSER_DATABASE_NAME}
user = ${PARSER_DATABASE_USER}
password = ${PARSER_DATABASE_PASSWORD}
sslmode = ${PARSER_DATABASE_SSL_MODE}

[migration]
path = internal/db/migrations
EOF

exec tern migrate --config tern.conf --migrations ./internal/db/migrations
