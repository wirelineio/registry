#!/bin/bash

PARAMS=""

if [[ -n "$MINING_ADDRESS" ]]; then
    PARAMS="$PARAMS --miningaddr=$MINING_ADDRESS"
fi

btcd -C ./btcd.conf --rpccert=./rpc/rpc.cert --rpckey=./rpc/rpc.key $PARAMS

