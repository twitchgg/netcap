package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"anyun.bitbucket.com/netcap/internal/api"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
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
	// bindAddr 绑定地址
	bindAddr string
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
		server, err := api.NewServer(&api.Config{
			BindAddr: envs.bindAddr,
		})
		if err != nil {
			logrus.WithField("prefix", "root_cmd").WithError(err).
				Fatal("应用程序创建错误: %s", err.Error())
		}
		errChan = server.Start()
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

	rootCmd.Flags().StringVar(&envs.dev, "dev", "", "监听网卡设备名称")
	rootCmd.Flags().StringVar(&envs.loggerLevel, "logger-level", "DEBUG", "日志级别")
	rootCmd.Flags().StringVar(&envs.ngrepPath, "ngrep-path", "", "ngrep执行程序路径")
	rootCmd.Flags().StringVar(&envs.bindAddr, "bind-addr", "", "API服务绑定地址")
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
