package httputil

import (
	"net/http"

	"github.com/go-playground/validator/v10"
)

type (
	// ValidatorProvider is an interface for validating input.
	ValidatorProvider interface {
		// Validate validates the struct returning any error encountered.
		Validate(input any) error
		// BindJSONAndValidate binds the JSON input and validates the struct returning any error encountered.
		BindJSONAndValidate(w http.ResponseWriter, req *http.Request, input any) error
	}

	DefaultValidator struct {
		p *validator.Validate
	}
)

// NewValidator creates a new validator provider.
// The default provider is the go-playground/validator.
func NewValidator(provider string) ValidatorProvider {
	return &DefaultValidator{
		p: validator.New(),
	}
}

func (gp *DefaultValidator) Validate(input any) error {
	return gp.p.Struct(input)
}

func (gp *DefaultValidator) ValidateInput(input any, tag string) error {
	return gp.p.Var(input, tag)
}

func (gp *DefaultValidator) BindJSONAndValidate(
	w http.ResponseWriter,
	req *http.Request,
	target any,
) error {
	if err := BindJSON(w, req, target); err != nil {
		return err
	}

	return gp.p.Struct(target)
}
