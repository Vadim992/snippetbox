package validator

import (
	"strings"
	"unicode/utf8"
)

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

func (v *Validator) AddFileError(key, message string) {
	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string, 3)
	}

	if _, exist := v.FieldErrors[key]; !exist {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFileError(key, message)
	}
}

func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func PermittedInt(value int, permittedValues ...int) bool {
	for _, item := range permittedValues {
		if value == item {
			return true
		}
	}
	return false
}
