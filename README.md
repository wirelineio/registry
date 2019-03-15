# Registry

Registry is a custom blockchain built using the Cosmos SDK. It currently houses SDK modules and CLI toolchain for Registry and Payments.

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

Setup accounts/keys for Alice, Bob & Charlie. Enter a passphrase for the key when prompted.

```
$ regcli keys add alice

$ regcli keys add bob

$ regcli keys add charlie
```

Add initial funds to the accounts.

```
$ registryd add-genesis-account $(regcli keys show alice --address) 1000000wire
$ registryd add-genesis-account $(regcli keys show bob --address) 1000000wire

```

Start the blockchain. You should see blocks being created every few seconds.

```
$ registryd start --gql-server --gql-playground
```

Run the following commands in another terminal.

## Registry module

Get Bob's address and public key.

```
$ regcli query registry key bob
Address   : 002aee66c9908426658a39d7e95a48646d172d0f
PubKey    : 61rphyED+i6I7SuuyeuX9Zgsww9WnXi3BOpxhyEWpnI4kZEfNGY=
```

Sign the resource with Bob's credentials.

```
$ regcli tx registry set service1.yml --from bob --sign-only
Password to sign with 'bob':
Address   : 002aee66c9908426658a39d7e95a48646d172d0f
PubKey    : 61rphyED+i6I7SuuyeuX9Zgsww9WnXi3BOpxhyEWpnI4kZEfNGY=
Signature : iYlLCgiqNL1vsm+3u7alGFNzZJD+u/vlM/YwdJfYAfZAwtChAOUQK3pWlIBIDsmqqwuqV5tK5pDrDcA5zT0swQ==
```

Update the resource payload (e.g. service1.yml) with Bob's address, public key and signature.

* Set `resource/owner/address` to Bob's address.
* Set signature `pubKey` and `sig` using output from the previous command.

```yaml
# service1.yml
resource:
  id: 05013527-30ef-4aee-85d5-a71e1722f255
  type: Service
  owner:
    address: 002aee66c9908426658a39d7e95a48646d172d0f
  systemAttributes:
    uri: https://api.example.org/service
  attributes:
    label: Weather
  links:

signatures:
  -
    pubKey: 61rphyED+i6I7SuuyeuX9Zgsww9WnXi3BOpxhyEWpnI4kZEfNGY=
    sig: iYlLCgiqNL1vsm+3u7alGFNzZJD+u/vlM/YwdJfYAfZAwtChAOUQK3pWlIBIDsmqqwuqV5tK5pDrDcA5zT0swQ==
```

Save resource record (will fail as we're not providing fees).

```
$ regcli tx registry set service1.yml --from alice
```

Save resource record, with fees.

```
$ regcli tx registry set service1.yml --from alice --fee 201wire
```

Verify that the fees have been deducted from Alice's account.

```
$ regcli query account $(regcli keys show alice --address) --indent --chain-id=wireline
```

Get resource record by ID.

```
$ regcli query registry get 05013527-30ef-4aee-85d5-a71e1722f255
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
$ regcli query registry graph f9557e0b-fde4-48ce-923f-7288268473c1 | dot -Tpng  > test.png && eog test.png
```

Delete resource record.

```
$ regcli tx registry delete service1.yml --from alice --fee 201wire
```

Clear all resource records (Warning: This bypasses all access checks and is for local testing purposes only).

```
$ regcli tx registry clear --from alice --fee 201wire
```

### GQL Server API

The GQL server is controlled using the following `registryd` flags:

* `--gql-server` - Enable GQL server.
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:8080/).
* `--gql-port` - Port to run the GQL server on (default 8080).

See `registryd/x/registry/gql/schema.graphql` for the GQL schema.



## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
