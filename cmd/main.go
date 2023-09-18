package main

import (
	"l0/internal/config"
	"l0/internal/domain/service"
	"l0/internal/storage"
	"l0/internal/storage/cache"
	"l0/internal/transport/broker"
	"l0/internal/transport/http_server/router"
	"log/slog"
	"os"
	"os/signal"
)

func main() {
	const fn = "main.main"
	slog.Info("start app")

	cnf, err := config.MustLoadConfig()
	if err != nil {
		slog.Error(fn, slog.String("failed to init config error: ", err.Error()))
		os.Exit(1)
	}
	slog.Info("config is loaded")
	s, err := storage.NewStorage(cnf.StorageConnectString)
	if err != nil {
		slog.Error(fn, slog.String("failed to init storage error: ", err.Error()))
		os.Exit(1)
	}
	slog.Info("storage is loaded")
	c, err := cache.NewMemoryCash(s)
	if err != nil {
		slog.Error(fn, slog.String("failed to init cache error: ", err.Error()))
		os.Exit(1)
	}
	slog.Info("cache is loaded")
	b := broker.NewStan()
	err = b.Connect(cnf.Broker.ClusterID, cnf.Broker.ClientID, cnf.Broker.URL)
	if err != nil {
		slog.Error(fn, slog.String("failed to init broker error: ", err.Error()))
		os.Exit(1)
	}
	slog.Info("broker is loaded")
	orderService := service.NewOrderProcessing(s, c, b)
	err = orderService.Save()
	if err != nil {
		slog.Error("failed to init service error: ", err)
		os.Exit(1)
	}
	slog.Info("service is loaded")
	echoRouter := router.NewEchoRouter(orderService)
	slog.Error(fn, slog.String("failed to start server error: ", echoRouter.Start(cnf.Host, cnf.Port).Error()))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)
	select {
	case <-stop:
		slog.Info("app is stopped")
		err = b.Close()
		if err != nil {
			slog.Error(fn, slog.String("failed to close broker error: ", err.Error()))
		}
		err = echoRouter.Stop()
		if err != nil {
			slog.Error(fn, slog.String("failed to close server error: ", err.Error()))
		}
		err = s.Close()
		if err != nil {
			slog.Error(fn, slog.String("failed to close storage error: ", err.Error()))
		}
	}
}
