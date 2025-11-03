package utils

import (
    "fmt"
    "log"
    "os"
    "sync"
)

type Level uint8

const (
    INFO Level = iota
    WARNING
    ERROR
    FATAL
)

func (l Level) String() string {
    switch l {
    case INFO: return "INFO"
    case WARNING: return "WARNING"
    case ERROR: return "ERROR"
    case FATAL: return "FATAL"
    default: return "UNKNOWN"
    }
}

type Logger struct {
    mu  sync.Mutex
    log *log.Logger
}

func NewLogger() *Logger {
    return &Logger{
        log: log.New(os.Stdout, "", log.Ldate|log.Ltime),
    }
}

func (l *Logger) logf(level Level, msg string, v ...any) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    oldPrefix := l.log.Prefix()
    l.log.SetPrefix(fmt.Sprintf("[%s] - ", level))
    defer l.log.SetPrefix(oldPrefix)
    
    l.log.Printf(msg, v...)
}

func (l *Logger) logln(level Level, msg string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    oldPrefix := l.log.Prefix()
    l.log.SetPrefix(fmt.Sprintf("[%s] - ", level))
    defer l.log.SetPrefix(oldPrefix)
    
    l.log.Println(msg)
}

// Métodos específicos (corrigidos)
func (l *Logger) Infof(msg string, v ...any) {
    l.logf(INFO, msg, v...)
}

func (l *Logger) Info(msg string) {
    l.logln(INFO, msg)
}

func (l *Logger) Warningf(msg string, v ...any) {
    l.logf(WARNING, msg, v...)
}

func (l *Logger) Warning(msg string) {
    l.logln(WARNING, msg)
}

func (l *Logger) Errorf(msg string, v ...any) {
    l.logf(ERROR, msg, v...)
}

func (l *Logger) Error(msg string) {
    l.logln(ERROR, msg)
}

func (l *Logger) Fatalf(msg string, v ...any) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    l.log.SetPrefix(fmt.Sprintf("%s - ", FATAL))
    l.log.Fatalf(msg, v...)
}

func (l *Logger) Fatal(msg string) {
    l.mu.Lock()
    defer l.mu.Unlock()
    
    l.log.SetPrefix(fmt.Sprintf("%s - ", FATAL))
    l.log.Fatalln(msg)
}