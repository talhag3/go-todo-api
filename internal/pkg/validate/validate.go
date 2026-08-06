// pkg/validate/validate.go
//
// Wraps go-playground/validator so the handler layer never builds
// error maps by hand. One call: validate.Struct(req). If it fails, you
// get back a ready-to-use *apperror.AppError with a details map keyed
// by JSON field name — the same object response.Error already knows
// how to serialize.
package validate

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/talhag3/todoapp/internal/pkg/apperror"
)

// validator.Validate is safe for concurrent use once built, so we build
// it exactly once at package init time — not per-request. Same idea as
// keeping a single pgxpool.Pool alive instead of opening a connection
// per query.
var v *validator.Validate

func init() {
	v = validator.New(validator.WithRequiredStructEnabled())

	// By default validator reports errors using the Go struct field
	// name ("Title"), not your JSON key ("title"). Your API consumers
	// only know about "title" — they never see the Go struct. This
	// tells validator: when building an error, read the `json` tag
	// instead of the field name.
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// Struct validates s against its `validate:"..."` tags. Returns nil if
// valid, otherwise an *apperror.AppError ready to hand straight to
// response.Error — the handler doesn't need to know validator exists.
func Struct(s any) error {
	if err := v.Struct(s); err == nil {
		return nil
	} else {
		var fieldErrs validator.ValidationErrors
		if !errors.As(err, &fieldErrs) {
			// Not a validation failure in the expected shape (e.g. you
			// passed a non-struct) — treat as our fault, not the client's.
			return apperror.Internal(err)
		}

		details := make(map[string]string, len(fieldErrs))
		for _, fe := range fieldErrs {
			details[fe.Field()] = message(fe)
		}
		return apperror.Validation("validation failed", details)
	}
}

// message turns validator's tag ("required", "min", etc) into a
// sentence a client can actually read, instead of exposing validator's
// internal tag names raw.
func message(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "this field is required"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	default:
		return "invalid value"
	}
}
