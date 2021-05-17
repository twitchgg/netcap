package api

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"anyun.bitbucket.com/netcap/pkg/ngrep"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) router(e *echo.Echo) error {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	v1 := e.Group("/cap/")
	v1.GET("", s.v1DescRouter)
	v1.POST("ip", s.v1IPsCap)
	return nil
}

func (s *Server) v1DescRouter(c echo.Context) error {
	return nil
}

func (s *Server) v1IPsCap(c echo.Context) (err error) {
	ipBegin := strings.ToUpper(c.FormValue("ip_begin"))
	ipEnd := strings.ToUpper(c.FormValue("ip_end"))
	stopTime := strings.ToUpper(c.FormValue("stop_time"))
	dumpSize := strings.ToUpper(c.FormValue("dump_size"))
	if ipBegin == "" {
		return fmt.Errorf("缺失起始IP参数参数 [ip_begin]")
	}
	if _, err := net.ResolveIPAddr("ip4", ipBegin); err != nil {
		return fmt.Errorf("起始IP解析错误: %s", err.Error())
	}
	if ipEnd == "" {
		return fmt.Errorf("缺失结束IP参数 [ip_end]")
	}
	if _, err := net.ResolveIPAddr("ip4", ipEnd); err != nil {
		return fmt.Errorf("结束IP解析错误: %s", err.Error())
	}
	if stopTime == "" {
		return fmt.Errorf("缺失任务结束时间参数 [stop_time]")
	}
	if dumpSize == "" {
		return fmt.Errorf("缺失dump文件大小参数 [dump_size]")
	}
	logrus.WithField("prefix", "api").
		Debugf("ip_begin [%s] ip_end [%s] stop_time [%s] dump_size [%s]",
			ipBegin, ipEnd, stopTime, dumpSize)
	if s.ng == nil {
		if s.ng, err = ngrep.NewApplication(&ngrep.AppConfig{}); err != nil {
			return fmt.Errorf("创建ngrep应用程序错误: %s", err.Error())
		}
	}
	if s.ng.IsStart {
		return fmt.Errorf("ngrep应用程序在运行中")
	}
	go func() {
		ips := GetIpsFromRange(ipBegin, ipEnd)
		params := GenNgrepParams("", "", ips, "./dump.pcap")
		if err := <-s.ng.Start(params); err != nil {
			logrus.WithField("prefix", "api").WithError(err).Error("ngrep应用程序启动失败")
		}
	}()
	return c.String(http.StatusOK, "ok")
}

// GetIpsFromRange 通过起始结束IP获取所有IP地址
func GetIpsFromRange(begin, end string) []string {
	ips := make([]string, 0)
	p11 := strings.Split(begin, ".")
	p21 := strings.Split(end, ".")
	_p1, _ := strconv.Atoi(p11[3])
	_p2, _ := strconv.Atoi(p21[3])
	for i := _p1; i <= _p2; i++ {
		fp := fmt.Sprintf("%s.%s.%s.%d", p11[0], p11[1], p11[2], i)
		ips = append(ips, fp)
	}
	return ips
}
