#!/bin/bash

set -e

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 user@host /remote/path"
  exit 1
fi

make build-linux

REMOTE_HOST="$1"
REMOTE_PATH="$2"

FILES=("checker" "start.sh" "stop.sh" "readme.md" ".env")

CONFIG_FILE="configs/domains.example.json"
TEMP_CONFIG="domains.json"

CLEANUP_REQUIRED=0

# Проверяем, существует ли локальный domains.json
if [ ! -f "$TEMP_CONFIG" ]; then
  echo "No local domains.json found → creating temp from example"
  cp "$CONFIG_FILE" "$TEMP_CONFIG"
  CLEANUP_REQUIRED=1
else
  echo "Using existing domains.json (no temp copy needed)"
fi

echo "Creating dir $REMOTE_HOST:$REMOTE_PATH..."
ssh "$REMOTE_HOST" "mkdir -p '$REMOTE_PATH'"

echo "Sending files..."

UPLOAD_FILES=("${FILES[@]}")

if [ -f "$TEMP_CONFIG" ]; then
  UPLOAD_FILES+=("$TEMP_CONFIG")
fi

scp "${UPLOAD_FILES[@]}" "$REMOTE_HOST:$REMOTE_PATH/"

if [ "$CLEANUP_REQUIRED" -eq 1 ]; then
  echo "Cleaning up temp domains.json"
  rm -f "$TEMP_CONFIG"
fi

echo "Done."