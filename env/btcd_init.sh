#!/bin/bash

source ./vars.sh

btcctl -C ./$BASE_DIR/btcd/btcctl.conf --rpcserver=$RPC_HOST:8334 --rpccert=./$BASE_DIR/btcd/rpc/rpc.cert --simnet generate 400

btcctl -C ./$BASE_DIR/btcd/btcctl.conf --rpcserver=$RPC_HOST:8334 --rpccert=./$BASE_DIR/btcd/rpc/rpc.cert --simnet  getblockchaininfo | grep -A 1 segwit
