package web

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

const (
	contextKeyRequestID = "gin-context-request-id"
)

// RequestIDMiddleware 用于在请求到来时生成请求的 uuid 并存入上下文中，便于请求追踪。
func RequestIDMiddleware(c *gin.Context) {
	// 设置 requestId。
	c.Set(contextKeyRequestID, uuid.New().String())
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// LoggerMiddleware 使用 logrus 的 std 实例打印请求和响应相关信息的日志，便于 debug。
func LoggerMiddleware(printBody bool) func(*gin.Context) {
	return LoggerMiddlewareWithLogger(printBody, logrus.StandardLogger())
}

// LoggerMiddlewareWithLogger 使用指定的 logrus 的 Logger 实例打印请求和响应相关信息的日志，便于 debug。
func LoggerMiddlewareWithLogger(printBody bool, logger *logrus.Logger) func(*gin.Context) {
	return func(c *gin.Context) {
		var (
			entry     *logrus.Entry
			reqURI    string
			startTime time.Time
		)

		// 打印请求前缀日志。
		reqURI = c.Request.RequestURI
		if reqURI == "" {
			reqURI = c.Request.URL.RequestURI()
		}

		// 健康检查接口不打印日志。
		if reqURI != "/healthz" {
			entry = logrus.NewEntry(logger)
			if requestID := c.GetString(contextKeyRequestID); requestID != "" {
				entry = entry.WithField(logKeyRequestID, c.GetString(contextKeyRequestID))
			}

			if printBody {
				var (
					body []byte
					buf  bytes.Buffer
				)
				if c.Request.Body != nil {
					tee := io.TeeReader(c.Request.Body, &buf)
					body, _ = ioutil.ReadAll(tee)
					c.Request.Body = ioutil.NopCloser(&buf)
				}

				entry.Infof(
					"|RequestLog| %s %s HTTP/%d.%d, body: %s, client: %s",
					c.Request.Method, reqURI, c.Request.ProtoMajor, c.Request.ProtoMinor,
					string(body), c.ClientIP(),
				)
			} else {
				entry.Infof(
					"|RequestLog| %s %s HTTP/%d.%d, client: %s",
					c.Request.Method, reqURI, c.Request.ProtoMajor, c.Request.ProtoMinor,
					c.ClientIP(),
				)
			}

			startTime = time.Now()
		}

		var w *responseBodyWriter
		if printBody {
			// 替换 gin 的 Writer，便于打印 Response Body。
			w = &responseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
			c.Writer = w
		}

		c.Next()

		if reqURI != "/healthz" && entry != nil {
			if printBody {
				entry.Infof(
					"|ResponseLog| HTTP/%d.%d %d %s, path: %s, costs: %s, body: %s, client: %s",
					c.Request.ProtoMajor, c.Request.ProtoMinor,
					c.Writer.Status(), http.StatusText(c.Writer.Status()),
					reqURI, time.Since(startTime).String(), w.body.String(), c.ClientIP(),
				)
			} else {
				entry.Infof(
					"|ResponseLog| HTTP/%d.%d %d %s, path: %s, costs: %s, client: %s",
					c.Request.ProtoMajor, c.Request.ProtoMinor,
					c.Writer.Status(), http.StatusText(c.Writer.Status()),
					reqURI, time.Since(startTime).String(), c.ClientIP(),
				)
			}
		}
	}
}
