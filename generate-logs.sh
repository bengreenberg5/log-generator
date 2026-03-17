#!/bin/sh
set -e

RATE_KiB=${1:-5}
BYTES_PER_SEC=$((RATE_KiB * 1024))

# Fixed-width padding to make every log line exactly 256 bytes (including newline)
# Prefix: 24 (timestamp) + 1 (space) + 4 (msg ) + 12 (i=0000000001) + 1 (space) = 42 chars, +1 newline = 43 fixed. Pad = 213.
PAD="AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA"

echo "Starting log generator at ~${RATE_KiB}KiB/s"

INDEX=0
while true; do
    BYTES_THIS_SEC=0

    while [ "$BYTES_THIS_SEC" -lt "$BYTES_PER_SEC" ]; do
        INDEX=$((INDEX + 1))
        TIMESTAMP=$(date -u '+%Y-%m-%dT%H:%M:%S.000Z')
        # printf ensures fixed width: 24 char timestamp + space + 4 "msg " + 10-digit index + space + 60 pad + newline = 100 bytes
        printf '%s msg i=%010d %s\n' "$TIMESTAMP" "$INDEX" "$PAD"
        BYTES_THIS_SEC=$((BYTES_THIS_SEC + 256))
    done

    sleep 1
done
