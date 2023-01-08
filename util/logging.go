package util

import "fmt"

type UtilLogger struct {
}

func NewLogger() *UtilLogger {
	return &UtilLogger{}
}

func (l *UtilLogger) Tracef(format string, values ...interface{}) {
	fmt.Printf("[Tracef] %s", fmt.Sprintf(format, values...))
}
func (l *UtilLogger) Debugf(format string, values ...interface{}) {
	fmt.Printf("[Debugf] %s", fmt.Sprintf(format, values...))
}
func (l *UtilLogger) Infof(format string, values ...interface{}) {
	fmt.Printf("[Infof] %s", fmt.Sprintf(format, values...))
}
func (l *UtilLogger) Warnf(format string, values ...interface{}) {
	fmt.Printf("[Warnf] %s", fmt.Sprintf(format, values...))
}
func (l *UtilLogger) Errorf(format string, values ...interface{}) {
	fmt.Printf("[Errorf] %s", fmt.Sprintf(format, values...))
}
func (l *UtilLogger) Fatalf(format string, values ...interface{}) {
	fmt.Printf("[Fatalf] %s", fmt.Sprintf(format, values...))
}
func (l *UtilLogger) Panicf(format string, values ...interface{}) {
	fmt.Printf("[Panicf] %s", fmt.Sprintf(format, values...))
}
func (l *UtilLogger) Trace(values ...interface{}) { fmt.Printf("[Trace] %s", fmt.Sprint(values...)) }
func (l *UtilLogger) Debug(values ...interface{}) { fmt.Printf("[Debug] %s", fmt.Sprint(values...)) }
func (l *UtilLogger) Info(values ...interface{})  { fmt.Printf("[Info] %s", fmt.Sprint(values...)) }
func (l *UtilLogger) Warn(values ...interface{})  { fmt.Printf("[Warn] %s", fmt.Sprint(values...)) }
func (l *UtilLogger) Error(values ...interface{}) { fmt.Printf("[Error] %s", fmt.Sprint(values...)) }
func (l *UtilLogger) Fatal(values ...interface{}) { fmt.Printf("[Fatal] %s", fmt.Sprint(values...)) }
func (l *UtilLogger) Panic(values ...interface{}) { fmt.Printf("[Panic] %s", fmt.Sprint(values...)) }

type Logger interface {
	Tracef(format string, values ...interface{})
	Debugf(format string, values ...interface{})
	Infof(format string, values ...interface{})
	Warnf(format string, values ...interface{})
	Errorf(format string, values ...interface{})
	Fatalf(format string, values ...interface{})
	Panicf(format string, values ...interface{})
	Trace(values ...interface{})
	Debug(values ...interface{})
	Info(values ...interface{})
	Warn(values ...interface{})
	Error(values ...interface{})
	Fatal(values ...interface{})
	Panic(values ...interface{})
}
