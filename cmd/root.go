package cmd

import (
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"anyun.bitbucket.com/netcap/pkg/ngrep"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

// envs ngrep执行环境变量
var envs struct {
	// ngrep执行路径
	ngrepPath string
	// ngrep监听设备
	dev string
	// 日志级别
	loggerLevel string
	// 过滤关键字
	keyword string
	// 监听端口
	ports string
	// pcap dump文件路径
	dumpPath string
	// 监听IP地址
	hostIps string
	// 监听网络
	hostNets string
}
var errChan chan error
var rootCmd = &cobra.Command{
	Use:   "netcap",
	Short: "网络抓包过滤节点",
	PreRun: func(cmd *cobra.Command, args []string) {
		logrus.SetOutput(os.Stdout)
		formatter := new(prefixed.TextFormatter)
		logrus.SetFormatter(formatter)
		lvl, err := logrus.ParseLevel(envs.loggerLevel)
		if err != nil {
			logrus.WithField("prefix", "root_cmd").WithError(err).
				Fatal("日志级别解析错误")
		}
		logrus.SetLevel(lvl)
	},
	Run: func(cmd *cobra.Command, args []string) {
		ports1 := strings.Split(envs.ports, ",")
		ports2 := make([]int, len(ports1))
		for i, v := range ports1 {
			aport, err := strconv.Atoi(v)
			if err != nil {
				logrus.WithField("prefix", "root_cmd").WithError(err).
					Fatal("监听端口解析错误")
			}
			ports2[i] = aport
		}
		hostIps := []string{}
		if envs.hostIps != "" {
			hostIps = strings.Split(envs.hostIps, ",")
		}
		hostNets := []string{}
		if envs.hostNets != "" {
			hostNets = strings.Split(envs.hostNets, ",")
		}
		app, err := ngrep.NewApplication(&ngrep.AppConfig{
			NgrepPath: envs.ngrepPath,
			Dev:       envs.dev,
			Keyword:   envs.keyword,
			DumpPath:  envs.dumpPath,
			HostIPS:   hostIps,
			HostNets:  hostNets,
			Ports:     ports2,
		})
		if err != nil {
			logrus.WithField("prefix", "root_cmd").WithError(err).
				Fatal("应用程序创建错误: %s", err.Error())
		}
		errChan = app.Start()
	},
	PostRun: func(cmd *cobra.Command, args []string) {
		go func() {
			if err := <-errChan; err != nil {
				logrus.WithField("prefix", "root_cmd").WithError(err).
					Fatalf("应用程序运行错误")
			}
		}()
		RunWithSysSignal(nil)
	},
}

func init() {
	cobra.OnInitialize(func() {})
	viper.AutomaticEnv()
	viper.SetEnvPrefix("NP")

	rootCmd.Flags().StringVar(&envs.dev, "dev", "", "监听网卡设备名称")
	rootCmd.Flags().StringVar(&envs.keyword, "keyword", "", "过滤关键字正则表达式")
	rootCmd.Flags().StringVar(&envs.loggerLevel, "logger-level", "DEBUG", "日志级别")
	rootCmd.Flags().StringVar(&envs.ngrepPath, "ngrep-path", "", "ngrep执行程序路径")
	rootCmd.Flags().StringVar(&envs.ports, "ports", "", "抓包监听端口")
	rootCmd.Flags().StringVar(&envs.dumpPath, "dump-path", "./dump.pcap", "pcap dump文件路径")
	rootCmd.Flags().StringVar(&envs.hostIps, "host-ip", "", "监听IP地址")
	rootCmd.Flags().StringVar(&envs.hostNets, "host-net", "", "监听IP网络")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.WithError(err).Error("系统运行错误")
		os.Exit(1)
	}
}

// RunWithSysSignal 等待系统关闭信号
func RunWithSysSignal(clearFunc func(os.Signal)) {
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		logrus.WithField("prefix", "root_cmd").
			Infof("系统信号: %s", sig)
		if clearFunc != nil {
			clearFunc(sig)
		}
		done <- true
	}()
	<-done
}
