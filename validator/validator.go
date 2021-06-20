package utils

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"

	opa_core "github.com/MagalixTechnologies/opa-core"
	ut "github.com/go-playground/universal-translator"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

var Validator *validator.Validate
var Translator ut.Translator

const alphaNumericDashBracesRegexString = "^[a-zA-Z0-9-_() ]+$"
const alphaDashRegex = "^[a-zA-Z0-9-_ ]+$"
const whitespaceRegex = "^\\S+$"
const alphaNumericUnderscoresRegexString = "^[a-zA-Z0-9_]+$"

func init() {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatal("translator not found")
	}

	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	_ = v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})

	_ = v.RegisterTranslation("email", trans, func(ut ut.Translator) error {
		return ut.Add("email", "{0} must be a valid email", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("email", fe.Field())
		return t
	})

	// register custom validation: rfe(Required if Field is Equal to some value).
	_ = v.RegisterValidation(`rfe`, func(fl validator.FieldLevel) bool {
		param := strings.Split(fl.Param(), `:`)
		paramField := param[0]
		paramValue := param[1]

		if paramField == `` {
			return true
		}

		// param field reflect.Value.
		var paramFieldValue reflect.Value

		if fl.Parent().Kind() == reflect.Ptr {
			paramFieldValue = fl.Parent().Elem().FieldByName(paramField)
		} else {
			paramFieldValue = fl.Parent().FieldByName(paramField)
		}

		if isEq(paramFieldValue, paramValue) == false {
			return true
		}

		return hasValue(fl)
	})

	alphaNumericUnderscoresRegex := regexp.MustCompile(alphaNumericUnderscoresRegexString)
	_ = v.RegisterValidation(`alphanumericunderscores`, func(fl validator.FieldLevel) bool {
		return alphaNumericUnderscoresRegex.MatchString(fl.Field().String())
	})

	_ = v.RegisterTranslation("alphanumericunderscores", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumericunderscores", "{0} should only have alphanumeric, and underscores", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericunderscores", fe.Field())
		return t
	})

	// register custom validation: title(required if string should be alphanum and includes (), -, _).
	alphaNumericDashBraceRegex := regexp.MustCompile(alphaNumericDashBracesRegexString)
	_ = v.RegisterValidation(`alphanumdashbraces`, func(fl validator.FieldLevel) bool {
		return alphaNumericDashBraceRegex.MatchString(fl.Field().String())
	})

	_ = v.RegisterTranslation("alphanumdashbraces", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumdashbraces", "{0} should only have alphanumeric, (, ), -, _, spaces", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumdashbraces", fe.Field())
		return t
	})

	alphaDashRegex := regexp.MustCompile(alphaDashRegex)
	_ = v.RegisterValidation(`alphadash`, func(fl validator.FieldLevel) bool {
		return alphaDashRegex.MatchString(fl.Field().String())
	})

	_ = v.RegisterTranslation("alphadash", trans, func(ut ut.Translator) error {
		return ut.Add("alphadash", "{0} should only have alphanumeric, -, _", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphadash", fe.Field())
		return t
	})

	// register custom validation: does not include whitespaces.
	whitespaceRegex := regexp.MustCompile(whitespaceRegex)
	_ = v.RegisterValidation(`nowhitespace`, func(fl validator.FieldLevel) bool {
		fmt.Println(fl.Field().String())
		return whitespaceRegex.MatchString(fl.Field().String())
	})

	_ = v.RegisterTranslation("nowhitespace", trans, func(ut ut.Translator) error {
		return ut.Add("nowhitespace", "{0} should not have spaces", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("nowhitespace", fe.Field())
		return t
	})

	// register custom validation: rfe(Required if Field is Equal to some value).
	_ = v.RegisterValidation(`rego`, func(fl validator.FieldLevel) bool {

		field := fl.Field()
		switch field.Kind() {
		case reflect.String:
			_, err := opa_core.Parse(field.Interface().(string), "violation")
			if err != nil {
				return false
			}
		default:
			return false
		}
		return true
	})

	_ = v.RegisterTranslation("rfe", trans, func(ut ut.Translator) error {
		return ut.Add("rfe", "{0} is required if value of {1} is {2}", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		param := strings.Split(fe.Param(), ":")
		t, _ := ut.T("rfe", fe.Field(), param[0], param[1])
		return t
	})

	_ = v.RegisterTranslation("rego", trans, func(ut ut.Translator) error {
		return ut.Add("rego", "{0} invalid rego policy: {1}", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		_, err := opa_core.Parse(fe.Value().(string), "violation")
		t, _ := ut.T("rego", fe.Field(), err.Error())
		return t
	})

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	Validator = v
	Translator = trans
}

func hasValue(fl validator.FieldLevel) bool {
	return requireCheckFieldKind(fl, "")
}

func requireCheckFieldKind(fl validator.FieldLevel, param string) bool {
	field := fl.Field()
	if len(param) > 0 {
		if fl.Parent().Kind() == reflect.Ptr {
			field = fl.Parent().Elem().FieldByName(param)
		} else {
			field = fl.Parent().FieldByName(param)
		}
	}
	switch field.Kind() {
	case reflect.Slice, reflect.Map, reflect.Ptr, reflect.Interface, reflect.Chan, reflect.Func:
		return !field.IsNil()
	default:
		_, _, nullable := fl.ExtractType(field)
		if nullable && field.Interface() != nil {
			return true
		}
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}

func isEq(field reflect.Value, value string) bool {
	switch field.Kind() {

	case reflect.String:
		return field.String() == value

	case reflect.Slice, reflect.Map, reflect.Array:
		p := asInt(value)

		return int64(field.Len()) == p

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p := asInt(value)

		return field.Int() == p

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p := asUint(value)

		return field.Uint() == p

	case reflect.Float32, reflect.Float64:
		p := asFloat(value)

		return field.Float() == p
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}

func asInt(param string) int64 {

	i, err := strconv.ParseInt(param, 0, 64)
	panicIf(err)

	return i
}

func asUint(param string) uint64 {

	i, err := strconv.ParseUint(param, 0, 64)
	panicIf(err)

	return i
}

func asFloat(param string) float64 {

	i, err := strconv.ParseFloat(param, 64)
	panicIf(err)

	return i
}

func panicIf(err error) {
	if err != nil {
		panic(err.Error())
	}
}
