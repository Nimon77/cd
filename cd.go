// cd is a Go library designed to interface with the "USB TRIGGER for CASH DRAWER" BT-100U. It simplifies the process of opening a cash drawer using Go code.
//
// ## Installation
//
// To use the `cd` library in your Go project, you need to install it using:
//
//	go get github.com/Nimon77/cd
//
// ## Usage
//
// Here is a quick guide on how to use the `cd` library to control your cash drawer.
//
// ### Opening the Cash Drawer
//
//	package main
//
//	import (
//
//	 "context"
//	 "log"
//	 "github.com/Nimon77/cd"
//
//	)
//
//	func main() {
//	    ctx := context.Background()
//
//	    // Initialize the CashDrawer manually with a specific port and baud rate
//	    drawer, err := cd.New("/dev/ttyUSB0", 9600)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer drawer.Close()
//
//	    // Open the cash drawer
//	    err = drawer.Open(ctx)
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    log.Println("Cash drawer opened successfully!")
//	}
//
// ### Automatically Detecting and Opening the Cash Drawer (linux only)
//
// If you prefer the library to automatically detect the correct USB port, you can use the `NewAuto` function.
//
//	package main
//
//	import (
//
//	 "log"
//	 "github.com/Nimon77/cd"
//
//	)
//
//	func main() {
//	    // Automatically detect the cash drawer and initialize it
//	    drawer, err := cd.NewAuto()
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//	    defer drawer.Close()
//
//	    // Open the cash drawer
//	    err = drawer.Open(context.Background())
//	    if err != nil {
//	        log.Fatal(err)
//	    }
//
//	    log.Println("Cash drawer opened successfully!")
//	}
//
// ## Dependencies
//
// This library depends on the following packages:
//
// - `github.com/citilinkru/libudev` - For scanning and detecting USB devices. (linux only)
// - `github.com/jacobsa/go-serial/serial` - For handling serial port communications.
//
// ## License
//
// This project is licensed under the MIT License. See the LICENSE file for more details.
//
// ## Contributions
//
// Contributions are welcome! Please fork the repository and submit a pull request with your changes.
//
// ## Support
//
// If you encounter any issues or have questions, feel free to open an issue on GitHub.
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
	_, err := io.WriteString(r.serialPort, "\x1B\x70\x00\x30")
	if err != nil {
		return err
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
