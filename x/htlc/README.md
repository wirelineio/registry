# HTLC module

## Redeem Flow

Alice creates a HTLC that will pay 25wire to Bob if he can produce the preimage of the SHA256 hash `6815f3c300383519de8e437497e2c3e97852fe8d717a5419d5aafb00cb43c494` within 20 blocks.

```
$ regcli tx htlc add 25wire 6815f3c300383519de8e437497e2c3e97852fe8d717a5419d5aafb00cb43c494 20 $(regcli keys show bob --address) --from=$(regcli keys show alice --address) --chain-id=wireline
```

Check balances. Note that Alice's balance has decreased by 25wire, now locked in the HTLC.

```
$ regcli query account $(regcli keys show alice --address) --indent --chain-id=wireline
$ regcli query account $(regcli keys show bob --address) --indent --chain-id=wireline
```

Bob redeems the HTLC by presenting the preimage `mango` within 20 blocks.

```
$ regcli tx htlc redeem mango --from=$(regcli keys show bob --address) --chain-id=wireline
```

Check balances. Note that Bob's balance has increased by 25wire.

```
$ regcli query account $(regcli keys show alice --address) --indent --chain-id=wireline
$ regcli query account $(regcli keys show bob --address) --indent --chain-id=wireline
```


## Timeout Flow

Alice creates a HTLC that will pay 25wire to Bob if he can produce the preimage of the SHA256 hash `CB4595E84361629F0BB6A3EAA79220E553A8360191433538FF0E0C41F44E30DB ` within 20 blocks.

```
$ regcli tx htlc add 25wire CB4595E84361629F0BB6A3EAA79220E553A8360191433538FF0E0C41F44E30DB 20 $(regcli keys show bob --address) --from=$(regcli keys show alice --address) --chain-id=wireline
```

Alice claims a timeout after waiting for 20 blocks (also try running this before the timeout and see what happens).

```
$ regcli tx htlc fail CB4595E84361629F0BB6A3EAA79220E553A8360191433538FF0E0C41F44E30DB --from=$(regcli keys show alice --address) --chain-id=wireline
```

Check balances. Note that Alice has been refunded 25wire.

```
$ regcli query account $(regcli keys show alice --address) --indent --chain-id=wireline
$ regcli query account $(regcli keys show bob --address) --indent --chain-id=wireline
```