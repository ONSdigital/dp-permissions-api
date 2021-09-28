package models

// EntityIDToPolicies maps an entity ID to a slice of policies.
type EntityIDToPolicies map[string][]*Policy

// PermissionToEntityLookup maps a permission ID to the next level in the lookup table - the EntityIDToPolicies map.
type PermissionToEntityLookup map[string]EntityIDToPolicies

// Bundle is the optimised lookup table for permissions.
type Bundle struct {
	PermissionToEntityLookup
}
