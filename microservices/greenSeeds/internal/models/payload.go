package models

// For validating request data
type Payload struct {
	FullName string `validate:"required"`
}
