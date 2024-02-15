#!/bin/sh
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP
set -eu

SCRIPT_DIR=$(dirname "$0")
. $SCRIPT_DIR/.env

docker compose --profile nats exec -it -u nats -w /home/nats nats-box /bin/sh /build/configure-nats-box.sh
