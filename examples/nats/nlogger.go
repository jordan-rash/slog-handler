package natslogger

import (
	"log/slog"
	"os"

	"github.com/jordan-rash/slog-handler"
	"github.com/nats-io/nats.go"
)

func NewLogger(nc *nats.Conn) *slog.Logger {
	nstdout := NewNatsLogger(nc, "stdout")
	nstderr := NewNatsLogger(nc, "stderr")

	logger := slog.New(
		handler.NewHandler(
			handler.WithTimeFormat("15:04"),
			handler.WithStdOut(os.Stdout, nstdout),
			handler.WithStdErr(os.Stderr, nstderr),
		),
	)

	return logger
}

type NatsLogger struct {
	nc       *nats.Conn
	OutTopic string
}

func NewNatsLogger(nc *nats.Conn, topic string) *NatsLogger {
	return &NatsLogger{
		OutTopic: topic,
		nc:       nc,
	}
}

func (nl *NatsLogger) Write(p []byte) (int, error) {
	err := nl.nc.Publish(nl.OutTopic, p)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}
