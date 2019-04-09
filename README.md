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

If you just wanted to run the Registry locally, stop now. Read on if you want to explore basic usage of the in-built CLI.


## Record Registration

Setup an account for Alice. Enter a passphrase for the key (e.g. test12345) when prompted. Note the generated mnemonic to restore the private key later.

```
$ regcli keys add alice
```

Get Alice's address and public key.

```
$ regcli query registry key alice
Address   : 002aee66c9908426658a39d7e95a48646d172d0f
PubKey    : 61rphyED+i6I7SuuyeuX9Zgsww9WnXi3BOpxhyEWpnI4kZEfNGY=
```

Create a payload file (e.g. service1.yml) with Alice's address as the `owner`.

```yaml
# service1.yml
record:
  id: wrn:record:05013527-30ef-4aee-85d5-a71e1722f255
  type: wrn:registry-type:service
  owner: 02e840ed2d4c3e0b4e068f0d4be811b095ec78d5
  attributes:
    label: Weather
```

Sign the payload with Alice's private key.

```
$ regcli tx registry set service1.yml --from alice --sign-only
Password to sign with 'alice':
Address   : 02e840ed2d4c3e0b4e068f0d4be811b095ec78d5
PubKey    : 61rphyEDI/Iy96OBr9fn11ADRfDPUgAiEW5MdETVuK9PohsxWMU=
Signature : r3J9Hi+1nyO86Gbdo0jRuxzU1zHRzEvtK3EqH2x9owQ9NNvzQp7BeBLyInASgwEDHu4Iec21fzRR8klHbDN5Sw==
```

Update the resource payload (e.g. service1.yml) with Alice's public key (`pubKey`) and signature (`sig`), using output from the previous command.

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

Publish resource record to the blockchain using the `root` account.

```
$ regcli tx registry set service1.yml --from root
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
$ regcli tx registry delete service1.yml --from root
```

Clear all resource records (Warning: This bypasses all access checks and is for local testing purposes only).

```
$ regcli tx registry clear --from root
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
