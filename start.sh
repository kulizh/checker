#!/bin/bash

set -e

APP="./checker"
CONFIG_FILE="domains.json"
EXAMPLE_CONFIG="configs/domains.example.json"
INTERVAL="30s"
RENOTIFY="1h"
PID_FILE="checker.pid"

# Resolve config file
if [ ! -f "$CONFIG_FILE" ] && [ -f "$EXAMPLE_CONFIG" ]; then
  echo "Config not found. Falling back to $EXAMPLE_CONFIG"
  CONFIG_FILE="$EXAMPLE_CONFIG"
elif [ ! -f "$CONFIG_FILE" ]; then
  echo "ERROR: No config found!"
  echo "Expected: $CONFIG_FILE or $EXAMPLE_CONFIG"
  exit 1
fi

# Exit if already running
if [ -f "$PID_FILE" ]; then
  PID=$(cat "$PID_FILE")
  if kill -0 "$PID" 2>/dev/null; then
    echo "Already running (PID $PID)"
    exit 0
  fi
  rm -f "$PID_FILE"
fi

echo "Starting: $APP --config $CONFIG_FILE --interval $INTERVAL --renotify $RENOTIFY"

nohup "$APP" --config "$CONFIG_FILE" --interval "$INTERVAL" --renotify "$RENOTIFY" > checker.log 2>&1 &
PID=$!
echo $PID > "$PID_FILE"
echo "Started (PID $PID)"