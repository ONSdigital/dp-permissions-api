package sdk

// EntityIDToPolicies maps an entity ID to a slice of policies.
type EntityIDToPolicies map[string][]Policy

// Bundle is the optimised lookup table for permissions.
type Bundle map[string]EntityIDToPolicies

// Policy is the policy model as stored in the permissions API.
type Policy struct {
	ID        string    `json:"id"`
	Condition Condition `json:"condition"`
}

// Condition is used within a policy to match additional attributes.
type Condition struct {
	Attribute string   `json:"attribute"`
	Operator  Operator `json:"operator"`
	Values    []string `json:"values"`
}

// Operator is used to define a set of supported Condition operators
type Operator string

// EntityData groups the different entity types into a single parameter
type EntityData struct {
	UserID string
	Groups []string
}

const (
	OperatorStringEquals Operator = "StringEquals"
	OperatorStartsWith   Operator = "StartsWith"
)
