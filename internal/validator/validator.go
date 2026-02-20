package validator

import (
	"slices"
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

//Check if there are errors atall to determine display of errors at html level
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
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

