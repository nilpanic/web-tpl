package valid

import (
	"github.com/go-playground/validator/v10"
	"regexp"
)

var Mobile validator.Func = func(fl validator.FieldLevel) bool {
	pattern := `^1(3\d|4[5-9]|5[0-35-9]|6[567]|7[0-8]|8\d|9[0-35-9])\d{8}$`
	v, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	r := regexp.MustCompile(pattern)
	return r.MatchString(v)
}
