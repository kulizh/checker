#!/bin/bash

if [ ! -f checker.pid ]; then
    echo "checker.pid not found."
    exit 1
fi

PID=$(cat checker.pid)

if kill -0 $PID > /dev/null 2>&1; then
    kill $PID
    echo "Kill sent to PID $PID"
    rm checker.pid
else
    echo "PID $PID not found."
    rm checker.pid
fi