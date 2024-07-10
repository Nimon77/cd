package main

import (
	"context"
	"io"

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

func (r *CashDrawer) Close() error {
	return r.serialPort.Close()
}
