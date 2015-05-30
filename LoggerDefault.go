package jsonapi;

import (
    "log";
    "os";
    "io"
)

type LoggerDefault struct{
    Output io.Writer
    Logger *log.Logger
}

func NewLoggerDefault(output io.Writer) LoggerDefault {
    if output == nil {
        output = os.Stdout;
    }
    return LoggerDefault{
        Output: output,
        Logger: log.New(output,"",log.LstdFlags | log.Lmicroseconds | log.Lshortfile),
    }
}

func(l *LoggerDefault) Debugf(fmt string, args ...interface{}) {
    l.Logger.Printf("[DEBUG] "+fmt, args);
}
func(l *LoggerDefault) Infof(fmt string, args ...interface{}) {
    l.Logger.Printf("[INFO ] "+fmt, args);
}
func(l *LoggerDefault) Warnf(fmt string, args ...interface{}) {
    l.Logger.Printf("[WARN ] "+fmt, args);
}
func(l *LoggerDefault) Errorf(fmt string, args ...interface{}) {
    l.Logger.Printf("[ERROR] "+fmt, args);
}
func(l *LoggerDefault) Criticalf(fmt string, args ...interface{}) {
    l.Logger.Printf("[CRITICAL] "+fmt, args);
}
