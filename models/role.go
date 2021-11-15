package models

//Roles represents an array of the role model
type Roles struct {
	Count      int    `json:"count"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
	Items      []Role `json:"items"`
	TotalCount int    `json:"total_count"`
}

//Role represents the structure for a role
type Role struct {
	ID          string   `bson:"_id" json:"id"`
	Name        string   `bson:"name" json:"name"`
	Permissions []string `bson:"permissions" json:"permissions"`
}

// roles permissions
const (
	RolesRead string = "roles:read"
	RolesCreate      = "roles:create"
	RolesUpdate      = "roles:update"
	RolesDelete      = "roles:delete"
)
