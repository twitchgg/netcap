package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

const (
	DEFAULT_HTTP_READ_TIMEOUT  = 5 * time.Second
	DEFAULT_HTTP_WRITE_TIMEOUT = 5 * time.Second
)

// SimpleAPIErrorHandler API错误处理器
func SimpleAPIErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
	}
	if err := c.String(code, err.Error()); err != nil {
		c.Logger().Error(err)
	}
}

// Config HTTP API 服务配置
type Config struct {
	BindAddr string
}

// Check 检查API服务配置
func (c *Config) Check() error {
	return nil
}

// GenNgrepParams 创建ngrep参数
func GenNgrepParams(dev string, keyword string, hostIPS []string, dumpPath string) []string {
	params := make([]string, 0)
	if dev != "" {
		params = append(params, "-d")
		params = append(params, dev)
	}

	params = append(params, "-i")
	params = append(params, keyword)
	if len(hostIPS) > 0 {
		for _, ip := range hostIPS {
			params = append(params, "src")
			params = append(params, "host")
			params = append(params, ip)
			params = append(params, "or")
		}
		params = params[:len(params)-1]
	}
	params = append(params, "-x")
	params = append(params, "-O")
	params = append(params, dumpPath)
	return params
}
