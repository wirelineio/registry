#!/bin/bash

btcctl -C ./btcctl.conf --rpcserver=localhost:8334 --rpccert=../btcd/rpc/rpc.cert --simnet generate 3
