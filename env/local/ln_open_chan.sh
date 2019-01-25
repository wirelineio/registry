#!/bin/bash

RPC_PORT=""
LNDIR=""

if [ "$1" -eq "1" ]
then 
    RPC_PORT="10009"
    LNDIR="lnd1"
fi 

if [ "$1" -eq "2" ]
then 
    RPC_PORT="10109"
    LNDIR="lnd2"
fi 

lncli   --network simnet \
        --rpcserver 127.0.0.1:$RPC_PORT \
        --tlscertpath=$LNDIR/.lnd/tls.cert \
        --macaroonpath=$LNDIR/.lnd/data/chain/bitcoin/simnet/admin.macaroon \
        openchannel --node_key=$2 --local_amt=1000000
