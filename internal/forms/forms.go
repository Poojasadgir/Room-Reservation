package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form creates a custom form struct and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// Valid returns true if there are no errors in the form, false otherwise.
func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

// New creates a new Form instance with the provided data and an empty errors map.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Required checks if the specified fields in the form are not blank.
// If any of the fields are blank, an error message is added to the form's Errors field.
func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has returns true if the form field with the given name has a non-empty value.
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	return x != ""
}

// MinLength checks if the value of the given field is at least of a certain length.
// If the value is shorter than the given length, an error message is added to the form errors.
// Returns true if the value is at least of the given length, false otherwise.
func (f *Form) MinLength(field string, length int) bool {
	x := f.Get(field)
	if len(x) < length {
		f.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}
	return true
}

// IsEmail checks if the value of the given field is a valid email address.
// If the value is not a valid email address, an error message is added to the form errors.
func (f *Form) IsEmail(field string) {
	if !govalidator.IsEmail(f.Get(field)) {
		f.Errors.Add(field, "Invalid email address")
	}
}
