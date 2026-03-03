package main

import (
	"fmt"
)

func generateEventBusApp() (string, error) {
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
	eventBusTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/%s/", rootDir, appEventBus)

	info := fmt.Sprintf(
		`Read more about the app at %[1]sdocs/template/app/event-bus.md%[2]s`,
		green,
		reset,
	)

	return info, generateApp(
		preProcessAppTemplatePath,
		postGenerateEnableCommandsFn(map[string]string{
			"// <!-- BEGIN_COMMANDS -->": "rootCmd.AddCommand(createEventBusCmd())",
		}),
		commonTemplatesDir,
		busCommonTemplatesDir,
		workerCommonTemplatesDir,
		eventBusTemplatesDir,
	)
}
