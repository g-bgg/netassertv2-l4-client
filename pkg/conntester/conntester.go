// Package conntester is a TCP / UDP client for generating traffic
package conntester

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/controlplaneio/netassertv2-l4-client/pkg/config"
	"github.com/controlplaneio/netassertv2-l4-client/pkg/log"
)

// ErrTestFailed is returned when the passed connection attempts are below the success threshold
var ErrTestFailed = errors.New("test failed")

// ConnTester performs connections against a certain destination
type ConnTester struct {
	protocol   string
	targetHost string
	targetPort uint16
	message    string
	timeout    uint
	attempts   uint
	period     uint
	succThrPec uint

	logger *log.Logger

	mutex       sync.Mutex
	succCounter uint
}

// New initializes a new ConnTester instance
func New(config *config.Config, logger *log.Logger) (*ConnTester, error) {
	if config.TargetHost == "" || config.TargetPort == 0 {
		return nil, errors.New("at least target-host and target-port must be provided")
	}
	if config.Attempts == 0 {
		return nil, errors.New("attempts must be greater than 0")
	}

	ct := &ConnTester{
		protocol:   config.Protocol,
		targetHost: config.TargetHost,
		targetPort: config.TargetPort,
		message:    config.Message,
		timeout:    config.Timeout,
		attempts:   config.Attempts,
		period:     config.Period,
		succThrPec: config.SuccThrPec,
	}
	ct.logger = logger
	return ct, nil
}

// Start the client
func (c *ConnTester) Start(ctx context.Context, wait chan error) {
	defer close(wait)

	ticker := time.NewTicker(time.Duration(c.period) * time.Millisecond)
	wg := sync.WaitGroup{}

	send := func() {
		wg.Add(1)
		go c.send(ctx, &wg)
	}

	cleanandwait := func() {
		c.logger.Info("waiting for connections to stop...")
		ticker.Stop()
		wg.Wait()
		c.logger.Info("all connections have finished")
	}

	send()
	for count := 1; count < int(c.attempts); {
		select {
		case <-ctx.Done():
			cleanandwait()
			wait <- context.Canceled
			return
		case <-ticker.C:
			count++
			send()
		}
	}

	c.logger.Info("done creating connections")
	cleanandwait()

	if err := ctx.Err(); err != nil {
		wait <- err
		return
	}

	succperc := (c.succCounter * 100) / c.attempts
	c.logger.Info(fmt.Sprintf("success rate of: %d", succperc))
	if succperc < c.succThrPec {
		c.logger.Info(fmt.Sprintf("success rate lower than threshold: %d", c.succThrPec))
		wait <- ErrTestFailed
		return
	}
	c.logger.Info(fmt.Sprintf("success rate greater than threshold: %d", c.succThrPec))
	wait <- nil
}

func (c *ConnTester) send(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	dialer := net.Dialer{Timeout: time.Duration(c.timeout) * (time.Millisecond)}
	var conn net.Conn
	var err error
	protocol := strings.ToLower(c.protocol)
	dstEndpoint := fmt.Sprintf("%s:%d", c.targetHost, c.targetPort)

	if conn, err = dialer.DialContext(ctx, protocol, dstEndpoint); err != nil {
		c.logger.Info(fmt.Sprintf("cannot connect to %s: %s", dstEndpoint, err))
		return
	}

	if err = conn.SetDeadline(time.Now().Add(time.Duration(c.timeout) * time.Millisecond)); err != nil {
		c.logger.Info(fmt.Sprintf("cannot set deadline to connection to %s: %s", dstEndpoint, err))
		return
	}
	if _, err := conn.Write([]byte(c.message)); err != nil {
		c.logger.Info(fmt.Sprintf("cannot send data to %s: %s", dstEndpoint, err))
		return
	}
	c.logger.Info(fmt.Sprintf("successful connection and data sent to %s", dstEndpoint))

	if err := conn.Close(); err != nil {
		c.logger.Info(fmt.Sprintf("error while closing connection to %s: %s", dstEndpoint, err))
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.succCounter++
}
