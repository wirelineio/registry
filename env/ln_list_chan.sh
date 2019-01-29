#!/bin/bash

source ./vars.sh 

lncli   --network simnet \
        --rpcserver $RPC_HOST:$RPC_PORT \
        --tlscertpath=$BASE_DIR/$LNDIR/.lnd/tls.cert \
        --macaroonpath=$BASE_DIR/$LNDIR/.lnd/data/chain/bitcoin/simnet/admin.macaroon \
        listchannels