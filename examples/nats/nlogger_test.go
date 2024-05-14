package natslogger

import (
	"fmt"
	"testing"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

func natsServer(t testing.TB) *server.Server {
	t.Helper()

	s, err := server.NewServer(
		&server.Options{
			Port: -1,
		},
	)
	if err != nil {
		server.PrintAndDie(err.Error())
	}

	s.ConfigureLogger()

	return s
}

func TestLogger(t *testing.T) {
	s := natsServer(t)
	defer s.Shutdown()
	go s.Start()

	// wait for server to start
	time.Sleep(time.Second)

	done := make(chan struct{}, 1)

	nc, err := nats.Connect(s.ClientURL(), nats.ClosedHandler(func(nc *nats.Conn) {
		done <- struct{}{}
	}))
	if err != nil {
		t.Fatalf("Failed to connect to server: %v", err)
	}

	stdout := []byte{}
	stderr := []byte{}

	_, err = nc.Subscribe("stdout", func(m *nats.Msg) {
		stdout = append(stdout, m.Data...)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to stdout: %v", err)
	}
	_, err = nc.Subscribe("stderr", func(m *nats.Msg) {
		stderr = append(stderr, m.Data...)
	})
	if err != nil {
		t.Fatalf("Failed to subscribe to stderr: %v", err)
	}

	now := time.Now().Format("15:04")

	// create logger
	logger := NewLogger(nc)
	logger.Info("info")
	logger.Error("error")

	err = nc.Drain()
	if err != nil {
		t.Fatalf("Failed to drain: %v", err)
	}
	<-done

	if string(stdout) != fmt.Sprintf("[INFO] %s - info\n", now) {
		t.Errorf("Expected stdout to be 'info\n', got '%s'", stdout)
		t.Fail()
	}
	if string(stderr) != fmt.Sprintf("[ERROR] %s - error\n", now) {
		t.Errorf("Expected stderr to be 'error\n', got '%s'", stderr)
		t.Fail()
	}
}
