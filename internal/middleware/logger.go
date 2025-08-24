package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger 返回一个使用 slog 记录请求的中间件
func Logger() gin.HandlerFunc {
	logger := slog.Default()

	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算请求耗时
		latency := time.Since(start)

		// 获取请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		// 构建日志字段
		logFields := []any{
			"client_ip", clientIP,
			"method", method,
			"path", path,
			"status", statusCode,
			"latency_ms", latency.Milliseconds(),
			"user_agent", userAgent,
		}

		// 如果有查询参数，也记录下来
		if raw != "" {
			logFields = append(logFields, "query", raw)
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			logFields = append(logFields, "errors", c.Errors.String())
		}

		// 根据状态码决定日志级别
		switch {
		case statusCode >= 500:
			logger.Error("Server error", logFields...)
		case statusCode >= 400:
			logger.Warn("Client error", logFields...)
		case statusCode >= 300:
			logger.Info("Redirect", logFields...)
		default:
			logger.Info("Request completed", logFields...)
		}
	}
}

// LoggerWithSlog 使用提供的 slog.Logger 记录请求
func LoggerWithSlog(logger *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 计算请求耗时
		latency := time.Since(start)

		// 获取请求信息
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		userAgent := c.Request.UserAgent()

		// 构建日志字段
		logFields := []any{
			"client_ip", clientIP,
			"method", method,
			"path", path,
			"status", statusCode,
			"latency_ms", latency.Milliseconds(),
			"user_agent", userAgent,
		}

		// 如果有查询参数，也记录下来
		if raw != "" {
			logFields = append(logFields, "query", raw)
		}

		// 如果有错误，记录错误信息
		if len(c.Errors) > 0 {
			logFields = append(logFields, "errors", c.Errors.String())
		}

		// 根据状态码决定日志级别
		switch {
		case statusCode >= 500:
			logger.Error("Server error", logFields...)
		case statusCode >= 400:
			logger.Warn("Client error", logFields...)
		case statusCode >= 300:
			logger.Info("Redirect", logFields...)
		default:
			logger.Info("Request completed", logFields...)
		}
	}
}
