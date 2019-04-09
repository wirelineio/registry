// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

type Account struct {
	Address  string  `json:"address"`
	PubKey   *string `json:"pubKey"`
	Number   BigUInt `json:"number"`
	Sequence BigUInt `json:"sequence"`
	Balance  []Coin  `json:"balance"`
}

type Bot struct {
	Record    *Record `json:"record"`
	Name      string  `json:"name"`
	AccessKey *string `json:"accessKey"`
}

type Coin struct {
	Type   string  `json:"type"`
	Amount BigUInt `json:"amount"`
}

type KeyValue struct {
	Key   string `json:"key"`
	Value Value  `json:"value"`
}

type KeyValueInput struct {
	Key   string     `json:"key"`
	Value ValueInput `json:"value"`
}

type Record struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Owner      string      `json:"owner"`
	Attributes []*KeyValue `json:"attributes"`
}

type Status struct {
	Version string `json:"version"`
}

type Value struct {
	Null    *bool    `json:"null"`
	Int     *int     `json:"int"`
	Float   *float64 `json:"float"`
	String  *string  `json:"string"`
	Boolean *bool    `json:"boolean"`
	Values  []*Value `json:"values"`
}

type ValueInput struct {
	Null    *bool         `json:"null"`
	Int     *int          `json:"int"`
	Float   *float64      `json:"float"`
	String  *string       `json:"string"`
	Boolean *bool         `json:"boolean"`
	Values  []*ValueInput `json:"values"`
}
