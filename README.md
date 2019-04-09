# Registry

Registry is a custom blockchain built using the Cosmos SDK. It currently houses SDK modules and CLI toolchain for Registry and Payments.

## Getting Started

Install golang 1.11+ for your platform.

```
$ go version
go version go1.11 darwin/amd64
```

Clone the repo.

NOTE: The repo must be created in the specified directory under `GOPATH`.

```
$ mkdir -p $GOPATH/src/github.com/wirelineio
$ cd $GOPATH/src/github.com/wirelineio
$ git clone git@github.com:wirelineio/registry.git
```

Initialize dep and install dependencies.

```
$ cd registry
$ make get_tools && make get_vendor_deps
```

Install the app into your `$GOBIN`.

```
$ make install
```

Test that the commands are installed.

```
$ registryd help
$ regcli help
```

## Initialize and Start Blockchain

Delete existing blockchain and config.

```
$ rm -rf ~/.registryd/ ~/.regcli
```

Initialize the chain.

```
$ registryd init --chain-id wireline
```

Setup the genesis account `root` which can be used to transfer funds to other accounts once the blockchain is running. Enter a passphrase for the key when prompted. Write down the generated mnemonic to restore the private key at a later date.

```
$ regcli keys add root
$ registryd add-genesis-account $(regcli keys show root --address) 1000000wire
```

Start the blockchain. You should see blocks being created every few seconds.

```
$ registryd start --gql-server --gql-playground
```

Check that the Registry is up and running by querying the GQL endpoint in another terminal.

```
$ curl -s -X POST -H "Content-Type: application/json" \
  -d '{ "query": "{ getStatus { version } }" }' http://localhost:9473/query | jq
```

### GQL Server API

The GQL server is controlled using the following `registryd` flags:

* `--gql-server` - Enable GQL server.
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:9473/).
* `--gql-port` - Port to run the GQL server on (default 9473).

See `registryd/x/registry/gql/schema.graphql` for the GQL schema.

## Testnets

### Development

Endpoints

* GQL: https://registry-testnet.dev.wireline.ninja/query
* GQL Playground: https://registry-testnet.dev.wireline.ninja/
* RPC: tcp://registry-testnet.dev.wireline.ninja:26657

### Production

Endpoints

* GQL: https://registry-testnet.wireline.ninja/query
* GQL Playground: https://registry-testnet.wireline.ninja/
* RPC: tcp://registry-testnet.wireline.ninja:26657

Note: The `regcli` command accepts a `--node` flag for the RPC endpoint.

## Faucet

The testnets come with a genesis account (`root`) that can be used to transfer funds to a new account. Run these commands locally to restore the keys on your own machine.

Note: Access to the mnemonic means access to all funds in the account. Don't share or use this mnemonic for non-testing purposes.

```
$ regcli keys add root --recover

# Use the following mnemonic for recovery:
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple

$ regcli tx send --amount 1000wire --to cosmos1lpzffjhasv5qhn7rn6lks9u4dvpzpuj922tdmy --from root --chain-id=wireline --node tcp://registry-testnet.wireline.ninja:26657

# Replace cosmos1lpzffjhasv5qhn7rn6lks9u4dvpzpuj922tdmy with the address you want to transfer funds to.
```

## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
