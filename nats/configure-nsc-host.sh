#!/bin/sh
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP
set -eu

SCRIPT_DIR=$(dirname "$0")
. $SCRIPT_DIR/.env

docker compose --profile nats exec -it -u nats -w /home/nats nats-box /bin/sh -c /build/configure-nsc-box.sh
mkdir -p "$SCRIPT_DIR/keys"
cp "$SCRIPT_DIR/nsc/nats/nsc/stores/$NATS_OPERATOR/$NATS_OPERATOR.jwt" "$SCRIPT_DIR/keys/operator.jwt"
cp "$SCRIPT_DIR/nsc/nkeys/creds/$NATS_OPERATOR/$NATS_ACCOUNT/$NATS_USER.creds" "$SCRIPT_DIR/keys/user.creds"
cp "$SCRIPT_DIR/nsc/nkeys/creds/$NATS_OPERATOR/$NATS_ACCOUNT/$TEST_USER.creds" "$SCRIPT_DIR/keys/tester.creds"
