package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/manifoldco/promptui"
	"golang.org/x/tools/imports"
)

const (
	paramDBName              = "DBName"
	paramServiceName         = "ServiceName"
	paramModuleName          = "ModuleName"
	paramTracerServiceName   = "TracerServiceName"
	paramSQSQueueName        = "SQSQueueName"
	paramRedisCacheNamespace = "RedisCacheNamespace"
	paramCustomAppName       = "CustomAppName"
	paramCustomAppPkgName    = "CustomAppPkgName"
	paramTemporalQueueName   = "TemporalQueueName"

	dirPerm  = 0775
	filePerm = 0664

	appWeb            = "web"
	appCommandBus     = "command-bus"
	appEventBus       = "event-bus"
	appTemporalWorker = "temporal-worker"
	appCustom         = "custom"

	finishCreation = "exit"

	redisMessage  = "Redis cache namespace. It will be used as key prefix for SQS message deduplication keys stored in cache. Make sure it's unique between all eneba services"
	tracerMessage = "Tracer service name.This will be visible on jaeger UI. Make sure it's unique between all eneba services"
)

var (
	serviceName = ""
	rootDir     = ""

	templateParams = map[string]interface{}{}
	tmplRegEx      = regexp.MustCompile("\\.tmpl$")

	green = "\033[32m"
	reset = "\033[0m"

	commonInfo = fmt.Sprintf(`Read more about example applications at %[1]sdocs/template/apps/%[2]s

Happy development!`, green, reset)
)

type (
	inputValidateFn = func(input string) error
)

func main() {
	var err error
	rootDir, err = getRootDir()
	if err != nil {
		log.Fatalf("Failed to read working dir %v\n", err)
	}

	serviceName, err = readServiceName()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	moduleName, err := readModuleName()
	if err != nil {
		log.Fatalf("Prompt failed %v\n", err)
	}

	dbName := readConsole("Database name", validateNonEmpty)

	appTypes := []string{appWeb, appCommandBus, appEventBus, appTemporalWorker, appCustom}

loop:
	for {
		prompt := promptui.Select{
			Label: "Choose application type",
			Items: append(appTypes, finishCreation),
		}

		_, appType, err := prompt.Run()
		if err != nil {
			log.Fatalf("Prompt failed %v\n", err)
		}

		// common params
		templateParams[paramServiceName] = serviceName
		templateParams[paramModuleName] = moduleName
		templateParams[paramDBName] = dbName

		info := ""

		switch appType {
		case appWeb:
			info, err = generateWebApp()
			appTypes = removeElement(appTypes, appWeb)
		case appCommandBus:
			info, err = generateCommandBusApp()
			appTypes = removeElement(appTypes, appCommandBus)
		case appEventBus:
			info, err = generateEventBusApp()
			appTypes = removeElement(appTypes, appEventBus)
		case appTemporalWorker:
			info, err = generateTemporalWorkerApp()
			appTypes = removeElement(appTypes, appTemporalWorker)
		case appCustom:
			info, err = generateCustomApp()
		case finishCreation:
			break loop
		}

		if err != nil {
			log.Fatalln(err.Error())
		}
		fmt.Printf(`App is generated. %s`, info)
		fmt.Println()
	}

	fmt.Println()
	fmt.Println()
	fmt.Println(commonInfo)
}

func getRootDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to read working dir %v\n", err)
	}

	directoriesToCheck := []string{"cmd", "internal", "pkg", "translations"}

	for _, dir := range directoriesToCheck {
		pathToCheck := fmt.Sprintf("%s/%s", wd, dir)
		if _, err := os.Stat(pathToCheck); os.IsNotExist(err) {
			return "", fmt.Errorf(`invalid root dir: sub dir "%s" not found`, dir)
		}
	}

	return wd, nil
}

func readModuleName() (string, error) {
	path := fmt.Sprintf("%s/go.mod", rootDir)

	file, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", path, err)
	}

	defer file.Close()

	moduleName := ""
	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "module") {
			lines = append(lines, line)

			continue
		}

		if line == "" || strings.TrimLeft(line, "module ") == "change-me" {
			// module name is not set yet

			msg := fmt.Sprintf("Enter module name, e.g. gitlab.com/d5100/some-team/services/go/some-service")

			moduleName = readConsole(msg, validateNonEmpty)

			lines = append(lines, fmt.Sprintf("module %s", moduleName))
		} else {
			// module name is already set
			moduleName = strings.TrimPrefix(line, "module ")
			lines = append(lines, line)
		}
	}

	if moduleName == "" {
		return "", errors.New("invalid module name")
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), filePerm)
	if err != nil {
		return "", fmt.Errorf("failed to change module name: %w", err)
	}

	return moduleName, nil
}

func readStringInput(param, msg string) (string, error) {
	res, ok := templateParams[param]
	if !ok {
		return promptAndSaveInput(param, msg)
	}

	if res == "" {
		return promptAndSaveInput(param, msg)
	}

	return res.(string), nil
}

func promptAndSaveInput(path, msg string) (string, error) {
	input := readConsole(msg, validateNonEmpty)

	templateParams[path] = input

	return input, nil
}

func readServiceName() (string, error) {
	path := fmt.Sprintf("%s/cmd", rootDir)

	files, err := os.ReadDir(path)
	if err != nil {
		return "", fmt.Errorf("failed to open %s: %w", path, err)
	}

	directories := make([]string, 0)

	for _, f := range files {
		if f.IsDir() && f.Name() != "hello" {
			directories = append(directories, f.Name())
		}
	}

	if len(directories) > 1 {
		return "", fmt.Errorf("%s must not have more than 1 subdirectory", path)
	}

	name := ""

	if len(directories) == 1 {
		name = directories[0]

		err := validateServiceName(name)
		if err != nil {
			return "", err
		}

		return name, nil
	}

	msg := fmt.Sprintf("Enter service name, e.g. some-service")

	name = readConsole(msg, validateServiceName)

	serviceDir := fmt.Sprintf("%s/%s", path, name)

	err = os.Mkdir(serviceDir, dirPerm)
	if err != nil {
		return "", fmt.Errorf("failed to create service dir %s: %w", serviceDir, err)
	}

	return name, nil
}

func readConsole(msg string, validate inputValidateFn) string {
	prompt := promptui.Prompt{
		Label:    msg,
		Validate: validate,
	}

	input, err := prompt.Run()

	if err != nil {
		return readConsole(msg, validate)
	}

	return input
}

func generateApp(
	pathPreProcessFn func(path string) string,
	postGenerateFn func() error,
	templateDirs ...string,
) error {
	templates, err := collectTemplates(templateDirs...)
	if err != nil {
		return fmt.Errorf("failed to collect templates: %w", err)
	}

	for path, tmpl := range templates {
		path = pathPreProcessFn(path)
		fileName := fmt.Sprintf("%s/%s", rootDir, strings.TrimSuffix(path, ".tmpl"))

		if err := processTemplate(tmpl, fileName); err != nil {
			return fmt.Errorf("failed to generate code: %w", err)
		}
	}

	err = postGenerateFn()

	return err
}

func collectTemplates(templatesDirectories ...string) (map[string]*template.Template, error) {
	templates := map[string]*template.Template{}

	for _, templatesDir := range templatesDirectories {
		err := filepath.Walk(templatesDir, func(path string, info os.FileInfo, err error) error {
			if err == nil && tmplRegEx.MatchString(info.Name()) {
				key := strings.TrimPrefix(path, templatesDir)

				templateFile, errT := ioutil.ReadFile(path)
				if errT != nil {
					return fmt.Errorf("failed to read template file %s: %w", path, errT)
				}

				tmpl, errT := template.New(info.Name()).Parse(string(templateFile))
				if errT != nil {
					return fmt.Errorf("failed to pars template file %s: %w", path, errT)
				}

				_, ok := templates[key]
				if ok {
					return fmt.Errorf("duplicate template %s", path)
				}

				templates[key] = tmpl
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	return templates, nil
}

func processTemplate(tmpl *template.Template, fileName string) error {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		// ensure directory exists
		targetDir := filepath.Dir(fileName)

		err := os.MkdirAll(targetDir, dirPerm)
		if err != nil {
			return fmt.Errorf("failed to create directory %s: %w", targetDir, err)
		}

		buf := bytes.NewBuffer([]byte{})

		err = tmpl.Execute(buf, templateParams)
		if err != nil {
			return fmt.Errorf("failed to execute template: %w", err)
		}

		content := buf.Bytes()

		if strings.HasSuffix(fileName, ".go") {
			res, err := imports.Process(fileName, buf.Bytes(), nil)
			if err != nil {
				return fmt.Errorf("failed to process imports: %w", err)
			}

			content = res
		}

		return ioutil.WriteFile(fileName, content, filePerm)
	}

	return nil
}

func validateNonEmpty(input string) error {
	if input == "" {
		return errors.New("value should not be empty")
	}

	return nil
}

func validateServiceName(input string) error {
	if input == "" {
		return errors.New("value should not be empty")
	}

	if !isServiceName(input) {
		return errors.New("invalid service name, should be lower-kebab-case without spaces, dashes, underscored, e.g. some-service")
	}

	return nil
}

func preProcessAppTemplatePath(path string) string {
	if strings.Contains(path, "cli-app-dir") {
		return strings.Replace(path, "cli-app-dir", serviceName, -1)
	}

	return path
}

func removeElement(l []string, item string) []string {
	for i, value := range l {
		if value == item {
			return append(l[:i], l[i+1:]...)
		}
	}

	return l
}
