// Package settings provides initialization helpers for the settings service router.
package settings

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/csgot/helis-market-base-client-go/pkg/client"
	settingsclient "bitbucket.org/csgot/helis-market-settings-client-go/pkg/client"
	"bitbucket.org/csgot/helis-market-settings-client-go/pkg/stn"
	"github.com/sirupsen/logrus"
)

// NewRouter creates a settings router backed by the settings service reachable at serviceURI.
func NewRouter(serviceURI string, httpClient *http.Client, logger *logrus.Logger) (*stn.Router, error) {
	path := "/graphql/"

	settingsClient := settingsclient.New(&settingsclient.Config{
		Executor: client.NewGQLExecutor(&client.GQLExecutorConfig{
			BaseURL:    serviceURI,
			Path:       &path,
			HTTPClient: httpClient,
		}),
	})

	settingsLoader, err := stn.NewCachedLoader(settingsClient, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize settings loader: %w", err)
	}

	settingsRouter, err := stn.NewRouter(settingsLoader, logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create settings router: %w", err)
	}

	return settingsRouter, nil
}
