package validate

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// V is the validator single instance.
//
// It is a singleton so to cache the structs info.
var V *validator.Validate

// T is the universal translator for validatiors.
var T ut.Translator

func init() {
	V = validator.New()

	// Register a function to get the field name from "name" tags.
	V.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("name"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	V.RegisterValidation("semver", isSemVer)

	eng := en.New()
	uni := ut.New(eng, eng)
	T, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(V, T)

	V.RegisterTranslation(
		"semver",
		T,
		func(ut ut.Translator) error {
			return ut.Add("semver", "{0} must be a semver-ish string", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(fe.Tag(), fe.Field())

			return t
		},
	)
}
