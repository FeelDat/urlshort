// Package models provides data structures and types that are used across the application.
package models

// CtxKey represents a type for context keys.
//
// This type is used to ensure that the application's context values are stored and accessed
// in a type-safe manner, thereby avoiding potential collisions with values from other packages.
//
// Usage:
//
//	cntx := context.WithValue(r.Context(), models.CtxKey("userID"), userID)
type CtxKey string
