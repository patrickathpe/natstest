#!/bin/sh
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP
set -eu

SCRIPT_DIR=$(dirname "$0")
. $SCRIPT_DIR/.env

nsc describe account
nsc push -A -u $NATS_URL
