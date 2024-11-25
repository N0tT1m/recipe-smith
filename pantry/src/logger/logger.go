package logger

import (
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type LogSeverity int

const (
	DEBUG   LogSeverity = iota
	ERROR   LogSeverity = iota
	WARNING LogSeverity = iota
	INFO    LogSeverity = iota
)

func GetSeverityString(logType LogSeverity) string {
	switch logType {
	case DEBUG:
		return "[DEBUG]"
	case ERROR:
		return "[ERROR]"
	case WARNING:
		return "[WARNING]"
	case INFO:
		return "[INFO]"
	default:
		return "[INFO]"
	}
}

func write(logType LogSeverity, message interface{}) {
	f, err := os.OpenFile("high-seas.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer f.Close()

	mw := io.MultiWriter(os.Stdout, f)
	log.SetOutput(mw)
	log.Printf("%s: %s\n", GetSeverityString(logType), message)
}

func WriteError(message string, err error) {
	write(ERROR, fmt.Sprintf("%s %s", message, err))
}

func WriteWarning(message string) {
	write(WARNING, message)
}

func WriteInfo(message interface{}) {
	write(INFO, message)
}

func WriteFatal(errMsg string, err error) {
	f, err := os.OpenFile("high-seas.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {

	}

	log.SetOutput(f)

	log.SetFormatter(&log.JSONFormatter{})

	log.SetLevel(log.InfoLevel)

	log.WithFields(log.Fields{
		"Reason for Error": errMsg,
		"Error":            err,
	}).Fatal("Can't finish running program will crash")
}

func WriteCMDInfo(cmd string, output string) {
	f, err := os.OpenFile("high-seas.log", os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		WriteError("Couldn't open file.", err)
	}

	log.SetOutput(f)

	log.SetFormatter(&log.JSONFormatter{})

	log.SetLevel(log.InfoLevel)

	log.WithFields(log.Fields{
		"Command":        cmd,
		"Command Output": output,
	}).Info("Command output.")
}
