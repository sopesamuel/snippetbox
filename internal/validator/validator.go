package validator

import (
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validator struct {
	NonFieldErrors []string
	FieldErrors map[string]string
}

//Check if there are errors atall to determine display of errors at html level
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0 && len(v.NonFieldErrors) == 0
}

func (v *Validator) AddNonField(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddField(key string, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	_, exists := v.FieldErrors[key]
	if !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key string, message string) {
	if !ok {
		v.AddField(key, message)
	}
}

//Necessary checks blank, wrong expiry value, more than 100
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChar(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedValue[T comparable](value T, permittedValue ...T) bool {
	return slices.Contains(permittedValue, value)
}

func MinChars(value string, n int) bool{
	return utf8.RuneCountInString(value) >= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}
func MatchOthers(actual string, expected string) bool {
	return actual == expected
}