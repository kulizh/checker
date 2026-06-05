#!/bin/bash
# Deploy script — builds binary and copies files to remote server.
# Usage: ./deploy.sh user@host /remote/path
# Variables (override as needed):
#   CONFIG     — path to local config (default: domains.json)
#   REMOTE_DIR — same as $2, can be set via env

set -e

if [ "$#" -lt 2 ] && [ -z "$REMOTE_DIR" ]; then
  echo "Usage: $0 user@host /remote/path"
  echo "       REMOTE_DIR=/remote/path $0 user@host"
  exit 1
fi

REMOTE_HOST="${1}"
REMOTE_PATH="${2:-$REMOTE_DIR}"
CONFIG="${CONFIG:-domains.json}"

# Build for linux
make build-linux

# Ensure config exists
if [ ! -f "$CONFIG" ]; then
  if [ -f "configs/domains.example.json" ]; then
    echo "Copying configs/domains.example.json → $CONFIG"
    cp configs/domains.example.json "$CONFIG"
  else
    echo "ERROR: no $CONFIG found"
    exit 1
  fi
fi

echo "Deploying to $REMOTE_HOST:$REMOTE_PATH"
ssh "$REMOTE_HOST" "mkdir -p '$REMOTE_PATH'"

scp \
  checker \
  start.sh \
  stop.sh \
  .env \
  "$CONFIG" \
  "$REMOTE_HOST:$REMOTE_PATH/"

echo "Done. On the server:"
echo "  cd $REMOTE_PATH"
echo "  # edit .env with your tokens"
echo "  ./start.sh"