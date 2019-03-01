# Multisig module

## Normal Operation

Alice creates a multisig contract and contributes 10wire, requiring Bob to deposit the same amount.

```
$ wirecli tx multisig init test1 10wire 10wire $(wirecli keys show bob --address) --from=$(wirecli keys show alice --address) --chain-id=wireline
```

Anyone can view the state of the contract.

```
$ wirecli query multisig view test1 --chain-id=wireline
```

Bob joins the contract.

```
$ wirecli tx multisig join test1 10wire --from=$(wirecli keys show bob --address) --chain-id=wireline
```

Spending funds from the contract now requires both Alice's and Bob's signatures.

```
$ wirecli tx multisig spend test1 5wire $(wirecli keys show charlie --address) $(wirecli keys show bob --address) --from=$(wirecli keys show alice --address) --chain-id=wireline
Password to sign with 'alice':
ERROR: {"codespace":"sdk","code":4,"message":"wrong number of signers"}
```

To spend the funds, Alice must first generate a transaction and save it to a file.

```
$ wirecli tx multisig spend test1 5wire $(wirecli keys show charlie --address) $(wirecli keys show bob --address) --from=$(wirecli keys show alice --address) --chain-id=wireline --generate-only > tx.json
```

Alice must then sign the transaction, saving the result to another file.

```
$ wirecli tx sign tx.json --chain-id=wireline --name=alice > tx_signed_alice.json
```

Alice sends this file to Bob to get his signature (e.g. over email or chat). Bob inspects the file and verifies the contract ID, amount and target address.

```
{
  "type": "auth/StdTx",
  "value": {
    "msg": [
      {
        "type": "multisig/SpendMultiSig",
        "value": {
          "ID": "test1",
          "Amount": {
            "denom": "wire",
            "amount": "5"
          },
          "ToAddress": "cosmos1nfagqae3hmph7ac88tzc0vwsp9msn4hn57svl2",
          "AliceAddress": "cosmos1gq9vx70vqlcnak37sh4mva5lece576zdrt4nzv",
          "BobAddress": "cosmos1pdqu33vnu8y98q9yqy3x8jlcusxeuezx8rutan"
        }
      }
    ],
    "fee": {
      "amount": [
        {
          "denom": "",
          "amount": "0"
        }
      ],
      "gas": "200000"
    },
    "signatures": [
      {
        "pub_key": {
          "type": "tendermint/PubKeySecp256k1",
          "value": "AmyY3MoPUJ1rpYE+snIvfSGPEoxCmo3XYm2h0L3W1JaF"
        },
        "signature": "eCXBtG/44KipehmzS5Scule29IBT82+aY+0BCZKWm34epVveHZ9HTPzRdjHhzXrzpYy54WarI/+pAHhQuBDgIg==",
        "account_number": "0",
        "sequence": "19"
      }
    ],
    "memo": ""
  }
}
```

Bob signs the transaction, saving the result to another file.

```
$ wirecli tx sign tx_signed_alice.json --chain-id=wireline --name=bob > tx_signed_alice_and_bob.json
```

The transaction can now be broadcast to the blockchain by any party and doesn't need further signatures.

```
$ wirecli tx broadcast tx_signed_alice_and_bob.json --chain-id=wireline
```

Verify that the target account (Charlie) has got the funds.

```
$ wirecli query account $(wirecli keys show charlie --address) --indent --chain-id=wireline
```

Check that the funds have been deducted from the contract.

```
$ wirecli query multisig view test1 --chain-id=wireline
```

## Contract Abort

Alice creates another multisig contract and contributes 10wire, requiring Bob to deposit 100wire.

```
$ wirecli tx multisig init test2 10wire 100wire $(wirecli keys show bob --address) --from=$(wirecli keys show alice --address) --chain-id=wireline
```

Note the contract details, especially the 'State', and Alice's balance.

```
$ wirecli query multisig view test2 --chain-id=wireline
$ wirecli query account $(wirecli keys show alice --address) --indent --chain-id=wireline
```

Before Bob joins, Alice can abort the contract and recover her funds.

```
$ wirecli tx multisig abort test2 --from=$(wirecli keys show alice --address) --chain-id=wireline
```

Confirm that the contract is deleted and Alice has been refunded the funds locked in the contract.

```
$ wirecli query multisig view test2 --chain-id=wireline
$ wirecli query account $(wirecli keys show alice --address) --indent --chain-id=wireline
```
