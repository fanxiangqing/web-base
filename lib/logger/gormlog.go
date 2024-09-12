package logger

import (
	"context"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
	"time"
)

type GormLogger struct {
	config logger.Config
}

func NewGormLogger(config logger.Config) *GormLogger {
	return &GormLogger{
		config: config,
	}
}

func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.config.LogLevel = level
	return &newLogger
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= logger.Info {
		logrus.WithContext(ctx).Infof(msg, data)
	}
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= logger.Warn {
		logrus.WithContext(ctx).Warnf(msg, data...)
	}
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.config.LogLevel >= logger.Error {
		logrus.WithContext(ctx).Errorf(msg, data...)
	}
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.config.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	fields := logrus.Fields{
		"elapsed": elapsed,
		"rows":    rows,
		"sql":     sql,
	}

	switch {
	case err != nil && l.config.LogLevel >= logger.Error:
		fields["error"] = err
		logrus.WithContext(ctx).WithFields(fields).Error("GORM error")
	case elapsed > l.config.SlowThreshold && l.config.SlowThreshold != 0 && l.config.LogLevel >= logger.Warn:
		logrus.WithContext(ctx).WithFields(fields).Warn("GORM slow query")
	case l.config.LogLevel >= logger.Info:
		logrus.WithContext(ctx).WithFields(fields).Info("GORM query")
	}
}
