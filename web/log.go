package web

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const (
	logKeyRequestID = "requestId"
	logKeyModule    = "module"
)

// LogContext 包装 gin 的上下文，用于打印日志，并给日志添加唯一的请求 ID。
type LogContext struct {
	Ctx    *gin.Context
	Module string
	Logger *logrus.Logger
}

func (c *LogContext) entry() (entry *logrus.Entry) {
	if c.Logger != nil {
		entry = logrus.NewEntry(c.Logger)
	} else {
		entry = logrus.NewEntry(logrus.StandardLogger())
	}
	if c.Ctx != nil {
		if requestID := c.Ctx.GetString(contextKeyRequestID); requestID != "" {
			entry = entry.WithField(logKeyRequestID, c.Ctx.GetString(contextKeyRequestID))
		}
	}
	if c.Module != "" {
		entry = entry.WithField(logKeyModule, c.Module)
	}
	return entry
}

// Info 打印 Info 级别日志打印。
func (c *LogContext) Info(args ...interface{}) {
	c.entry().Info(args...)
}

// Infof 格式化打印 Info 级别日志。
func (c *LogContext) Infof(format string, args ...interface{}) {
	c.entry().Infof(format, args...)
}

// Warn 打印 Warn 级别日志打印。
func (c *LogContext) Warn(args ...interface{}) {
	c.entry().Warn(args...)
}

// Warnf 格式化打印 Warn 级别日志。
func (c *LogContext) Warnf(format string, args ...interface{}) {
	c.entry().Warnf(format, args...)
}

// Error 打印 Error 级别日志打印。
func (c *LogContext) Error(args ...interface{}) {
	c.entry().Error(args...)
}

// Errorf 格式化打印 Error 级别日志。
func (c *LogContext) Errorf(format string, args ...interface{}) {
	c.entry().Errorf(format, args...)
}
