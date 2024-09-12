package logger

import (
	"bufio"
	"github.com/fanxiangqing/web-base/internal/lib/utils"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func Init(logFile, logLevel string, maxRemainCnt uint, rotationTime time.Duration) {
	// 设置日志
	setupLoggingJson(logFile, maxRemainCnt, rotationTime)
	l, _ := logrus.ParseLevel(logLevel)
	logrus.SetLevel(l)
}

// setupLoggingJson
func setupLoggingJson(logFile string, maxRemainCnt uint, rotationTime time.Duration) {
	// 日志文件不为空则关闭stdout
	if logFile != "" {
		logDir, err := filepath.Abs(filepath.Dir(logFile))
		if err != nil {
			logrus.Panicf("fail to get dir of LogFile(%s)![%v]\n", logFile, err)
		}

		if ok, _ := utils.IsFileExist(logDir); !ok {
			if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
				logrus.Panicf("mkdir logDir[%s] failed![%v]\n", logDir, err)
			}
		}

		configFileLoggerJson(logFile, maxRemainCnt, rotationTime*time.Hour)
		// configFileLoggerJson(logFile, maxRemainCnt, rotationTime*time.Second)

		src, err := os.OpenFile(os.DevNull, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		if err != nil {
			logrus.Panicf("err", err)
		}
		writer := bufio.NewWriter(src)
		logrus.SetOutput(writer)
		logrus.SetOutput(os.Stdout)
	}

	formatter := logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableHTMLEscape: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	}
	logrus.SetFormatter(&formatter)

}

// config logrus log to local filesystem, with file rotation
func configFileLoggerJson(logFileName string, maxRemainCnt uint, rotationTime time.Duration) {
	writer, err := rotatelogs.New(
		logFileName+"-%Y%m%d%H%M%S.log",
		rotatelogs.WithLinkName(logFileName+".log"), // 生成软链，指向最新日志文件
		rotatelogs.WithRotationCount(maxRemainCnt),  // 设置文件清理前最多保存的个数
		rotatelogs.WithRotationTime(rotationTime),   // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.JSONFormatter{
		TimestampFormat:   "2006-01-02 15:04:05",
		DisableHTMLEscape: true,
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyMsg: "message",
		},
	})

	logrus.AddHook(lfHook)
}
