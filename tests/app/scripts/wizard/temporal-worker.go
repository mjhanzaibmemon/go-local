package main

import "fmt"

func generateTemporalWorkerApp() (string, error) {
	tracerServiceName, err := readStringInput(paramTracerServiceName, tracerMessage)
	if err != nil {
		panic(err)
	}

	msg := fmt.Sprintf(`Temporal queue name which is going to be used in development.`)
	temporalQueueName := readConsole(msg, validateNonEmpty)

	redisCacheNamespace, err := readStringInput(paramRedisCacheNamespace, redisMessage)
	if err != nil {
		panic(err)
	}

	templateParams[paramTracerServiceName] = tracerServiceName
	templateParams[paramTemporalQueueName] = temporalQueueName
	templateParams[paramRedisCacheNamespace] = redisCacheNamespace

	commonTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/common/", rootDir)
	workerCommonTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/worker-common/", rootDir)
	temporalWorkerTemplatesDir := fmt.Sprintf("%s/scripts/wizard/templates/%s/", rootDir, appTemporalWorker)

	info := fmt.Sprintf(
		`Read more about the app at %[1]sdocs/template/app/temporal-worker.md%[2]s`,
		green,
		reset,
	)

	return info, generateApp(
		preProcessAppTemplatePath,
		postGenerateEnableCommandsFn(map[string]string{
			"// <!-- BEGIN_DEV_COMMANDS -->": "devCmdGroup.AddCommand(createDevTemporalActivityTriggerCmd())",
			"// <!-- BEGIN_COMMANDS -->":     "rootCmd.AddCommand(createTemporalWorkerCmd())",
		}),
		commonTemplatesDir,
		workerCommonTemplatesDir,
		temporalWorkerTemplatesDir,
	)
}
