package forms

type errors map[string][]string

// Add adds an error message to the given field.
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message for the given field, or an empty string if there are no errors.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
