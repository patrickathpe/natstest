#!/bin/sh
# (C) Copyright 2024 Hewlett Packard Enterprise Development LP
set -eu

SCRIPT_DIR=$(dirname "$0")
. $SCRIPT_DIR/.env

nsc add operator $NATS_OPERATOR
nsc edit operator --service-url=$NATS_URL
nsc add account --name=SYS
nsc edit operator --system-account=SYS
nsc add account $NATS_ACCOUNT
nsc edit account \
	--name $NATS_ACCOUNT \
	--js-mem-storage=-1 \
	--js-disk-storage=-1 \
	--js-streams=-1 \
	--js-consumer=-1
nsc add user \
	$NATS_USER \
	--allow-pub "$NATS_RTM_EVENTS_SUBJECT" \
	--allow-sub "_INBOX.>"
nsc add user \
	$TEST_USER \
	--allow-pub "$NATS_RTM_EVENTS_SUBJECT" \
	--allow-sub "$NATS_RTM_EVENTS_SUBJECT" \
	--allow-sub "_INBOX.>" \
	--allow-pub "\$JS.API.INFO" \
	--allow-pub "\$JS.API.STREAM.CREATE.$NATS_RTM_STREAM" \
	--allow-pub "\$JS.API.STREAM.DELETE.$NATS_RTM_STREAM" \
	--allow-pub "\$JS.API.STREAM.INFO.$NATS_RTM_STREAM" \
	--allow-pub "\$JS.API.STREAM.UPDATE.$NATS_RTM_STREAM" \
	--allow-pub "\$JS.API.STREAM.LIST" \
	--allow-pub "\$JS.API.STREAM.NAMES" \
	--allow-pub "\$JS.ACK.$NATS_RTM_STREAM.>" \
	--allow-pub "\$JS.API.CONSUMER.NAMES.$NATS_RTM_STREAM" \
	--allow-pub "\$JS.API.CONSUMER.INFO.$NATS_RTM_STREAM.>" \
	--allow-pub "\$JS.API.CONSUMER.CREATE.$NATS_RTM_STREAM.>" \
	--allow-pub "\$JS.API.CONSUMER.DELETE.$NATS_RTM_STREAM.>" \
	--allow-pub "\$JS.API.CONSUMER.DURABLE.CREATE.$NATS_RTM_STREAM.>" \
	--allow-pub "\$JS.API.CONSUMER.UPDATE.$NATS_RTM_STREAM.>" \
	--allow-pub "\$JS.API.CONSUMER.MSG.NEXT.$NATS_RTM_STREAM.>"
