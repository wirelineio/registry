# Environment Setup

Installing lnd and btcd.

```bash
go get -d github.com/lightningnetwork/lnd
cd $GOPATH/src/github.com/lightningnetwork/lnd
make && make install
make btcd
```

Installing btcctl.

```bash
cd $GOPATH/src/github.com/btcsuite/btcd
GO111MODULE=on go install -v . ./cmd/...
```

Generate certificates for btcd (if required).

```bash
cd $GOPATH/src/github.com/wirelineio/wirechain/env/local/btcd
gencerts --host="*" --directory=./rpc --force
```
