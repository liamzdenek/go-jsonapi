package jsonapi;

import (
    "os";
    "io"
    "runtime"
    "time"
    "fmt"
    "strings"
)

type LoggerDefault struct{
    Output io.Writer
}

func NewLoggerDefault(output io.Writer) *LoggerDefault {
    if output == nil {
        output = os.Stdout;
    }
    return &LoggerDefault{
        Output: output,
    }
}

func(l *LoggerDefault) GetCaller(depth int) (string) {
    _, file, line, ok := runtime.Caller(depth+1);
    if !ok {
        file = "???";
        line = 0;
    } else {
        file_parts := strings.Split(file, "/");
        file = file_parts[len(file_parts)-1];
    }
    return fmt.Sprintf("%s:%d", file, line);
}

func(l *LoggerDefault) GetTime() string {
    return time.Now().Local().Format(
        "2006/01/02 15:04:05.000000",
    );
}

func(l *LoggerDefault) PrepareArgs(args []interface{}) []interface{} {
    return append([]interface{}{
        l.GetTime(),
        //l.GetCaller(2),
    }, args...);
}

func(l *LoggerDefault) Debugf(format string, args ...interface{}) {
    l.Output.Write(
        []byte(fmt.Sprintf("%s [DEBUG] "+format, l.PrepareArgs(args)...)),
    );
}
func(l *LoggerDefault) Infof(format string, args ...interface{}) {
    l.Output.Write(
        []byte(fmt.Sprintf("%s [INFO ] "+format, l.PrepareArgs(args)...)),
    );
}
func(l *LoggerDefault) Warnf(format string, args ...interface{}) {
    l.Output.Write(
        []byte(fmt.Sprintf("%s [WARN ] "+format, l.PrepareArgs(args)...)),
    );
}
func(l *LoggerDefault) Errorf(format string, args ...interface{}) {
    l.Output.Write(
        []byte(fmt.Sprintf("%s [ERROR] "+format, l.PrepareArgs(args)...)),
    );
}
func(l *LoggerDefault) Criticalf(format string, args ...interface{}) {
    l.Output.Write(
        []byte(fmt.Sprintf("%s [CRITICAL] "+format, l.PrepareArgs(args)...)),
    );
}
