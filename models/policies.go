package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

// policies permissions
const (
	PoliciesRead   string = "policies:read"
	PoliciesCreate        = "policies:create"
	PoliciesUpdate        = "policies:update"
	PoliciesDelete        = "policies:delete"

	OperatorStringEquals Operator = "StringEquals"
	OperatorStartsWith   Operator = "StartsWith"
)

// A list of errors returned from package
var (
	ErrorReadingBody = errors.New("failed to read message body")
	ErrorParsingBody = errors.New("failed to parse json body")
)

type Operator string

//Condition represents the conditions to be applied for a policy
type Condition struct {
	Attribute string   `bson:"attribute" json:"attribute"`
	Operator  Operator `bson:"operator" json:"operator"`
	Values    []string `bson:"Values" json:"values"`
}

//Policy represent a structure for a policy in DB
type Policy struct {
	ID        string    `bson:"_id"          json:"id,omitempty"`
	Entities  []string  `bson:"entities"   json:"entities"`
	Role      string    `bson:"role"      json:"role"`
	Condition Condition `bson:"condition" json:"condition,omitempty"`
}

//UpdateResult represent a result of the upsert policy
type UpdateResult struct {
	ModifiedCount int
	UpsertedCount int
}

//PolicyInfo contains properties required to create or update a policy
type PolicyInfo struct {
	Entities  []string  `json:"entities"`
	Role      string    `json:"role"`
	Condition Condition `json:"condition,omitempty"`
}

func (operator Operator) IsValid() bool {
	operators := map[Operator]struct{}{
		OperatorStringEquals: {},
		OperatorStartsWith:   {},
	}
	_, ok := operators[operator]
	return ok
}

func (operator Operator) String() string {
	return string(operator)
}

// GetPolicy creates a policy object with ID
func (policy *PolicyInfo) GetPolicy(id string) *Policy {
	return &Policy{
		ID:        id,
		Entities:  policy.Entities,
		Role:      policy.Role,
		Condition: policy.Condition,
	}
}

// ValidatePolicy checks that all the mandatory fields are non-empty and non-empty fields contain valid values
func (policy *PolicyInfo) ValidatePolicy() error {

	var missingFields, invalidFields, validationErrors []string

	if len(policy.Entities) == 0 {
		missingFields = append(missingFields, "entities")
	}
	if len(policy.Role) == 0 {
		missingFields = append(missingFields, "role")
	}
	if len(missingFields) > 0 {
		validationErrors = append(validationErrors, fmt.Sprintf("missing mandatory fields: %v", strings.Join(missingFields, ", ")))
	}

	if len(policy.Condition.Operator) > 0 {
		if !policy.Condition.Operator.IsValid() {
			invalidFields = append(invalidFields, "condition operator "+policy.Condition.Operator.String())
		}
	}
	if len(invalidFields) > 0 {
		validationErrors = append(validationErrors, fmt.Sprintf("invalid field values: %v", strings.Join(invalidFields, ", ")))
	}

	if len(validationErrors) > 0 {
		return fmt.Errorf(strings.Join(validationErrors, ". "))
	}
	return nil
}

//CreatePolicy manages the creation of a filter from reader
func CreatePolicy(reader io.Reader) (*PolicyInfo, error) {

	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, ErrorReadingBody
	}

	var policy PolicyInfo
	err = json.Unmarshal(bytes, &policy)
	if err != nil {
		return nil, ErrorParsingBody
	}

	return &policy, nil
}
