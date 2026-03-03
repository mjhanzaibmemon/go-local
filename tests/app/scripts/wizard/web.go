package main

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func generateWebApp() (string, error) {
	tracerServiceName, err := readStringInput(paramTracerServiceName, tracerMessage)
	if err != nil {
		panic(err)
	}

	templateParams[paramTracerServiceName] = tracerServiceName

	commonTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/common/", rootDir)
	webTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/%s/", rootDir, appWeb)

	info := fmt.Sprintf(`Read more about the app at %[1]sdocs/template/app/web.md%[2]s`, green, reset)

	return info, generateApp(
		preProcessAppTemplatePath,
		postGenerateEnableCommandsFn(map[string]string{
			"// <!-- BEGIN_COMMANDS -->": "rootCmd.AddCommand(createWebCmd())",
		}),
		commonTemplatesDir,
		webTemplatesDir,
	)
}

func postGenerateEnableCommandsFn(cfg map[string]string) func() error {
	return func() error {
		return enableCommands(cfg)
	}
}

func enableCommands(cfg map[string]string) error {
	path := fmt.Sprintf("%s/cmd/%s/main.go", rootDir, serviceName)

	input, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Errorf("failed to read file %s: %w", path, err)
	}

	lines := strings.Split(string(input), "\n")

	newLines := make([]string, 0)

	for _, line := range lines {
		for match, newLine := range cfg {
			if strings.Contains(line, match) {
				newLines = append(newLines, fmt.Sprintf("\t%s", newLine))
			}
		}

		newLines = append(newLines, line)
	}

	output := strings.Join(newLines, "\n")

	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {
		return fmt.Errorf("failed to update %s: %w", path, err)
	}

	return nil
}
