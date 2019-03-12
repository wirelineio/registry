// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package gql

type Link struct {
	ID         string  `json:"id"`
	Attributes *string `json:"attributes"`
}

type Owner struct {
	ID      *string `json:"id"`
	Address *string `json:"address"`
}

type Resource struct {
	ID               string  `json:"id"`
	Type             string  `json:"type"`
	Owner            Owner   `json:"owner"`
	SystemAttributes *string `json:"systemAttributes"`
	Attributes       *string `json:"attributes"`
	Links            []Link  `json:"links"`
}
