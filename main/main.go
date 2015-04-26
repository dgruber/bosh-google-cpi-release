package main

import (
	"flag"
	"os"

	boshlog "github.com/cloudfoundry/bosh-agent/logger"
	boshsys "github.com/cloudfoundry/bosh-agent/system"
	boshuuid "github.com/cloudfoundry/bosh-agent/uuid"

	bgcaction "github.com/frodenas/bosh-google-cpi/action"
	bgcdisp "github.com/frodenas/bosh-google-cpi/api/dispatcher"
	bgctrans "github.com/frodenas/bosh-google-cpi/api/transport"

	gclient "github.com/frodenas/bosh-google-cpi/google/client"
)

const mainLogTag = "main"

var (
	configPathOpt = flag.String("configPath", "", "Path to configuration file")
)

func main() {
	logger, fs, cmdRunner, uuidGen := basicDeps()

	defer logger.HandlePanic("Main")

	flag.Parse()

	config, err := NewConfigFromPath(*configPathOpt, fs)
	if err != nil {
		logger.Error(mainLogTag, "Loading config - %s", err.Error())
		os.Exit(1)
	}

	dispatcher, err := buildDispatcher(config, logger, fs, cmdRunner, uuidGen)
	if err != nil {
		logger.Error(mainLogTag, "Building Dispatcher - %s", err)
		os.Exit(1)
	}

	cli := bgctrans.NewCLI(os.Stdin, os.Stdout, dispatcher, logger)

	err = cli.ServeOnce()
	if err != nil {
		logger.Error(mainLogTag, "Serving once %s", err)
		os.Exit(1)
	}
}

func basicDeps() (boshlog.Logger, boshsys.FileSystem, boshsys.CmdRunner, boshuuid.Generator) {
	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr, os.Stderr)

	fs := boshsys.NewOsFileSystem(logger)

	cmdRunner := boshsys.NewExecCmdRunner(logger)

	uuidGen := boshuuid.NewGenerator()

	return logger, fs, cmdRunner, uuidGen
}

func buildDispatcher(
	config Config,
	logger boshlog.Logger,
	fs boshsys.FileSystem,
	cmdRunner boshsys.CmdRunner,
	uuidGen boshuuid.Generator,
) (bgcdisp.Dispatcher, error) {
	googleClient, err := gclient.NewGoogleClient(config.Google.Project, config.Google.JsonKey, config.Google.DefaultZone, config.Google.AccessKeyId, config.Google.SecretAccessKey)
	if err != nil {
		return nil, err
	}

	actionFactory := bgcaction.NewConcreteFactory(
		googleClient,
		uuidGen,
		config.Actions,
		logger,
	)

	caller := bgcdisp.NewJSONCaller()

	return bgcdisp.NewJSON(actionFactory, caller, logger), nil
}
