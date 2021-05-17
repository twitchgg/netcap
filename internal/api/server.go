package api

import (
	"fmt"
	"net/http"
	"time"

	"anyun.bitbucket.com/netcap/pkg/ngrep"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

type Server struct {
	conf             *Config
	server           *echo.Echo
	httpServer       *http.Server
	ng               *ngrep.Application
	local            *time.Location
	lastDumpFileName string
}

// NewServer 创建接口服务器
func NewServer(conf *Config) (*Server, error) {
	if conf == nil {
		return nil, fmt.Errorf("HTTP API 服务配置未定义")
	}
	if err := conf.Check(); err != nil {
		return nil, fmt.Errorf("HTTP API 服务配置检查错误: %s", err.Error())
	}
	e := echo.New()
	e.HideBanner = true
	e.HTTPErrorHandler = SimpleAPIErrorHandler
	bindAddr := conf.BindAddr
	httpServer := &http.Server{
		ReadTimeout:  DEFAULT_HTTP_READ_TIMEOUT,
		WriteTimeout: DEFAULT_HTTP_WRITE_TIMEOUT,
		Addr:         bindAddr,
	}
	l, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return nil, fmt.Errorf("获取时区信息错误: %s", err.Error())
	}
	return &Server{
		conf:       conf,
		httpServer: httpServer,
		server:     e,
		local:      l,
	}, nil
}

// Start 启动API服务
func (s *Server) Start() chan error {
	errChan := make(chan error)
	if err := s.router(s.server); err != nil {
		errChan <- err
		return errChan
	}
	go func() {
		logrus.WithField("prefix", "http_api").
			Infof("启动HTTP API服务,服务监听地址 [%s]", s.conf.BindAddr)
		errChan <- s.server.StartServer(s.httpServer)
	}()
	return errChan
}
