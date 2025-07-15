package handler

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"butterfly.orx.me/core/log"
	"github.com/gin-gonic/gin"
	"github.com/orvice/simpleproxy/internal/conf"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
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

			// set host header to upstream host
			req.Host = upstreamURL.Host
			req.URL.Host = upstreamURL.Host
			req.URL.Scheme = upstreamURL.Scheme
			req.URL.Path = upstreamURL.Path + req.URL.Path

			logger.Info("proxy request",
				"new_url", req.URL.String(),
				"host", proxyConf.Host, "upstream", proxyConf.Upstream)
		}
		proxy.Transport = otelhttp.NewTransport(http.DefaultTransport)
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
