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

$ wirecli keys add charlie
```

Add initial funds to the accounts.

```
$ wirechaind add-genesis-account $(wirecli keys show alice --address) 1000wire
$ wirechaind add-genesis-account $(wirecli keys show bob --address) 1000wire

```

Start the blockchain. You should see blocks being created every few seconds.


```
$ wirechaind start
```

Run the following commands in another terminal.

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

## Multisig module

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

## References

* https://golang.org/doc/install
* https://github.com/cosmos/cosmos-sdk
* https://cosmos.network/docs/tutorial/
* https://github.com/cosmos/sdk-application-tutorial
