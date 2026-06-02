#!/bin/bash

set -e

CONFIG_FILE="domains.json"
EXAMPLE_CONFIG="configs/domains.example.json"

INTERVAL="15s"

if [ -f "$CONFIG_FILE" ]; then
  echo "Using config: $CONFIG_FILE"
elif [ -f "$EXAMPLE_CONFIG" ]; then
  echo "Config not found. Falling back to $EXAMPLE_CONFIG"
  CONFIG_FILE="$EXAMPLE_CONFIG"
else
  echo "ERROR: No config found!"
  echo "Expected:"
  echo "  - $CONFIG_FILE"
  echo "  - $EXAMPLE_CONFIG"
  exit 1
fi

echo "Starting checker..."
echo "Config: $CONFIG_FILE"
echo "Interval: $INTERVAL"

PID_FILE="checker.pid"

if [ -f "$PID_FILE" ]; then
    PID=$(cat "$PID_FILE")
    if ps -p "$PID" > /dev/null 2>&1; then
        exit 0
    else
        rm "$PID_FILE"
    fi
fi

nohup ./checker > /dev/null 2>&1 &
echo $! > "$PID_FILE"