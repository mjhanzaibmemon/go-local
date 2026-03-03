// Package translations loads and bundles localization message files (go-i18n) from a directory.
package translations

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v2"
)

// BundleConfig defines where translation files are located.
type BundleConfig struct {
	Dir string
}

// Bundle loads translations files into i18n.Bundle
func Bundle(c *BundleConfig) (*i18n.Bundle, error) {
	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	err := filepath.Walk(c.Dir, func(path string, info fs.FileInfo, _ error) error {
		if info.IsDir() {
			return nil
		}

		if _, err := bundle.LoadMessageFile(path); err != nil {
			return fmt.Errorf("failed to load %s translations file: %w", path, err)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to load translations: %w", err)
	}

	return bundle, nil
}
