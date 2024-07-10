package main

import (
	"context"
	"io"
	"log"

	"github.com/citilinkru/libudev"
	"github.com/citilinkru/libudev/matcher"
	"github.com/jacobsa/go-serial/serial"
)

type CashDrawer struct {
	Context context.Context

	serialPort io.ReadWriteCloser
}

func (r *CashDrawer) Open(ctx context.Context) error {
	i, err := io.WriteString(r.serialPort, "\x1B\x70\x00\x30")
	if err != nil {
		return err
	}
	if i != 4 {
		return io.ErrShortWrite
	}
	return nil
}

func New(port string, baud int) (*CashDrawer, error) {
	options := serial.OpenOptions{
		PortName:        port,
		BaudRate:        uint(baud),
		DataBits:        8,
		StopBits:        1,
		MinimumReadSize: 1,
	}
	s, err := serial.Open(options)
	if err != nil {
		return nil, err
	}
	return &CashDrawer{
		serialPort: s,
	}, nil
}

func NewAuto() (*CashDrawer, error) {
	sc := libudev.NewScanner()
	err, devices := sc.ScanDevices()
	if err != nil {
		log.Fatal(err)
	}

	m := matcher.NewMatcher()
	m.SetStrategy(matcher.StrategyAnd)
	m.AddRule(matcher.NewRuleEnv("ID_VENDOR", "Prolific_Technology_Inc."))
	m.AddRule(matcher.NewRuleEnv("DEVNAME", "tty"))

	filteredDevices := m.Match(devices)

	if len(filteredDevices) == 0 {
		return nil, io.EOF
	}

	return New("/dev/"+filteredDevices[0].Env["DevName"], 9600)

}

func (r *CashDrawer) Close() error {
	return r.serialPort.Close()
}
