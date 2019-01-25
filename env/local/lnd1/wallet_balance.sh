#!/bin/bash

lncli --network simnet --rpcserver 127.0.0.1:10009 --tlscertpath=./.lnd/tls.cert --macaroonpath=./.lnd/data/chain/bitcoin/simnet/admin.macaroon walletbalance
