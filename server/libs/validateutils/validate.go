package validateutils

import (
	"fmt"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type EnglishValidator struct {
	validator  *validator.Validate
	translator *ut.Translator
}

func NewEnglishValidator() *EnglishValidator {
	v := validator.New()
	trans := engErrTranslator(v)
	return &EnglishValidator{
		validator:  v,
		translator: trans,
	}
}

func (v *EnglishValidator) Struct(s interface{}) (errs []error) {
	err := v.validator.Struct(s)
	return translateError(err, v.translator)
}

func engErrTranslator(validate *validator.Validate) *ut.Translator {
	eng := en.New()
	uni := ut.New(eng, eng)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	return &trans
}

func translateError(err error, trans *ut.Translator) (errs []error) {
	if err == nil {
		return nil
	}
	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		translatedErr := fmt.Errorf(e.Translate(*trans))
		errs = append(errs, translatedErr)
	}
	return errs
}
