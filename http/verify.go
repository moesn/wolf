package http

import (
	"errors"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtrans "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"strings"
)

var (
	trans    ut.Translator
	validate *validator.Validate
)


func VerifyInit() {
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ = uni.GetTranslator("zh")
	validate = validator.New()
	validate = validator.New()

	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("trans")

		if len(name)==0 {
			name = fld.Tag.Get("trans_")
		}
		if len(name)==0 {
			name = fld.Tag.Get("trans__")
		}
		if name == "-" {
			return ""
		}

		return "<" + name + ">"
	})

	zhtrans.RegisterDefaultTranslations(validate, trans)
}

func Verify(data interface{}) error {
	err := validate.Struct(data)

	if err != nil {
		var errMsgs []string

		for _, err := range err.(validator.ValidationErrors) {
			errMsgs = append(errMsgs, err.Translate(trans))
		}

		return errors.New(strings.Join(errMsgs, "，"))
	}
	return nil
}
