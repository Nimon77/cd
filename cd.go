package cd

import (
	"context"
	"io"

	"github.com/jacobsa/go-serial/serial"
)

// CashDrawer represents a cash drawer, with its serial port.
type CashDrawer struct {
	Context context.Context

	serialPort io.ReadWriteCloser
}

// Open sends the command to open the cash drawer. It writes the necessary bytes to the serial port to trigger the drawer to open.
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

// New initializes a new CashDrawer with the given serial port (i.e /dev/ttyUSB0) and baud rate (i.e 9600).
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

// Close closes the serial port connection to the cash drawer.
func (r *CashDrawer) Close() error {
	return r.serialPort.Close()
}
