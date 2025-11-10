// Package logger
package logger

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Log = InitLogger("dev")

func InitLogger(env string) *logrus.Logger {
	log := logrus.New()
	log.SetOutput(io.MultiWriter(
		os.Stdout,
		&lumberjack.Logger{
			Filename:   "logs/go-trading-bot.log",
			MaxSize:    10, // MB
			MaxBackups: 5,
			MaxAge:     7, // days
			Compress:   true,
		}))

	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		ForceQuote:      true,
		DisableColors:   false,
		DisableQuote:    false,
	})
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)

	return log
}

func Info(args ...interface{})          { Log.Info(args...) }
func Warn(args ...interface{})          { Log.Warn(args...) }
func Error(args ...interface{})         { Log.Error(args...) }
func Debug(args ...interface{})         { Log.Debug(args...) }
func Infof(f string, a ...interface{})  { Log.Infof(f, a...) }
func Warnf(f string, a ...interface{})  { Log.Warnf(f, a...) }
func Errorf(f string, a ...interface{}) { Log.Errorf(f, a...) }
func Debugf(f string, a ...interface{}) { Log.Debugf(f, a...) }
