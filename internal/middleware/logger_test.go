package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogger(t *testing.T) {
	// 设置 gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))
	slog.SetDefault(logger)

	// 创建路由
	r := gin.New()
	r.Use(Logger())

	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	r.GET("/error", func(c *gin.Context) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
	})

	// 测试正常请求
	t.Run("Normal Request", func(t *testing.T) {
		buf.Reset()
		req, _ := http.NewRequest("GET", "/test?foo=bar", nil)
		req.Header.Set("User-Agent", "test-agent")
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// 验证日志输出
		logOutput := buf.String()
		if !strings.Contains(logOutput, "Request completed") {
			t.Error("Log should contain 'Request completed'")
		}
		if !strings.Contains(logOutput, "foo=bar") {
			t.Error("Log should contain query parameters")
		}
	})

	// 测试错误请求
	t.Run("Error Request", func(t *testing.T) {
		buf.Reset()
		req, _ := http.NewRequest("GET", "/error", nil)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, w.Code)
		}

		// 验证日志输出
		logOutput := buf.String()
		if !strings.Contains(logOutput, "Server error") {
			t.Error("Log should contain 'Server error' for 500 status")
		}
	})
}

func TestLoggerWithSlog(t *testing.T) {
	// 设置 gin 为测试模式
	gin.SetMode(gin.TestMode)

	// 创建一个 buffer 来捕获日志输出
	var buf bytes.Buffer
	logger := slog.New(slog.NewJSONHandler(&buf, nil))

	// 创建路由
	r := gin.New()
	r.Use(LoggerWithSlog(logger))

	// 添加测试路由
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "test"})
	})

	// 测试请求
	req, _ := http.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}

	// 验证日志输出是否为有效的 JSON
	var logEntry map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &logEntry); err != nil {
		t.Errorf("Failed to parse log output as JSON: %v", err)
	}

	// 检查必要的字段
	if logEntry["msg"] != "Request completed" {
		t.Error("Log message should be 'Request completed'")
	}
	if logEntry["method"] != "GET" {
		t.Error("Log should contain correct method")
	}
	if logEntry["path"] != "/test" {
		t.Error("Log should contain correct path")
	}
	if logEntry["status"] != float64(200) {
		t.Error("Log should contain correct status code")
	}
}
