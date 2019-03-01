//
// Copyright 2019 Wireline, Inc.
//

package registry

import (
	"encoding/json"
)

// WirelineChainID is the Cosmos SDK chain ID.
const WirelineChainID = "wireline"

// ID for resources.
type ID string

// Owner represents a resource owner.
type Owner struct {
	// If ID is populated, that will be used (ID of Owner resource record). Else, Address will be used.
	// One of the two MUST be populated.
	ID      ID     `json:"id"`
	Address string `json:"address"`
}

// Resource represents a registry record that can be serialized from/to YAML.
type Resource struct {
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

// PayloadObj represents a signed resource payload.
type PayloadObj struct {
	Resource   ResourceObj `json:"resource"`
	Signatures []Signature `json:"signatures"`
}

// ResourceObj represents a registry record.
type ResourceObj struct {
	ID               ID     `json:"id"`
	Type             string `json:"type"`
	Owner            Owner  `json:"owner"`
	SystemAttributes []byte `json:"systemAttributes"`
	Attributes       []byte `json:"attributes"`
	Links            []byte `json:"links"`
}

// Payload represents a signed resource payload that can be serialized from/to YAML.
type Payload struct {
	Resource   Resource    `json:"resource"`
	Signatures []Signature `json:"signatures"`
}

// ResourceToResourceObj convers Resource to ResourceObj.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func ResourceToResourceObj(resource Resource) ResourceObj {
	var resourceObj ResourceObj

	resourceObj.ID = resource.ID
	resourceObj.Type = resource.Type
	resourceObj.Owner = resource.Owner
	resourceObj.SystemAttributes = MarshalToJSONBytes(resource.SystemAttributes)
	resourceObj.Attributes = MarshalToJSONBytes(resource.Attributes)
	resourceObj.Links = MarshalToJSONBytes(resource.Links)

	return resourceObj
}

// PayloadToPayloadObj converts Payload to PayloadObj object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func PayloadToPayloadObj(payload Payload) PayloadObj {
	var payloadObj PayloadObj

	payloadObj.Resource = ResourceToResourceObj(payload.Resource)
	payloadObj.Signatures = payload.Signatures

	return payloadObj
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

// ResourceObjToResource convers ResourceObj to Resource.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func ResourceObjToResource(resourceObj ResourceObj) Resource {
	var resource Resource

	resource.ID = resourceObj.ID
	resource.Type = resourceObj.Type
	resource.Owner = resourceObj.Owner
	resource.SystemAttributes = UnMarshalJSONBytes(resourceObj.SystemAttributes)
	resource.Attributes = UnMarshalJSONBytes(resourceObj.Attributes)
	resource.Links = UnMarshalJSONBytes(resourceObj.Links)

	return resource
}

// PayloadObjToPayload converts Payload to PayloadObj object.
// Why? Because go-amino can't handle maps: https://github.com/tendermint/go-amino/issues/4.
func PayloadObjToPayload(payloadObj PayloadObj) Payload {
	var payload Payload

	payload.Resource = ResourceObjToResource(payloadObj.Resource)
	payload.Signatures = payloadObj.Signatures

	return payload
}
