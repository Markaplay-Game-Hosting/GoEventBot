package validator

import (
	duration "github.com/channelmeter/iso8601duration"
	"github.com/teambition/rrule-go"
	"regexp"
)

var (
	EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}

func (v *Validator) IsValidRRule(s string) {
	// A simple check to see if the string starts with "RRULE:"
	_, err := rrule.StrToRRule(s)
	if err != nil {
		v.AddError("rrule", "must be a valid RRULE format: "+err.Error())
	}
}

func (v *Validator) IsValidDurationRule(s string) bool {
	_, err := duration.FromString(s)
	if err != nil {
		v.AddError("duration", "must be a valid ISO 8601 duration format: "+err.Error())
		return false
	}
	return true
}
