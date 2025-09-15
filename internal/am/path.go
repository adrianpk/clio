package am

import (
	"fmt"

	"github.com/google/uuid"
)

// ListPath returns the path for listing resources
func ListPath(i Identifiable) string {
	return fmt.Sprintf("list-%s", Plural(i.Type()))
}

// NewPath returns the path for creating a new resource
func NewPath(i Identifiable) string {
	return fmt.Sprintf("new-%s", i.Type())
}

// CreatePath returns the path for creating a resource
func CreatePath(i Identifiable) string {
	return fmt.Sprintf("create-%s", i.Type())
}

// ShowPath returns the path for showing a resource
func ShowPath(i Identifiable, id uuid.UUID) string {
	return fmt.Sprintf("%s/%s", i.Type(), id)
}

// EditPath returns the path for editing a resource
func EditPath(i Identifiable, id uuid.UUID) string {
	return fmt.Sprintf("edit-%s?id=%s", i.Type(), id)
}

// UpdatePath returns the path for updating a resource
func UpdatePath(i Identifiable) string {
	return fmt.Sprintf("update-%s", i.Type())
}

// DeletePath returns the path for deleting a resource
func DeletePath(i Identifiable, id uuid.UUID) string {
	return fmt.Sprintf("delete-%s/%s", i.Type(), id)
}

// ListRelatedPath returns the path for listing related resources
func ListRelatedPath(i Identifiable, j Identifiable, id uuid.UUID) string {
	return fmt.Sprintf("list-%s-%s?id=%s", i.Type(), j.Type(), id)
}

// AddRelatedPath returns the path for adding a related resource
func AddRelatedPath(i Identifiable, j Identifiable) string {
	return fmt.Sprintf("add-%s-to-%s", j.Type(), i.Type())
}

// RemoveRelatedPath returns the path for removing a related resource
func RemoveRelatedPath(i Identifiable, j Identifiable) string {
	return fmt.Sprintf("remove-%s-from-%s", j.Type(), i.Type())
}
