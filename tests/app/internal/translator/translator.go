// Package translator provides helpers to translate messages and rich errors using go-i18n.
package translator

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Translator wraps an i18n.Bundle and provides helpers for translating strings and errors.
type Translator struct {
	bundle *i18n.Bundle
}

// New constructs a Translator using the provided i18n.Bundle.
func New(bundle *i18n.Bundle) *Translator {
	return &Translator{
		bundle: bundle,
	}
}

// Trans translates a message by ID for the given locale using optional template params.
func (t *Translator) Trans(locale, id string, params interface{}) string {
	localizer := i18n.NewLocalizer(t.bundle, locale)

	translated, _ := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    id,
		TemplateData: params,
	})

	return translated
}

// TransErr tries to translate a rich error (e.g., validator errors). Returns translated message and a flag.
func (t *Translator) TransErr(locale string, err error) (string, bool) {
	localizer := i18n.NewLocalizer(t.bundle, locale)

	var (
		validatorErrs validator.ValidationErrors
		// TODO: add app specific translatable errors here
	)

	if errors.As(err, &validatorErrs) {
		return translateValidateErrors(localizer, validatorErrs), true
	}

	// TODO: handle app specific translatable errors here

	return "", false
}

func translateValidateErrors(localizer *i18n.Localizer, validatorErrs validator.ValidationErrors) string {
	buff := bytes.NewBufferString("")

	for _, validatorErr := range validatorErrs {
		var pluralCount interface{}

		pluralCount = validatorErr.Param()
		if pluralCount == "" {
			pluralCount = nil
		}

		msg, _ := localizer.Localize(&i18n.LocalizeConfig{
			MessageID: fmt.Sprintf("validation.%s", validatorErr.Tag()),
			TemplateData: map[string]interface{}{
				"Field": validatorErr.Field(),
				"Param": validatorErr.Param(),
				"Value": validatorErr.Value(),
			},
			PluralCount: pluralCount,
		})

		buff.WriteString(msg)
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}
