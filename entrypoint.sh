#!/usr/bin/env bash
set -e

if [ -d "/certs" ] && [ "$(ls -A /certs/*.crt 2>/dev/null)" ]; then
  cp /certs/*.crt /usr/local/share/ca-certificates/
  update-ca-certificates
fi

# Bootstrap: create the first superuser if one doesn't exist yet.
# Throwaway credential meant to be rotated via the UI right after first login.
# Runs as picmd (same user as the main process) so it has permission to write pb_data.
if [ -n "$INITIAL_ADMIN_EMAIL" ] && [ -n "$INITIAL_ADMIN_PASSWORD" ]; then
  su-exec picmd:picmd picmd superuser create "$INITIAL_ADMIN_EMAIL" "$INITIAL_ADMIN_PASSWORD" --dir="/picmd/data}" || true
fi

exec su-exec picmd:picmd "$@"