package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"butterfly.orx.me/core/log"
	"github.com/gin-gonic/gin"
	"github.com/orvice/simpleproxy/internal/conf"
)

func Router(m *gin.Engine) {
	// 注册代理处理函数，使用通配符路由捕获所有请求
	m.NoRoute(proxy)
}

var (
	reserveProxy = make(map[string]*httputil.ReverseProxy)
)

func Init() error {
	for _, proxyConf := range conf.Conf.Proxy {
		upstreamURL, err := url.Parse(proxyConf.Upstream)
		if err != nil {
			return err
		}
		proxy := httputil.NewSingleHostReverseProxy(upstreamURL)
		proxy.Director = func(req *http.Request) {
			logger := log.FromContext(req.Context())

			req.Host = proxyConf.Host
			req.URL.Host = upstreamURL.Host
			req.URL.Scheme = upstreamURL.Scheme
			req.URL.Path = upstreamURL.Path + req.URL.Path

			// remove forward header
			req.Header.Del("X-Forwarded-For")
			req.Header.Del("X-Forwarded-Host")
			req.Header.Del("X-Forwarded-Proto")

			logger.Info("proxy request",
				"new_url", req.URL.String(),
				"host", proxyConf.Host, "upstream", proxyConf.Upstream)
		}
		reserveProxy[proxyConf.Host] = proxy
	}
	return nil
}

func proxy(c *gin.Context) {
	logger := log.FromContext(c.Request.Context())
	host := c.Request.Host
	logger.Info("new proxy request",
		"path", c.Request.URL.Path,
		"host", host)

	proxy, ok := reserveProxy[host]
	if !ok {
		logger.Error("proxy not found", "host", host)
		c.JSON(http.StatusNotFound, gin.H{"error": "proxy not found"})
		return
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
