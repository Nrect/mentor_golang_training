package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	colorReset         = "\033[0m"
	colorBrightRed     = "\033[91m"
	colorBrightGreen   = "\033[92m"
	colorBrightYellow  = "\033[93m"
	colorBrightBlue    = "\033[94m"
	colorBrightMagenta = "\033[95m"
	colorBrightCyan    = "\033[96m"
)

type Logger interface {
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	SetOutput(output *os.File)
}

type logger struct {
	debugEnabled bool
}

func New(debug bool) *logger {
	return &logger{debugEnabled: debug}
}

func (l *logger) Info(args ...interface{}) {
	log.Println(colorBrightGreen + "[INFO] " + fmt.Sprint(args...) + colorReset)
}

func (l *logger) Infof(format string, args ...interface{}) {
	log.Printf(colorBrightGreen+"[INFO] "+format+colorReset, args...)
}

func (l *logger) Debug(args ...interface{}) {
	if l.debugEnabled {
		log.Println(colorBrightCyan + "[DEBUG] " + fmt.Sprint(args...) + colorReset)
	}
}

func (l *logger) Debugf(format string, args ...interface{}) {
	if l.debugEnabled {
		log.Printf(colorBrightCyan+"[DEBUG] "+format+colorReset, args...)
	}
}

func (l *logger) Error(args ...interface{}) {
	log.Println(colorBrightRed + "[ERROR] " + fmt.Sprint(args...) + colorReset)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	log.Printf(colorBrightRed+"[ERROR] "+format+colorReset, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	log.Fatalln(colorBrightMagenta + "[FATAL] " + fmt.Sprint(args...) + colorReset)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	log.Fatalf(colorBrightMagenta+"[FATAL] "+format+colorReset, args...)
}

func (l *logger) Panic(args ...interface{}) {
	log.Panicln(colorBrightYellow + "[PANIC] " + fmt.Sprint(args...) + colorReset)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	log.Panicf(colorBrightYellow+"[PANIC] "+format+colorReset, args...)
}

func (l *logger) SetOutput(output *os.File) {
	log.SetOutput(output)
}
