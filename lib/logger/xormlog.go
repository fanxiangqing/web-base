package logger

import (
	"github.com/sirupsen/logrus"
	"xorm.io/xorm/log"
)

var _ log.Logger = new(XORMLoggerStruct)
var XORMLogger *XORMLoggerStruct = new(XORMLoggerStruct)

type XORMLoggerStruct struct{}

func (l *XORMLoggerStruct) Debug(v ...interface{}) {
	logrus.Debug(v...)
}
func (l *XORMLoggerStruct) Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}
func (l *XORMLoggerStruct) Error(v ...interface{}) {
	logrus.Error(v...)
}
func (l *XORMLoggerStruct) Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}
func (l *XORMLoggerStruct) Info(v ...interface{}) {
	logrus.Info(v...)
}
func (l *XORMLoggerStruct) Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}
func (l *XORMLoggerStruct) Warn(v ...interface{}) {
	logrus.Warning(v...)
}
func (l *XORMLoggerStruct) Warnf(format string, v ...interface{}) {
	logrus.Warningf(format, v...)
}

func (l *XORMLoggerStruct) Level() log.LogLevel {
	return log.LOG_DEBUG
}

func (l *XORMLoggerStruct) SetLevel(lv log.LogLevel) {
}

func (l *XORMLoggerStruct) ShowSQL(show ...bool) {
}
func (l *XORMLoggerStruct) IsShowSQL() bool {
	return true
}
