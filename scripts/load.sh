#!/usr/bin/env bash
# Continuously write random key/values into the clone to watch eviction kick in.
# KeysLimit=100, allkeys-random, EvictionRatio=0.40 -> store caps ~100, drops 40 at a time.
#
# Usage: ./scripts/load.sh [port] [delay_seconds]
#   port          defaults to 7379
#   delay_seconds pause between writes (default 0.05). Use 0 for max speed.

set -euo pipefail

PORT="${1:-7379}"
DELAY="${2:-0.05}"

echo "Writing random keys to localhost:$PORT (Ctrl-C to stop)..."
i=0
while true; do
  key="key:$RANDOM$RANDOM"
  val="val:$RANDOM"
  redis-cli -p "$PORT" set "$key" "$val" >/dev/null
  i=$((i + 1))
  if (( i % 50 == 0 )); then
    echo "wrote $i keys"
  fi
  if [[ "$DELAY" != "0" ]]; then
    sleep "$DELAY"
  fi
done
