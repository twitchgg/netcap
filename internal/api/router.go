package api

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"anyun.bitbucket.com/netcap/pkg/ngrep"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

func (s *Server) router(e *echo.Echo) error {
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "OK")
	})
	v1 := e.Group("/cap/")
	v1.GET("stat", s.v1Stat)
	v1.GET("file", s.v1File)
	v1.POST("ip", s.v1IPsCap)
	v1.POST("keyword", s.v1Keyword)
	v1.POST("stop", s.v1Stop)
	return nil
}

func (s *Server) v1File(c echo.Context) error {
	if s.lastDumpFileName == "" {
		return fmt.Errorf("dump文件不存在")
	}
	c.Response().Header().Add("Content-Type", "application/octet-stream")
	c.Response().Header().Add("Content-Disposition", "attachment;filename=stream.pcap")
	return c.File(s.lastDumpFileName)
}

func (s *Server) v1Stat(c echo.Context) error {
	if s.lastDumpFileName == "" {
		return fmt.Errorf("dump文件不存在")
	}
	fi, err := os.Stat(s.lastDumpFileName)
	if err != nil {
		return fmt.Errorf("读取文件信息错误: %s", err.Error())
	}
	return c.String(http.StatusOK, fmt.Sprintf("%d", fi.Size()))
}
func (s *Server) v1Stop(c echo.Context) (err error) {
	if err := s.ng.Stop(); err != nil {
		return fmt.Errorf("ngrep应用程序关闭失败: %s", err.Error())
	}
	return c.String(http.StatusOK, "ok")
}

func (s *Server) v1Keyword(c echo.Context) (err error) {
	stopTime := strings.ToLower(c.FormValue("stop_time"))
	if stopTime == "" {
		return fmt.Errorf("缺失任务结束时间参数 [stop_time]")
	}
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", stopTime, s.local)
	if err != nil {
		return fmt.Errorf("停止时间解析错误: %s", err.Error())
	}
	keyword := c.FormValue("keyword")
	if keyword == "" {
		return fmt.Errorf("缺失关键字参数 [keyword]")
	}
	width := strings.ToLower(c.FormValue("width"))
	if width == "" {
		width = "dynamic"
	}
	if width != "fixed" && width != "dynamic" {
		return fmt.Errorf("不支持的关键字类型: %s", width)
	}
	dumpSize := strings.ToUpper(c.FormValue("dump_size"))
	if dumpSize == "" {
		return fmt.Errorf("缺失dump文件大小参数 [dump_size]")
	}
	logrus.WithField("prefix", "api").
		Debugf("keyword [%s] width [%s] stop_time [%s] dump_size [%s]",
			keyword, width, stopTime, dumpSize)
	if s.ng == nil {
		if s.ng, err = ngrep.NewApplication(&ngrep.AppConfig{}); err != nil {
			return fmt.Errorf("创建ngrep应用程序错误: %s", err.Error())
		}
	}
	if s.ng.IsStart {
		return fmt.Errorf("ngrep应用程序在运行中")
	}

	go func() {
		s.lastDumpFileName = "./data/keyword_" + width + "_dump.pcap"
		os.Remove(s.lastDumpFileName)
		if width == "fixed" {
			keyword = "^" + keyword
		}
		params := GenNgrepParams(s.conf.Dev, keyword, nil, s.lastDumpFileName)
		go func() {
			for {
				time.Sleep(time.Second)
				if time.Since(t1) >= 0 {
					if !s.ng.IsStart {
						return
					}
					if err := s.ng.Stop(); err != nil {
						logrus.WithField("prefix", "api").WithError(err).Error("ngrep应用程序关闭失败")
					} else {
						logrus.WithField("prefix", "api").Info("ngrep应用程序关闭")
					}
					return
				}
			}
		}()
		if err := <-s.ng.Start(params); err != nil {
			logrus.WithField("prefix", "api").WithError(err).Error("ngrep应用程序启动失败")
		}

	}()
	return c.String(http.StatusOK, "ok")
}

func (s *Server) v1IPsCap(c echo.Context) (err error) {
	ipBegin := strings.ToLower(c.FormValue("ip_begin"))
	ipEnd := strings.ToLower(c.FormValue("ip_end"))
	stopTime := strings.ToLower(c.FormValue("stop_time"))
	t1, err := time.ParseInLocation("2006-01-02 15:04:05", stopTime, s.local)
	if err != nil {
		return fmt.Errorf("停止时间解析错误: %s", err.Error())
	}
	dumpSize := strings.ToLower(c.FormValue("dump_size"))
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
		s.lastDumpFileName = "./data/ip_range_dump.pcap"
		os.Remove(s.lastDumpFileName)
		params := GenNgrepParams(s.conf.Dev, "", ips, s.lastDumpFileName)
		go func() {
			for {
				time.Sleep(time.Second)
				if time.Since(t1) >= 0 {
					if !s.ng.IsStart {
						return
					}
					if err := s.ng.Stop(); err != nil {
						logrus.WithField("prefix", "api").WithError(err).Error("ngrep应用程序关闭失败")
					} else {
						logrus.WithField("prefix", "api").Info("ngrep应用程序关闭")
					}
					return
				}
			}
		}()
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
