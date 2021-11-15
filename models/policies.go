package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// A list of errors returned from package
var (
	ErrorReadingBody = errors.New("failed to read message body")
	ErrorParsingBody = errors.New("failed to parse json body")
)

type Condition struct {
	Attributes []string `bson:"attributes"          json:"attributes"`
	Operator   string   `bson:"operator"          json:"operator"`
	Values     []string `bson:"Values"          json:"values"`
}

type Policy struct {
	ID         string      `bson:"_id"          json:"id,omitempty"`
	Entities   []string    `bson:"entities"   json:"entities"`
	Role       string      `bson:"role"      json:"role"`
	Conditions []Condition `bson:"conditions" json:"conditions,omitempty"`
}

type NewPolicy struct {
	Entities   []string    `json:"entities"`
	Role       string      `json:"role"`
	Conditions []Condition `json:"conditions,omitempty"`
}

// policies permissions
const (
	PoliciesRead string = "policies:read"
	PoliciesCreate      = "policies:create"
	PoliciesUpdate      = "policies:update"
	PoliciesDelete      = "policies:delete"
)


// ValidateNewPolicy checks that all the mandatory fields are non-empty
func (policy *NewPolicy) ValidateNewPolicy() error {

	var invalidFields []string

	if len(policy.Entities) == 0 {
		invalidFields = append(invalidFields, "entities")
	}

	if len(policy.Role) == 0 {
		invalidFields = append(invalidFields, "role")
	}

	if invalidFields != nil {
		return fmt.Errorf("missing mandatory fields: %v", strings.Join(invalidFields, ", "))
	}

	return nil
}

//CreateNewPolicy manages the creation of a filter from reader
func CreateNewPolicy(reader io.Reader) (*NewPolicy, error) {

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, ErrorReadingBody
	}

	var policy NewPolicy
	err = json.Unmarshal(bytes, &policy)
	if err != nil {
		return nil, ErrorParsingBody
	}

	return &policy, nil
}
