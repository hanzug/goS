package shutdown

import (
	"context"
	"go.uber.org/zap"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func GracefullyShutdown(server *http.Server) {
	// 创建系统信号接收器接收关闭信号
	done := make(chan os.Signal, 1)
	/**
	os.Interrupt           -> ctrl+c 的信号
	syscall.SIGINT|SIGTERM -> kill 进程时传递给进程的信号
	*/
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-done:
		{
			zap.S().Infof("stopping service, because received signal:", sig)
			if err := server.Shutdown(context.Background()); err != nil {
				zap.S().Infof("closing http service gracefully failed: :%v", err)
			}
			zap.S().Infof("service has stopped")
			os.Exit(0)
		}
	default:
		os.Exit(0)
	}
}
