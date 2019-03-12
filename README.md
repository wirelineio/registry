# Wirechain

Wirechain is a custom blockchain built using the Cosmos SDK. It currently houses SDK modules and CLI toolchain for Registry and Payments.

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

Setup accounts/keys for Alice, Bob & Charlie. Enter a passphrase for the key when prompted.

```
$ wirecli keys add alice

$ wirecli keys add bob

$ wirecli keys add charlie
```

Add initial funds to the accounts.

```
$ wirechaind add-genesis-account $(wirecli keys show alice --address) 1000000wire
$ wirechaind add-genesis-account $(wirecli keys show bob --address) 1000000wire

```

Start the blockchain. You should see blocks being created every few seconds.

```
$ wirechaind start --gql-server --gql-playground
```

Run the following commands in another terminal.

## Registry module

Get Bob's address and public key.

```
$ wirecli query registry key bob
Address   : 002aee66c9908426658a39d7e95a48646d172d0f
PubKey    : 61rphyED+i6I7SuuyeuX9Zgsww9WnXi3BOpxhyEWpnI4kZEfNGY=
```

Sign the resource with Bob's credentials.

```
$ wirecli tx registry set service1.yml --from bob --sign-only
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
$ wirecli tx registry set service1.yml --from alice
```

Save resource record, with fees.

```
$ wirecli tx registry set service1.yml --from alice --fee 201wire
```

Verify that the fees have been deducted from Alice's account.

```
$ wirecli query account $(wirecli keys show alice --address) --indent --chain-id=wireline
```

Get resource record by ID.

```
$ wirecli query registry get 05013527-30ef-4aee-85d5-a71e1722f255
```

List resource records.

```
$ wirecli query registry list
```

Generate resource graph.

```
$ wirecli query registry graph | dot -Tpng  > test.png && eog test.png
```

Generate graph, starting from a particular resource.

```
$ wirecli query registry graph f9557e0b-fde4-48ce-923f-7288268473c1 | dot -Tpng  > test.png && eog test.png
```

Delete resource record.

```
$ wirecli tx registry delete service1.yml --from alice --fee 201wire
```

Clear all resource records (Warning: This bypasses all access checks and is for local testing purposes only).

```
$ wirecli tx registry clear --from alice --fee 201wire
```

### GQL Server API

The GQL server is controlled using the following `wirechaind` flags:

* `--gql-server` - Enable GQL server.
* `--gql-playground` - Enable GQL playground app (Available at http://localhost:8080/).
* `--gql-port` - Port to run the GQL server on (default 8080).

See `wirechaind/x/registry/gql/schema.graphql` for the GQL schema.

## UTXO module

Birth UTXO from account funds.

```
wirecli tx utxo birth 100wire --from alice --chain-id=wireline
```

Pay to an address from the UTXO (Account Output).

```
# Note the outpoint hash.
wirecli query utxo ls-account-outputs --chain-id=wireline

# Sign but don't broadcast.
wirecli tx utxo pay --from alice --chain-id=wireline $(wirecli keys show alice --address) $(wirecli keys show bob --address) 40 60 A447DEF319B76E111FA557BD6777B86486DAF180973A7842FB25D98A1891AC24 x CAFE --sign-only

# Broadcast with the signature from the previous step.
wirecli tx utxo pay --from alice --chain-id=wireline $(wirecli keys show alice --address) $(wirecli keys show bob --address) 40 60 A447DEF319B76E111FA557BD6777B86486DAF180973A7842FB25D98A1891AC24 x 8AA8B6335F9FB51BA798EBB44B0B4CF25CA279A77C79D306AD4FC34B73406356024B8390F6C0EDD16ED62BB6336526FFAC088D003FE9B9C4A64935FB4B3FBAC8
```

Pay to an address from the newly generated UTXOs.

```
# Sign offline.
wirecli tx utxo pay --from bob --chain-id=wireline $(wirecli keys show bob --address) $(wirecli keys show alice --address) 10 30 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 0 CAFE --sign-only

# Broadcast with signature from previous step.
wirecli tx utxo pay --from bob --chain-id=wireline $(wirecli keys show bob --address) $(wirecli keys show alice --address) 10 30 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 0 5A8350C473BBC066A9B87F715408BC57EFB8868B816FBF43E262DCC1EB6996A533E5E270C3D56CF507AF65A2864DCDB74275BFA42FF9DF1E62F0CC529C026644

# Sign only.
wirecli tx utxo pay --from alice --chain-id=wireline $(wirecli keys show alice --address) $(wirecli keys show bob --address) 40 20 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 1 CAFE --sign-only

# Broadcast with signature from previous step.
# Note: Anyone can broadcast the Cosmos SDK transaction. It's the witness/signature in the payload Tx that counts.
wirecli tx utxo pay --from bob --chain-id=wireline $(wirecli keys show alice --address) $(wirecli keys show bob --address) 40 20 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 1 B3EDC33220D37BC477758068471E1E23F5417E473F6C3194C6EE2EAF9A5FAA7E59B0BFE2DD1A5784F866715453DC2CE40C30F065E9FE6BC18A5B997E998FE8A4

```

List UTXO/Account Outputs.

```
wirecli query utxo ls-account-outputs --chain-id=wireline
wirecli query utxo ls --chain-id=wireline
```

List transactions.

```
wirecli query utxo ls-tx --chain-id=wireline
wirecli query utxo get-tx --chain-id=wireline 40154BD90424E96506B172489DC569E2657ECD7A5F145E6D8A1BE5B1F062C2FE
```

View wallet balance and UTXOs.

```
wirecli query utxo balance --chain-id=wireline $(wirecli keys show alice --address)
wirecli query utxo balance --chain-id=wireline $(wirecli keys show bob --address)
```

Generate transaction graph.

```
wirecli query utxo graph --chain-id=wireline | dot -Tpng  > test.png && eog test.png
```

## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
