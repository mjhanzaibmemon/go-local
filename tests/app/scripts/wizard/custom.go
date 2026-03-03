package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var isLetter = regexp.MustCompile(`^[a-zA-Z]+$`).MatchString
var isServiceName = regexp.MustCompile(`^([a-z][a-z0-9]*)(-[a-z0-9]+)*$`).MatchString
var customAppPkgName string

func validateCustomAppName(input string) error {
	if input == "" {
		return errors.New("value should not be empty")
	}

	if !isLetter(input) {
		return errors.New("invalid app name, should be CamelCase without spaces, dashes, underscored, e.g. SomeService")
	}

	return nil
}

func generateCustomApp() (string, error) {
	msg := fmt.Sprintf(`App name in CamelCase without spaces, dashes, underscored, e.g. MyApp`)
	appName := readConsole(msg, validateCustomAppName)
	customAppPkgName = strings.ToLower(appName)

	tracerServiceName, err := readStringInput(paramTracerServiceName, tracerMessage)
	if err != nil {
		panic(err)
	}

	templateParams[paramCustomAppName] = appName
	templateParams[paramCustomAppPkgName] = customAppPkgName
	templateParams[paramTracerServiceName] = tracerServiceName

	commonTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/common/", rootDir)
	appTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/%s/", rootDir, appCustom)

	info := fmt.Sprintf(`Before you run the app, do the following steps:
 - run %[1]s"make init"%[2]s
 - run %[1]s"make test"%[2]s
 - run %[1]s"make build"%[2]s
 - run %[1]s"make db_create"%[2]s

Entrypoint of this app is at %[1]scmd/%[3]s/main.go%[2]s.

Key files for later development:
 - entrypoint: %[1]scmd/%[3]s/main.go%[2]s
 - config file: %[1]s.env.%[3]s%[2]s
 - config struct: %[1]sinternal/config/%[3]s_cmd.go%[2]s
 - container configuration: %[1]sinternal/container/%[3]s_app.go%[2]s
`, green, reset, customAppPkgName)

	return info, generateApp(
		preProcessCustomAppTemplatePath,
		postGenerateEnableCommandsFn(map[string]string{
			"// <!-- BEGIN_COMMANDS -->": fmt.Sprintf("rootCmd.AddCommand(create%sCmd())", appName),
		}),
		commonTemplatesDir,
		appTemplatesDir,
	)
}

func preProcessCustomAppTemplatePath(path string) string {
	if strings.Contains(path, "example") {
		path = strings.Replace(path, "example", customAppPkgName, -1)
	}

	if strings.Contains(path, "cli-app-dir") {
		return strings.Replace(path, "cli-app-dir", serviceName, -1)
	}

	return path
}
