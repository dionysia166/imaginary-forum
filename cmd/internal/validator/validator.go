package validator

import (
    "unicode"
    "regexp"
	"strings"
	"unicode/utf8"
)

type Validator struct {
    NonFieldErrors []string
    FieldErrors map[string]string
}

// Valid() returns true if the FieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
    return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

// AddFieldError() adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) AddFieldError(key, message string) {
    if v.FieldErrors == nil {
        v.FieldErrors = make(map[string]string)
    }

    if _, exists := v.FieldErrors[key]; !exists {
        v.FieldErrors[key] = message
    }
}

// AddNonFieldError() adds an error message to the NonFieldErrors slice.
func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

// CheckField() adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
    if !ok {
        v.AddFieldError(key, message)
    }
}

// NotBlank() returns true if a value is not an empty string.
func NotBlank(value string) bool {
    return strings.TrimSpace(value) != ""
}

// MaxChars() returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
    return utf8.RuneCountInString(value) <= n
}

// MinChars() returns true if a value contains no less than n characters.
func MinChars(value string, n int) bool {
    return utf8.RuneCountInString(value) >= n
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// Matches() returns true if a value matches a regular expression.
func Matches(value string, rx *regexp.Regexp) bool {
    return rx.MatchString(value)
}

// UpperCase() returns true if a value contains at least one uppercase letter.
func UpperCase(value string) bool {
    for _, r := range value {
        if unicode.IsUpper(r) {
            return true
        }
    }
    return false
}

// ContainsNumber() returns true if a value contains at least one uppercase letter.
func ContainsNumber(value string) bool {
    for _, r := range value {
        if unicode.IsNumber(r) {
            return true
        }
    }
    return false
}