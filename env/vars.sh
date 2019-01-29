#!/bin/bash

RPC_HOST=""
RPC_PORT=""
BASE_DIR=""
LNDIR=""

if [ -z "$1" ]
then
    echo "Pass --local or --dev as first argument."
    exit
fi

if [ "$1" == "--local" ]
then
    BASE_DIR="local"
    RPC_HOST="127.0.0.1"
fi

if [ "$1" == "--dev" ]
then
    BASE_DIR="dev"
    RPC_HOST="ec2-34-227-79-96.compute-1.amazonaws.com"
fi

if (! [ -z "$2" ]) &&  [ "$2" -eq "1" ]
then
    RPC_PORT="10009"
    LNDIR="lnd1"
fi 

if (! [ -z "$2" ]) && [ "$2" -eq "2" ]
then 
    RPC_PORT="10109"
    LNDIR="lnd2"
fi 