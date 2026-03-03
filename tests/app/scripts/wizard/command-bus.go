package main

import (
	"fmt"
)

const (
	sqsMessage = "SQS queue name which is going to be used in development."
)

func generateCommandBusApp() (string, error) {
	tracerServiceName, err := readStringInput(paramTracerServiceName, tracerMessage)
	if err != nil {
		panic(err)
	}

	sqsQueueName, err := readStringInput(paramSQSQueueName, sqsMessage)
	if err != nil {
		panic(err)
	}

	redisCacheNamespace, err := readStringInput(paramRedisCacheNamespace, redisMessage)
	if err != nil {
		panic(err)
	}

	templateParams[paramTracerServiceName] = tracerServiceName
	templateParams[paramSQSQueueName] = sqsQueueName
	templateParams[paramRedisCacheNamespace] = redisCacheNamespace

	commonTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/common/", rootDir)
	busCommonTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/bus-common/", rootDir)
	workerCommonTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/worker-common/", rootDir)
	commandBusTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/%s/", rootDir, appCommandBus)

	info := fmt.Sprintf(
		`Read more about the app at %[1]sdocs/template/app/command-bus.md%[2]s`,
		green,
		reset,
	)

	err = generateApp(
		preProcessAppTemplatePath,
		postGenerateEnableCommandsFn(map[string]string{
			"// <!-- BEGIN_COMMANDS -->": "rootCmd.AddCommand(createCommandBusCmd())",
		}),
		commonTemplatesDir,
		busCommonTemplatesDir,
		workerCommonTemplatesDir,
		commandBusTemplatesDir,
	)
	if err != nil {

	}

	return info, err
}
