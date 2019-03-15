# UTXO module

Birth UTXO from account funds.

```
regcli tx utxo birth 100wire --from alice --chain-id=wireline
```

Pay to an address from the UTXO (Account Output).

```
# Note the outpoint hash.
regcli query utxo ls-account-outputs --chain-id=wireline

# Sign but don't broadcast.
regcli tx utxo pay --from alice --chain-id=wireline $(regcli keys show alice --address) $(regcli keys show bob --address) 40 60 A447DEF319B76E111FA557BD6777B86486DAF180973A7842FB25D98A1891AC24 x CAFE --sign-only

# Broadcast with the signature from the previous step.
regcli tx utxo pay --from alice --chain-id=wireline $(regcli keys show alice --address) $(regcli keys show bob --address) 40 60 A447DEF319B76E111FA557BD6777B86486DAF180973A7842FB25D98A1891AC24 x 8AA8B6335F9FB51BA798EBB44B0B4CF25CA279A77C79D306AD4FC34B73406356024B8390F6C0EDD16ED62BB6336526FFAC088D003FE9B9C4A64935FB4B3FBAC8
```

Pay to an address from the newly generated UTXOs.

```
# Sign offline.
regcli tx utxo pay --from bob --chain-id=wireline $(regcli keys show bob --address) $(regcli keys show alice --address) 10 30 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 0 CAFE --sign-only

# Broadcast with signature from previous step.
regcli tx utxo pay --from bob --chain-id=wireline $(regcli keys show bob --address) $(regcli keys show alice --address) 10 30 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 0 5A8350C473BBC066A9B87F715408BC57EFB8868B816FBF43E262DCC1EB6996A533E5E270C3D56CF507AF65A2864DCDB74275BFA42FF9DF1E62F0CC529C026644

# Sign only.
regcli tx utxo pay --from alice --chain-id=wireline $(regcli keys show alice --address) $(regcli keys show bob --address) 40 20 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 1 CAFE --sign-only

# Broadcast with signature from previous step.
# Note: Anyone can broadcast the Cosmos SDK transaction. It's the witness/signature in the payload Tx that counts.
regcli tx utxo pay --from bob --chain-id=wireline $(regcli keys show alice --address) $(regcli keys show bob --address) 40 20 FD1B20785812C1D1D8B776DB424ED79C8CD5A68AC94ED0C26AA8E191A78522CD 1 B3EDC33220D37BC477758068471E1E23F5417E473F6C3194C6EE2EAF9A5FAA7E59B0BFE2DD1A5784F866715453DC2CE40C30F065E9FE6BC18A5B997E998FE8A4

```

List UTXO/Account Outputs.

```
regcli query utxo ls-account-outputs --chain-id=wireline
regcli query utxo ls --chain-id=wireline
```

List transactions.

```
regcli query utxo ls-tx --chain-id=wireline
regcli query utxo get-tx --chain-id=wireline 40154BD90424E96506B172489DC569E2657ECD7A5F145E6D8A1BE5B1F062C2FE
```

View wallet balance and UTXOs.

```
regcli query utxo balance --chain-id=wireline $(regcli keys show alice --address)
regcli query utxo balance --chain-id=wireline $(regcli keys show bob --address)
```

Generate transaction graph.

```
regcli query utxo graph --chain-id=wireline | dot -Tpng  > test.png && eog test.png
```