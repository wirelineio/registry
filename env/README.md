# Environment Setup

## Install

Installing lnd and btcd.

```bash
go get -d github.com/lightningnetwork/lnd
cd $GOPATH/src/github.com/lightningnetwork/lnd
make && make install
make btcd
```

Installing btcctl.

```bash
cd $GOPATH/src/github.com/btcsuite/btcd
GO111MODULE=on go install -v . ./cmd/...
```

Generate certificates for btcd, if required. Note: You don't need this for `--local` mode.

```bash
cd $GOPATH/src/github.com/wirelineio/registry/env/dev/btcd
gencerts --host="*" --directory=./rpc --force
# Use -H to add additional host names (e.g. on EC2).
```

## Initialize Services

Run btcd.

```bash
cd btcd
./run.sh
```

Run lnd #1 and #2.

```bash
cd lnd1
./run.sh
```

```bash
cd lnd2
./run.sh
```

Use `--local` for local instances or `--dev` to connect to the instances on EC2.

Create wallet for lnd #1 and lnd #2.

```bash
./wallet_create.sh --local 1
Input wallet password:
Confirm wallet password:

Do you have an existing cipher seed mnemonic you want to use? (Enter y/n): n

Your cipher seed can optionally be encrypted.
Input your passphrase if you wish to encrypt it (or press enter to proceed without a cipher seed passphrase):

Generating fresh cipher seed...

lnd successfully initialized!
```

```bash
./wallet_create.sh --local 2
```

Create new address for lnd #1.

```bash
./wallet_new_address.sh --local 1
{
    "address": "rc3AKE1i7kdNCCZD2MWQTrJ3wMXXNpSsS2"
}
```

Restart btcd with mining enabled to pay the above address.

```bash
cd btcd
MINING_ADDRESS=rc3AKE1i7kdNCCZD2MWQTrJ3wMXXNpSsS2 ./run.sh
```

Generate 400 blocks (we need at least "100 >=" blocks because of coinbase block maturity and "300 ~=" in order to activate segwit):

```bash
./btcd_init.sh --local
```

Check wallet balance. You might have to restart lnd and unlock the wallet using ```./wallet_unlock.sh 1``` for the balance to show correctly.

```bash
./wallet_balance.sh --local 1
{
    "total_balance": "1505000000000",
    "confirmed_balance": "1505000000000",
    "unconfirmed_balance": "0"
}
```

## Connect LND instances

Get identity pubkey of lnd #2.

```bash
./ln_get_info.sh --local 2
{
    "version": "0.5.1-beta commit=v0.5.1-beta-509-gecd5541d55fbf1218662d0d95b220411f25415ed",
    "identity_pubkey": "0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2",
    "alias": "0393cc3fa60704da2568",
    "num_pending_channels": 0,
    "num_active_channels": 0,
    "num_inactive_channels": 0,
    "num_peers": 0,
    "block_height": 400,
    "block_hash": "08f40478fbdccf8c11d4fa67f820ba8b7e94fbc37292ce8dda7b5e6fa62c5c7e",
    "best_header_timestamp": 1548422737,
    "synced_to_chain": true,
    "testnet": false,
    "chains": [
        {
            "chain": "bitcoin",
            "network": "simnet"
        }
    ],
    "uris": [
        "0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2@127.0.0.1:9835"
    ]
}
```

Connect lnd #1 to lnd #2 using the URI in the output.

```bash
./ln_connect.sh --local 1 0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2@127.0.0.1:9835

./ln_list_peer.sh --local 1
{
    "peers": [
        {
            "pub_key": "0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2",
            "address": "127.0.0.1:9835",
            "bytes_sent": "137",
            "bytes_recv": "137",
            "sat_sent": "0",
            "sat_recv": "0",
            "inbound": false,
            "ping_time": "0"
        }
    ]
}
```

## Open Channel

Open channel from lnd #1 to #2.

```bash
./ln_open_chan.sh --local 1 0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2
{
    "funding_txid": "e32aa05b200b77e29123fd5df76c2867f078a59a03622318118e94b75d8b0aa3"
}
```

Include funding transaction in a block thereby opening the channel.

```bash
cd btcd
./btcd_blocks.sh --local
[
  "351676ac207bf57e6ff74f25bab31f5ff448ddd45cf264c8aa44cf34dc1d53e6",
  "0dfcad26e3feb21ea47cbdce01b6ae8973024d23cbf6d82e6866c6cd34b6c413",
  "68c1707ca2edcf50569bb50af671ad9263fefc6b2d098ddb8e5ed149ba8b6955"
]
```

List channels.

```bash
./ln_list_chan.sh --local 1
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2",
            "channel_point": "e32aa05b200b77e29123fd5df76c2867f078a59a03622318118e94b75d8b0aa3:0",
            "chan_id": "440904162803712",
            "capacity": "1000000",
            "local_balance": "990950",
            "remote_balance": "0",
            "commit_fee": "9050",
            "commit_weight": "600",
            "fee_per_kw": "12500",
            "unsettled_balance": "0",
            "total_satoshis_sent": "0",
            "total_satoshis_received": "0",
            "num_updates": "0",
            "pending_htlcs": [
            ],
            "csv_delay": 144,
            "private": false,
            "initiator": true
        }
    ]
}
```

## Payment

Add invoice on lnd #2.

```bash
./ln_add_invoice.sh --local 2
{
    "r_hash": "d5d3afc2aba122b2319547b0e4b70f95043fd6a3fc56dd75eb34521afa2ec954",
    "pay_req": "lnsb100u1pwykxe8pp56hf6ls4t5y3tyvv4g7cwfdc0j5zrl44rl3td6a0tx3fp473we92qdqqcqzyshx23aw889jufgpye4nrzcaqtqq5ey4njq84wwws0xe2hjjeqnqlzkakyya76j2qmcw5qj0xmj9z2rarlx8tdgwtec9mxapjvst7gx6cpzpmclj",
    "add_index": 1
}
```

Pay invoice from lnd #1.

```bash
./ln_pay_invoice.sh --local 1 lnsb100u1pwykxe8pp56hf6ls4t5y3tyvv4g7cwfdc0j5zrl44rl3td6a0tx3fp473we92qdqqcqzyshx23aw889jufgpye4nrzcaqtqq5ey4njq84wwws0xe2hjjeqnqlzkakyya76j2qmcw5qj0xmj9z2rarlx8tdgwtec9mxapjvst7gx6cpzpmclj
Description:
Amount (in satoshis): 10000
Destination: 0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2
Confirm payment (yes/no): yes
{
    "payment_error": "",
    "payment_preimage": "0601d96707470f643da39db88d2954cbdea48eddfa66d3a1468bb407043898b2",
    "payment_route": {
        "total_time_lock": 547,
        "total_amt": 10000,
        "hops": [
            {
                "chan_id": 440904162803712,
                "chan_capacity": 1000000,
                "amt_to_forward": 10000,
                "expiry": 547,
                "amt_to_forward_msat": 10000000,
                "pub_key": "0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2"
            }
        ],
        "total_amt_msat": 10000000
    }
}
```

Check channel balances.

```bash
./ln_chan_balance.sh --local 1
{
    "balance": "980950",
    "pending_open_balance": "0"
}

./ln_chan_balance.sh --local 2
{
    "balance": "10000",
    "pending_open_balance": "0"
}
```

## Close Channel

Get the funding tx outpoint.

```bash
./ln_list_chan.sh --local 1
{
    "channels": [
        {
            "active": true,
            "remote_pubkey": "0393cc3fa60704da2568216d0609938e1cf1820cf0e53738fd3b66b02ec40fbdb2",
            "channel_point": "e32aa05b200b77e29123fd5df76c2867f078a59a03622318118e94b75d8b0aa3:0",
            "chan_id": "440904162803712",
            "capacity": "1000000",
            "local_balance": "980950",
            "remote_balance": "10000",
            "commit_fee": "9050",
            "commit_weight": "724",
            "fee_per_kw": "12500",
            "unsettled_balance": "0",
            "total_satoshis_sent": "10000",
            "total_satoshis_received": "0",
            "num_updates": "2",
            "pending_htlcs": [
            ],
            "csv_delay": 144,
            "private": false,
            "initiator": true
        }
    ]
}
```

Note: "channel_point": "e32aa05b200b77e29123fd5df76c2867f078a59a03622318118e94b75d8b0aa3:0"

```bash
./ln_close_chan.sh --local 1 e32aa05b200b77e29123fd5df76c2867f078a59a03622318118e94b75d8b0aa3 0
{
    "closing_txid": "f526bc7c4a603450d186e378819e00d301279ac0d098c9eb938e86e220fef4f8"
}
```

Include close transaction in a block thereby closing the channel.

```bash
cd btcd
./btcd_blocks.sh --local
[
  "61a9fe6496bce13579dc9056a8eeb6291487dfa14db07c9a6b4b3ba90a43d915",
  "254f107099bcb5fcb4fdd28f94505f60b5595a05dc956ff5114b4b9262bf07a7",
  "26e1e1fcc06d619de503b06b00eb0770d08edf67667745175ce23456cf1af095"
]
```

Check wallet balances.

```bash
./wallet_balance.sh --local 1
{
    "total_balance": "1534999972163",
    "confirmed_balance": "1534999972163",
    "unconfirmed_balance": "0"
}

./wallet_balance.sh --local 2
{
    "total_balance": "10000",
    "confirmed_balance": "10000",
    "unconfirmed_balance": "0"
}

```
