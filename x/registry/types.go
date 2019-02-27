//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ID for resources.
type ID string

// Owner represents a resource owner.
type Owner struct {
	// If ID is populated, that will be used (ID of Owner resource record). Else, Address will be used.
	// One of the two MUST be populated.
	ID      ID             `json:"id"`
	Address sdk.AccAddress `json:"address"`
}

// ResourceYaml represents a registry record that can be serialized from/to YAML.
type ResourceYaml struct {
	ID               ID                     `json:"id"`
	Type             string                 `json:"type"`
	Owner            Owner                  `json:"owner"`
	SystemAttributes map[string]interface{} `json:"systemAttributes"`
	Attributes       map[string]interface{} `json:"attributes"`
	Links            map[string]interface{} `json:"links"`
}

// Signature represents a resource signature.
type Signature struct {
	PubKey    string `json:"pubKey"`
	Signature string `json:"sig"`
}

// Payload represents a signed resource payload.
type Payload struct {
	Resource   Resource    `json:"resource"`
	Signatures []Signature `json:"signatures"`
}

// Resource represents a registry record.
type Resource struct {
	ID               ID     `json:"id"`
	Type             string `json:"type"`
	Owner            Owner  `json:"owner"`
	SystemAttributes []byte `json:"systemAttributes"`
	Attributes       []byte `json:"attributes"`
	Links            []byte `json:"links"`
}

// PayloadYaml represents a signed resource payload that can be serialized from/to YAML.
type PayloadYaml struct {
	Resource   ResourceYaml `json:"resource"`
	Signatures []Signature  `json:"signatures"`
}

// ResourceYamlToResource convers ResourceYaml to Resource.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func ResourceYamlToResource(resourceYaml ResourceYaml) Resource {
	var resource Resource

	resource.ID = resourceYaml.ID
	resource.Type = resourceYaml.Type
	resource.Owner = resourceYaml.Owner
	resource.SystemAttributes = MarshalToJSONBytes(resourceYaml.SystemAttributes)
	resource.Attributes = MarshalToJSONBytes(resourceYaml.Attributes)
	resource.Links = MarshalToJSONBytes(resourceYaml.Links)

	return resource
}

// PayloadYamlToPayload converts PayloadYaml to Payload object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func PayloadYamlToPayload(payloadYaml PayloadYaml) Payload {
	var payload Payload

	payload.Resource = ResourceYamlToResource(payloadYaml.Resource)
	payload.Signatures = payloadYaml.Signatures

	return payload
}

// MarshalToJSONBytes converts map[string]interface{} to bytes.
func MarshalToJSONBytes(val map[string]interface{}) (bytes []byte) {
	bytes, err := json.Marshal(val)
	if err != nil {
		panic("Marshal error.")
	}

	return
}

// UnMarshalJSONBytes converts bytes to map[string]interface{}.
func UnMarshalJSONBytes(bytes []byte) map[string]interface{} {
	var val map[string]interface{}
	err := json.Unmarshal(bytes, &val)

	if err != nil {
		panic("Marshal error.")
	}

	return val
}

// ResourceToResourceYaml convers Resource to ResourceYaml.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func ResourceToResourceYaml(resource Resource) ResourceYaml {
	var resourceYaml ResourceYaml

	resourceYaml.ID = resource.ID
	resourceYaml.Type = resource.Type
	resourceYaml.Owner = resource.Owner
	resourceYaml.SystemAttributes = UnMarshalJSONBytes(resource.SystemAttributes)
	resourceYaml.Attributes = UnMarshalJSONBytes(resource.Attributes)
	resourceYaml.Links = UnMarshalJSONBytes(resource.Links)

	return resourceYaml
}

// PayloadToPayloadYaml converts PayloadYaml to Payload object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func PayloadToPayloadYaml(payload Payload) PayloadYaml {
	var payloadYaml PayloadYaml

	payloadYaml.Resource = ResourceToResourceYaml(payload.Resource)
	payloadYaml.Signatures = payload.Signatures

	return payloadYaml
}
