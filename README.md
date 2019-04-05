# Registry

Registry is a custom blockchain built using the Cosmos SDK. It currently houses SDK modules and CLI toolchain for Registry and Payments.

## Setup & install

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

## Initialize & start blockchain

Delete existing blockchain and config.

```
$ rm -rf ~/.registryd/ ~/.regcli
```

Initialize the chain.

```
$ registryd init --chain-id wireline
```

Setup accounts/keys for Alice & Bob. Enter a passphrase for the key when prompted.

```
$ regcli keys add alice

$ regcli keys add bob
```

Add initial funds to the accounts.

```
$ registryd add-genesis-account $(regcli keys show alice --address) 1000000wire

```

Start the blockchain. You should see blocks being created every few seconds.

```
$ registryd start --gql-server --gql-playground
```

Run the following commands in another terminal.

## Register module

Get Bob's address and public key.

```
$ regcli query registry key bob
Address   : 002aee66c9908426658a39d7e95a48646d172d0f
PubKey    : 61rphyED+i6I7SuuyeuX9Zgsww9WnXi3BOpxhyEWpnI4kZEfNGY=
```

Sign the resource with Bob's credentials.

TODO: service1.yml is not defined.

```
$ regcli tx registry set service1.yml --from bob --sign-only
Password to sign with 'bob':
Address   : 002aee66c9908426658a39d7e95a48646d172d0f
PubKey    : 61rphyED+i6I7SuuyeuX9Zgsww9WnXi3BOpxhyEWpnI4kZEfNGY=
Signature : iYlLCgiqNL1vsm+3u7alGFNzZJD+u/vlM/YwdJfYAfZAwtChAOUQK3pWlIBIDsmqqwuqV5tK5pDrDcA5zT0swQ==
```

Update the resource payload (e.g. service1.yml) with Bob's address, public key and signature.

* Set `resource/owner` to Bob's address.
* Set signature `pubKey` and `sig` using output from the previous command.

```yaml
# service1.yml
record:
  id: wrn:record:05013527-30ef-4aee-85d5-a71e1722f255
  type: wrn:registry-type:service
  owner: 02e840ed2d4c3e0b4e068f0d4be811b095ec78d5
  attributes:
    label: Weather

signatures:
  -
    pubKey: 61rphyEDI/Iy96OBr9fn11ADRfDPUgAiEW5MdETVuK9PohsxWMU=
    sig: r3J9Hi+1nyO86Gbdo0jRuxzU1zHRzEvtK3EqH2x9owQ9NNvzQp7BeBLyInASgwEDHu4Iec21fzRR8klHbDN5Sw==
```

Publish resource record.

```
$ regcli tx registry set service1.yml --from alice
```

Get resource record by ID.

```
$ regcli query registry get wrn:record:05013527-30ef-4aee-85d5-a71e1722f255
```

List resource records.

```
$ regcli query registry list
```

Generate resource graph.

```
$ regcli query registry graph | dot -Tpng  > test.png && eog test.png
```

Generate graph, starting from a particular resource.

```
$ regcli query registry graph wrn:record:05013527-30ef-4aee-85d5-a71e1722f255 | dot -Tpng  > test.png && eog test.png
```

Delete resource record.

```
$ regcli tx registry delete service1.yml --from alice
```

Clear all resource records (Warning: This bypasses all access checks and is for local testing purposes only).

```
$ regcli tx registry clear --from alice
```

### GQL Server API

The GQL server is controlled using the following `registryd` flags:

* `--gql-server` - Enable GQL server.
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:8080/).
* `--gql-port` - Port to run the GQL server on (default 8080).

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

The testnets come with a seed account (`alice`) that can be used to transfer funds to a new account. Run these commands locally to restore Alice's keys on your own machine to transfer funds.

```
$ regcli keys add alice --recover

# Use the following mnemonic for recovery:
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple

$ regcli tx send --amount 1000wire --to cosmos1lpzffjhasv5qhn7rn6lks9u4dvpzpuj922tdmy --from alice --chain-id=wireline --node tcp://registry-testnet.wireline.ninja:26657

# Replace cosmos1lpzffjhasv5qhn7rn6lks9u4dvpzpuj922tdmy with the address you want to transfer funds to.
```

## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
