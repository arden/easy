package logger

import (
    "github.com/gogf/gf/frame/g"
    "github.com/gogf/gf/os/glog"
)

var (

    errorLogPattern = "error-{Ymd}.log"
    errorLogger = g.Log().File(errorLogPattern)

    fatalLogPattern = "fatal-{Ymd}.log"
    fatalLogger = g.Log().File(fatalLogPattern)

    panicLogPattern = "panic-{Ymd}.log"
    panicLogger = g.Log().File(panicLogPattern)

    infoLogPattern = "info-{Ymd}.log"
    infoLogger = g.Log().File(infoLogPattern)

    debugLogPattern = "debug-{Ymd}.log"
    debugLogger = g.Log().File(debugLogPattern)

    noticeLogPattern = "notice-{Ymd}.log"
    noticeLogger = g.Log().File(noticeLogPattern)

    warningLogPattern = "warning-{Ymd}.log"
    warningLogger = g.Log().File(warningLogPattern)

    criticalLogPattern = "critical-{Ymd}.log"
    criticalLogger = g.Log().File(criticalLogPattern)
)

type Logger struct {
    name string
    *glog.Logger
}

func New(names...string) *Logger {
    name := ""
    if len(names) > 0 {
        name = names[0]
    }
    logger := &Logger{
        name: name,
    }
    return logger
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func (l *Logger) Fatal(v ...interface{}) {
    fatalLogger.Cat(l.name).Fatal(v)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func (l *Logger) Fatalf(format string, v ...interface{}) {
    fatalLogger.Cat(l.name).Fatalf(format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func (l *Logger) Panic(v ...interface{}) {
    panicLogger.Cat(l.name).Panic(v...)
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func (l *Logger) Panicf(format string, v ...interface{}) {
    panicLogger.Cat(l.name).Panicf(format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func (l *Logger) Info(v ...interface{}) {
    infoLogger.Cat(l.name).Info(v...)
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func (l *Logger) Infof(format string, v ...interface{}) {
    infoLogger.Cat(l.name).Infof(format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func (l *Logger) Debug(v ...interface{}) {
    debugLogger.Cat(l.name).Debug(v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func (l *Logger) Debugf(format string, v ...interface{}) {
    debugLogger.Cat(l.name).Debugf(format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Notice(v ...interface{}) {
    noticeLogger.Cat(l.name).Notice(v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Noticef(format string, v ...interface{}) {
    noticeLogger.Cat(l.name).Noticef(format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Warning(v ...interface{}) {
    warningLogger.Cat(l.name).Warning(v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Warningf(format string, v ...interface{}) {
    warningLogger.Cat(l.name).Warningf(format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Error(v ...interface{}) {
    errorLogger.Cat(l.name).Error(v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Errorf(format string, v ...interface{}) {
    errorLogger.Cat(l.name).Errorf(format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Critical(v ...interface{}) {
    criticalLogger.Cat(l.name).Critical(v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Criticalf(format string, v ...interface{}) {
    criticalLogger.Cat(l.name).Criticalf(format, v...)
}

