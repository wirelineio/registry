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
