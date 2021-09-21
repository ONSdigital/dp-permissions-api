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
	Attributes []string
	Operator   string
	Values     []string
}

type Policy struct {
	ID         string      `bson:"-"          json:"id,omitempty"`
	Entities   []string    `bson:"entities"   json:"entities"`
	Roles      []string    `bson:"roles"      json:"roles"`
	Conditions []Condition `bson:"conditions" json:"conditions,omitempty"`
}

type NewPolicy struct {
	Entities   []string    `json:"entities"`
	Roles      []string    `json:"roles"`
	Conditions []Condition `json:"conditions,omitempty"`
}

// ValidateNewPolicy checks that all the mandatory fields are non-empty
func (policy *NewPolicy) ValidateNewPolicy() error {

	var invalidFields []string

	if len(policy.Entities) == 0 {
		invalidFields = append(invalidFields, "entities")
	}

	if len(policy.Roles) == 0 {
		invalidFields = append(invalidFields, "roles")
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
