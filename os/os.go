package wolf

import (
	"github.com/moesn/wolf/db"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

// HandleSignal 优雅的退出GO守护进程
func HandleSignal(server *http.Server) {
	c := make(chan os.Signal)

	// 用户发送INTR字符(Ctrl+C)触发
	// 用户发送QUIT字符(Ctrl+/)触发
	// 结束程序(可以被捕获、阻塞或忽略)
	signal.Notify(c, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	go func() {
		s := <-c
		logrus.Infof("收到退出信号 [%s], 正在退出", s)

		if err := server.Close(); err != nil {
			logrus.Errorf("服务关闭失败: " + err.Error())
		}

		// 关闭数据库连接
		db.Close()

		logrus.Infof("已退出")
		os.Exit(0)
	}()
}
