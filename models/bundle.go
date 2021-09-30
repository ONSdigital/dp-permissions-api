package models

// EntityIDToPolicies maps an entity ID to a slice of policies.
type EntityIDToPolicies map[string][]*Policy

// Bundle is the optimised lookup table for permissions.
type Bundle map[string]EntityIDToPolicies
