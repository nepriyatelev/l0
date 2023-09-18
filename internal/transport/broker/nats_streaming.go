package broker

import (
	"github.com/nats-io/stan.go"
	"log/slog"
)

type Stan struct {
	sc  stan.Conn
	sub stan.Subscription
}

func NewStan() *Stan {
	return &Stan{}
}

func (s *Stan) Connect(clusterID, clientID, URL string) error {
	const fn = "broker.Stan.Connect"
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(URL))
	if err != nil {
		slog.Error(fn, slog.String("failed to connect to broker error: ", err.Error()))
		return err
	}
	s.sc = sc
	slog.Info("connected to broker")
	return nil
}

func (s *Stan) Subscribe(subject string, cb stan.MsgHandler) error {
	const fn = "broker.Stan.Subscribe"
	sub, err := s.sc.Subscribe(subject, cb, stan.SetManualAckMode())
	if err != nil {
		slog.Error(fn, slog.String("failed to subscribe to broker error: ", err.Error()))
		return err
	}
	s.sub = sub
	slog.Info("subscribed to broker")
	return nil
}

func (s *Stan) Close() error {
	const fn = "broker.Stan.Close"
	err := s.sub.Unsubscribe()
	if err != nil {
		slog.Error(fn, slog.String("failed to unsubscribe from broker error: ", err.Error()))
		return err
	}
	slog.Info("unsubscribed from broker")
	err = s.sc.Close()
	if err != nil {
		slog.Error(fn, slog.String("failed to close connection to broker error: ", err.Error()))
		return err
	}
	slog.Info("connection to broker is closed")
	return nil
}
