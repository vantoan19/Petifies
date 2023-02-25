package validateutils

import (
	"fmt"
	"sync"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

var lock = &sync.Mutex{}

type englishValidator struct {
	validator  *validator.Validate
	translator *ut.Translator
}

var englishValidatorInstance *englishValidator

func GetEnglishValidatorInstance() *englishValidator {
	if englishValidatorInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if englishValidatorInstance == nil {
			v := validator.New()
			trans := engErrTranslator(v)
			englishValidatorInstance = &englishValidator{
				validator:  v,
				translator: trans,
			}
		}
	}

	return englishValidatorInstance
}

func (v *englishValidator) Struct(s interface{}) (errs []error) {
	err := v.validator.Struct(s)
	return translateError(err, v.translator)
}

func (v *englishValidator) Var(va interface{}, rules string) error {
	err := v.validator.Var(va, rules)
	return err
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
