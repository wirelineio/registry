# Record Registration

Note: This README assumes you already have the Registry up and running. See the README in the repo root for install/setup instructions.

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

## Clear Remote Registry

To clear a remote registry, you need to know:

* The RPC endpoint of the remote registry (e.g. see https://github.com/wirelineio/registry#testnets).
* The mnemonic for an account that has funds on the registry.

The following example will work for https://registry-testnet.dev.wireline.ninja.

Create an account on your machine, using the mnemonic for the remote `root` account.

```
$ regcli keys add root-testnet-dev --recover
# Enter a passphrase for the new account, repeat it when prompted.
# Use the following mnemonic for recovery:
# salad portion potato insect unknown exile lion soft layer evolve flavor hollow emerge celery ankle sponsor easy effort flush furnace life maximum rotate apple
```

Clear the remote registry using the following command:

```
$ regcli tx registry clear --from root-testnet-dev --node tcp://registry-testnet.dev.wireline.ninja:26657
# Enter passphrase when prompted.
```

Note: You might see an error on the first couple of runs: 'ERROR: broadcast_tx_commit: Post http://registry-testnet.dev.wireline.ninja:26657: EOF'. It's not clear why this happens, but ignore the error and just try again.

Use the GQL playground (https://registry-testnet.dev.wireline.ninja) to confirm that all records are gone.
