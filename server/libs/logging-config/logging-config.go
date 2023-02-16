package logging

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	common "github.com/vantoan19/Petifies/server/libs/common-utils"
)

type LoggingLevel string

// Data represents any types of data to log, data is a collection of
// key and corresponding logging data.
type Data map[string]interface{}

// LoggingMessage structures how a logging message looks like
type LoggingMessage struct {
	Time       time.Time
	PID        int // Process ID to collect related logs in distributed logging system.
	Level      LoggingLevel
	LoggerName string
	Message    string
	Data       Data
}

type Logger struct {
	Level LoggingLevel
	Name  string
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

const (
	defaultTimestampFormat    = "2023-01-30 15:00:00.000"
	defaultLogLevelTemplate   = "%8s"
	defaultLoggerNameTemplate = "%-30s"
)

const (
	Debug     LoggingLevel = "DEBUG"
	Info      LoggingLevel = "INFO"
	Notice    LoggingLevel = "NOTICE"
	Warning   LoggingLevel = "WARNING"
	Error     LoggingLevel = "ERROR"
	Critical  LoggingLevel = "CRITICAL"
	Alert     LoggingLevel = "ALERT"
	Emergency LoggingLevel = "EMERGENCY"
	Undefined LoggingLevel = "UNDEFINED"
)

var levelCodes = map[LoggingLevel]int{
	Emergency: 0,
	Alert:     1,
	Critical:  2,
	Error:     3,
	Warning:   4,
	Notice:    5,
	Info:      6,
	Debug:     7,
}

func New(name string) *Logger {
	var level LoggingLevel

	if common.IsDevEnv() {
		level = Debug
	} else {
		level = getEnvLogLevel()
	}

	return &Logger{
		Level: level,
		Name:  name,
	}
}

func (l *Logger) logString(level LoggingLevel, message string, data *Data) (string, error) {
	loggerLevel, ok := levelCodes[l.Level]
	if !ok {
		return "", errors.New("unrecognized level")
	}

	level_, ok := levelCodes[level]
	if !ok {
		return "", errors.New("unrecognized level")
	}

	if level_ > loggerLevel {
		return "", errors.New("message log level is too high")
	}

	var (
		timestamp  = time.Now().Format(defaultTimestampFormat)
		pID        = strconv.Itoa(os.Getpid())
		loggerName = fmt.Sprintf(defaultLoggerNameTemplate, l.Name)
	)
	data_, err := json.Marshal(*data)
	if err != nil {
		return "", err
	}

	s := "" +
		colorizeMessge(timestamp, Green) + "\t" +
		colorizeMessge(pID, Blue) + "\t" +
		colorizeLogLevel(level) + "\t" +
		colorizeMessge(loggerName, Purple) + "\t" +
		message + "\t" +
		"data = " + string(data_)

	return s, nil
}

func (l *Logger) log(level LoggingLevel, message string, data *Data) {
	msg, err := l.logString(level, message, data)
	if err != nil {
		return
	}
	fmt.Println(msg)
}

func (l *Logger) DebugData(message string, data Data) {
	l.log(Debug, message, &data)
}

func (l *Logger) InfoData(message string, data Data) {
	l.log(Info, message, &data)
}

func (l *Logger) NoticeData(message string, data Data) {
	l.log(Notice, message, &data)
}

func (l *Logger) WarningData(message string, data Data) {
	l.log(Warning, message, &data)
}

func (l *Logger) ErrorData(message string, data Data) {
	l.log(Error, message, &data)
}

func (l *Logger) CriticalData(message string, data Data) {
	l.log(Critical, message, &data)
}

func (l *Logger) AlertData(message string, data Data) {
	l.log(Alert, message, &data)
}

func (l *Logger) EmergencyData(message string, data Data) {
	l.log(Emergency, message, &data)
}

func (l *Logger) Debug(message string) {
	l.log(Debug, message, &Data{})
}

func (l *Logger) Info(message string) {
	l.log(Info, message, &Data{})
}

func (l *Logger) Notice(message string) {
	l.log(Notice, message, &Data{})
}

func (l *Logger) Warning(message string) {
	l.log(Warning, message, &Data{})
}

func (l *Logger) Error(message string) {
	l.log(Error, message, &Data{})
}

func (l *Logger) Critical(message string) {
	l.log(Critical, message, &Data{})
}

func (l *Logger) Alert(message string) {
	l.log(Alert, message, &Data{})
}

func (l *Logger) Emergency(message string) {
	l.log(Emergency, message, &Data{})
}

func parseLogLevel(logLevel string) (LoggingLevel, error) {
	logLevel_ := LoggingLevel(strings.ToUpper(logLevel))

	if _, ok := levelCodes[logLevel_]; ok {
		return logLevel_, nil
	}

	return Undefined, errors.New("unrecognized logging level: " + logLevel + ", using INFO as default.")
}

func getEnvLogLevel() LoggingLevel {
	if envLogLevel, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if level, err := parseLogLevel(envLogLevel); err != nil {
			fmt.Fprintln(os.Stderr, err)
			return Info
		} else {
			return level
		}
	}

	return Info
}

func colorizeMessge(msg string, color string) string {
	if !common.IsDevEnv() {
		return msg
	}

	return color + msg + Reset
}

func colorizeLogLevel(level LoggingLevel) string {
	level_ := fmt.Sprintf(defaultLogLevelTemplate, level)
	if !common.IsDevEnv() {
		return level_
	}

	switch level {
	case Debug:
		return Gray + level_ + Reset
	case Info:
		return Green + level_ + Reset
	case Notice:
		return Cyan + level_ + Reset
	case Warning:
		return Yellow + level_ + Reset
	case Error, Critical, Alert, Emergency:
		return Red + level_ + Reset
	default:
		return level_
	}
}
