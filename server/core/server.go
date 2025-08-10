package core

import (
	"context"

	"net/http"
	"os"
	"os/signal"
	"server/global"
	"server/initialize"
	"syscall"
	"time"

	"go.uber.org/zap"
)

func StartServer() {
	router := initialize.InitRouter()

	srv := &http.Server{
		Addr:           global.Config.System.Addr(),
		Handler:        router,
		ReadTimeout:    global.Config.System.ReadTimeout,
		WriteTimeout:   global.Config.System.WriteTimeout,
		IdleTimeout:    global.Config.System.IdleTimeout,
		MaxHeaderBytes: global.Config.System.MaxHeaderBytes,
	}
	global.ZapLog.Info("Server is running on %s", zap.String("address", global.Config.System.Addr()))

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.ZapLog.Error("listen:%s\n", zap.Error(err))
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	global.ZapLog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.ZapLog.Error("Server forced to shutdown:%v\n", zap.Error(err))
	}

	global.ZapLog.Info("Server exiting")
}
