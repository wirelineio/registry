#!/bin/bash

btcctl -C ./btcctl.conf --rpcserver=127.0.0.1:8334 --rpccert=../btcd/rpc/rpc.cert --simnet generate 3
