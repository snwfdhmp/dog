package models

import "github.com/go-ozzo/ozzo-validation"

// {{.ResourceName}} represents an {{unexported .ResourceName}} record.
type {{.ResourceName}} struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Validate validates the {{.ResourceName}} fields.
func (m {{.ResourceName}}) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.Name, validation.Required, validation.Length(0, 120)),
	)
}
