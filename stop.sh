#!/bin/bash

PID_FILE="checker.pid"

if [ ! -f "$PID_FILE" ]; then
  echo "PID file not found ($PID_FILE)"
  exit 1
fi

PID=$(cat "$PID_FILE")

if kill -0 "$PID" 2>/dev/null; then
  kill "$PID"
  echo "Stopped (PID $PID)"
else
  echo "Process $PID not running"
fi

rm -f "$PID_FILE"