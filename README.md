# Wirechain

Wirechain is a custom blockchain built using the Cosmos SDK. It currently houses custom modules and toolchain for Settlements 2.0 PoC.

## Setup & install

Install golang 1.11+ for your platform.

```
$ go version
go version go1.11 darwin/amd64
```

Clone the repo.

```
$ mkdir -p $GOPATH/src/github.com/wirelineio
$ cd $GOPATH/src/github.com/wirelineio
$ git clone git@github.com:wirelineio/wirechain.git
```

Initialize dep and install dependencies.

```
$ cd wirechain
$ make get_tools && make get_vendor_deps
```

Install the app into your `$GOBIN`.

```
$ make install
```

Test that the commands are installed.

```
$ wirechaind help
$ wirecli help
```

## Initialize & start blockchain

Delete existing blockchain and config.

```
$ rm -rf ~/.wirechaind/ ~/.wirecli
```

Initialize the chain.

```
$ wirechaind init --chain-id wireline
```

Setup accounts/keys for Alice and Bob. Enter a passphrase for the key when prompted.

```
$ wirecli keys add alice

$ wirecli keys add bob
```

Add initial funds to the accounts.

```
$ wirechaind add-genesis-account $(wirecli keys show alice --address) 1000wire
$ wirechaind add-genesis-account $(wirecli keys show bob --address) 10wire

```

Start the blockchain. You should see blocks being created every few seconds.


```
$ wirechaind start
```

Run the following commands in another terminal.

## HTLC module

### Redeem Flow

Alice creates a HTLC that will pay 25wire to Bob if he can produce the preimage of the SHA256 hash `6815f3c300383519de8e437497e2c3e97852fe8d717a5419d5aafb00cb43c494` within 20 blocks.

```
$ wirecli tx htlc add 25wire 6815f3c300383519de8e437497e2c3e97852fe8d717a5419d5aafb00cb43c494 20 $(wirecli keys show bob --address) --from=$(wirecli keys show alice --address) --chain-id=wireline
```

Check balances. Note that Alice's balance has decreased by 25wire, now locked in the HTLC.

```
$ wirecli query account $(wirecli keys show alice --address) --indent --chain-id=wireline
$ wirecli query account $(wirecli keys show bob --address) --indent --chain-id=wireline
```

Bob redeems the HTLC by presenting the preimage `mango` within 20 blocks.

```
$ wirecli tx htlc redeem mango --from=$(wirecli keys show bob --address) --chain-id=wireline
```

Check balances. Note that Bob's balance has increased by 25wire.

```
$ wirecli query account $(wirecli keys show alice --address) --indent --chain-id=wireline
$ wirecli query account $(wirecli keys show bob --address) --indent --chain-id=wireline
```


### Timeout Flow

Alice creates a HTLC that will pay 25wire to Bob if he can produce the preimage of the SHA256 hash `CB4595E84361629F0BB6A3EAA79220E553A8360191433538FF0E0C41F44E30DB ` within 20 blocks.

```
$ wirecli tx htlc add 25wire CB4595E84361629F0BB6A3EAA79220E553A8360191433538FF0E0C41F44E30DB 20 $(wirecli keys show bob --address) --from=$(wirecli keys show alice --address) --chain-id=wireline
```

Alice claims a timeout after waiting for 20 blocks (also try running this before the timeout and see what happens).

```
$ wirecli tx htlc fail CB4595E84361629F0BB6A3EAA79220E553A8360191433538FF0E0C41F44E30DB --from=$(wirecli keys show alice --address) --chain-id=wireline
```

Check balances. Note that Alice has been refunded 25wire.

```
$ wirecli query account $(wirecli keys show alice --address) --indent --chain-id=wireline
$ wirecli query account $(wirecli keys show bob --address) --indent --chain-id=wireline
```


## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
