package models

// EntityIDToPolicies maps an entity ID to a slice of policies.
type EntityIDToPolicies map[string][]*BundlePolicy

// Bundle is the optimised lookup table for permissions.
type Bundle map[string]EntityIDToPolicies

// BundlePolicy represents a policy tailored for the permissions bundle.
// The permissions bundle json does not include the entities and role fields.
type BundlePolicy struct {
	ID        string    `bson:"_id"          json:"id,omitempty"`
	Entities  []string  `bson:"entities"   json:"-"`
	Role      string    `bson:"role"      json:"-"`
	Condition Condition `bson:"condition" json:"condition,omitempty"`
}
