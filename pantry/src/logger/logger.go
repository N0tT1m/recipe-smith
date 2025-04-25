package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var (
	// Define log levels
	LevelInfo    = "INFO"
	LevelWarning = "WARNING"
	LevelError   = "ERROR"

	// Mutex for thread-safe logging
	logMutex sync.Mutex

	// Log file
	logFile *os.File

	// Logger instances
	infoLogger    *log.Logger
	warningLogger *log.Logger
	errorLogger   *log.Logger

	// Initialize flag
	initialized bool
)

// init initializes the logger
func init() {
	setupLoggers()
}

// setupLoggers configures the loggers
func setupLoggers() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if initialized {
		return
	}

	// Create logs directory if it doesn't exist
	err := os.MkdirAll("logs", 0755)
	if err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Create log file with timestamp in filename
	timestamp := time.Now().Format("2006-01-02")
	filename := fmt.Sprintf("logs/recipe_crawler_%s.log", timestamp)

	logFile, err = os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Create multi-writer to log to both file and stdout
	multiWriter := MultiWriter{
		file:   logFile,
		stdout: os.Stdout,
	}

	// Initialize loggers
	infoLogger = log.New(multiWriter, fmt.Sprintf("[%s] ", LevelInfo), log.Ldate|log.Ltime)
	warningLogger = log.New(multiWriter, fmt.Sprintf("[%s] ", LevelWarning), log.Ldate|log.Ltime)
	errorLogger = log.New(multiWriter, fmt.Sprintf("[%s] ", LevelError), log.Ldate|log.Ltime)

	initialized = true
}

// MultiWriter implements io.Writer to write to both file and stdout
type MultiWriter struct {
	file   *os.File
	stdout *os.File
}

// Write implements io.Writer
func (mw MultiWriter) Write(p []byte) (n int, err error) {
	n, err = mw.file.Write(p)
	if err != nil {
		return n, err
	}

	n, err = mw.stdout.Write(p)
	return n, err
}

// WriteInfo logs an info message
func WriteInfo(message string) {
	if !initialized {
		setupLoggers()
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	infoLogger.Println(message)
}

// WriteWarning logs a warning message
func WriteWarning(message string) {
	if !initialized {
		setupLoggers()
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	warningLogger.Println(message)
}

// WriteError logs an error message
func WriteError(message string) {
	if !initialized {
		setupLoggers()
	}

	logMutex.Lock()
	defer logMutex.Unlock()

	errorLogger.Println(message)
}

// Flush ensures all logs are written
func Flush() {
	if logFile != nil {
		logFile.Sync()
	}
}

// Close closes the log file
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

// LogWithContext logs a message with context information
func LogWithContext(level, context, message string) {
	formatted := fmt.Sprintf("[%s] %s", context, message)

	switch level {
	case LevelInfo:
		WriteInfo(formatted)
	case LevelWarning:
		WriteWarning(formatted)
	case LevelError:
		WriteError(formatted)
	default:
		WriteInfo(formatted)
	}
}
