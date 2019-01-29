#!/bin/bash

source ./vars.sh

btcctl -C ./$BASE_DIR/btcd/btcctl.conf --rpcserver=$RPC_HOST:8334 --rpccert=./$BASE_DIR/btcd/rpc/rpc.cert --simnet generate 3
