#!/bin/bash

btcctl -C ./btcctl.conf --rpcserver=localhost:8334 --rpccert=../btcd/rpc/rpc.cert --simnet generate 400

btcctl -C ./btcctl.conf --rpcserver=localhost:8334 --rpccert=../btcd/rpc/rpc.cert --simnet  getblockchaininfo | grep -A 1 segwit
