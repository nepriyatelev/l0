package main

import (
	"bufio"
	"fmt"
	"github.com/nats-io/stan.go"
	"l0/internal/config"
	"log/slog"
	"os"
)

const (
	clientID = "client-1"
	subject  = "order"
)

func main() {
	cnf, _ := config.MustLoadConfig()
	slog.Info("config is loaded")
	st, err := stan.Connect(cnf.Broker.ClusterID, clientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		slog.Error("failed to init broker error: ", err)
		return
	}
	var fileName string
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("insert file name")
		fmt.Fscan(reader, &fileName)
		data, readErr := os.ReadFile(fileName)
		if readErr != nil {
			slog.Error("failed to read file error: ", readErr)
			continue
		}
		err = st.Publish(subject, data)
		if err != nil {
			slog.Error("failed to publish message error: ", err)
			continue
		}
		slog.Info("message published")
	}
}
