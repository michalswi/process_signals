#!/usr/bin/env bash
set -e

if [ $# -eq 0 ]; then
    echo "please provide 'true'.."
    exit 1
fi

if [ $1 == "true" ]; then
    while $BOOL; do
        date +"%H:%M:%S"
        echo "sleeping.."
        sleep 5
    done
else
    echo "not 'true'.."
    exit 1
fi        