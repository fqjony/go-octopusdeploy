package octopusdeploy

import (
	"github.com/go-playground/validator/v10"
	uuid "github.com/google/uuid"
)

type RootResource struct {
	Application          string     `json:"Application" validate:"required"`
	Version              string     `json:"Version" validate:"required"`
	APIVersion           string     `json:"ApiVersion" validate:"required"`
	InstallationID       *uuid.UUID `json:"InstallationId" validate:"required"`
	IsEarlyAccessProgram bool       `json:"IsEarlyAccessProgram"`
	HasLongTermSupport   bool       `json:"HasLongTermSupport"`

	resource
}

func NewRootResource() *RootResource {
	return &RootResource{
		resource: *newResource(),
	}
}

// Validate checks the state of the root resource and returns an error if
// invalid.
func (r *RootResource) Validate() error {
	return validator.New().Struct(r)
}
