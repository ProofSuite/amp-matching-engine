package utils

import (
	"io"
	"os"
	"path"
	"runtime"

	logging "github.com/op/go-logging"
)

var Logger = NewLogger("main", "./logs/main.log")
var OperatorLogger = NewLogger("operator", "./logs/operator.log")
var EngineLogger = NewLogger("engine", "./logs/engine.log")
var APILogger = NewLogger("api", "./logs/api.log")
var RabbitLogger = NewLogger("rabbitmq", "./logs/rabbit.log")
var TerminalLogger = NewColoredLogger()

func NewLogger(module string, logFile string) *logging.Logger {
	_, fileName, _, _ := runtime.Caller(1)
	logDir := path.Join(path.Dir(fileName), "../logs/")
	mainLogFile := path.Join(path.Dir(fileName), "../logs/main.log")
	logFile = path.Join(path.Dir(fileName), "../", logFile)

	logger, err := logging.GetLogger("api")
	if err != nil {
		panic(err)
	}

	var format = logging.MustStringFormatter(
		`%{level:.4s} %{time:15:04:05} at %{shortpkg}/%{shortfile} in %{shortfunc}():%{message}`,
	)

	if _, err := os.Stat(logDir); os.IsNotExist(err) {
		os.Mkdir(logDir, os.ModePerm)
	}

	mainLog, err := os.OpenFile(mainLogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	log, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	writer := io.MultiWriter(os.Stdout, mainLog, log)
	backend := logging.NewLogBackend(writer, "", 0)

	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)

	logger.SetBackend(leveledBackend)
	return logger
}

func NewColoredLogger() *logging.Logger {
	logger, err := logging.GetLogger("colored")
	if err != nil {
		panic(err)
	}

	var format = logging.MustStringFormatter(
		`%{color}%{level:.4s} %{time:15:04:05} at %{shortpkg}/%{shortfile} in %{shortfunc}():%{color:reset} %{message}`,
	)

	writer := io.MultiWriter(os.Stdout)
	backend := logging.NewLogBackend(writer, "", 0)

	formattedBackend := logging.NewBackendFormatter(backend, format)
	leveledBackend := logging.AddModuleLevel(formattedBackend)

	logger.SetBackend(leveledBackend)
	return logger
}
