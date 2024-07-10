# CashDrawer (cd)

`cd` is a Go library designed to interface with the "USB TRIGGER for CASH DRAWER" BT-100U. It simplifies the process of opening a cash drawer using Go code.

## Installation

To use the `cd` library in your Go project, you need to install it using:

```sh
go get github.com/Nimon77/cd
```

## Usage

Here is a quick guide on how to use the `cd` library to control your cash drawer.

### Opening the Cash Drawer

```go
package main

import (
    "context"
    "log"
    "github.com/Nimon77/cd"
)

func main() {
    ctx := context.Background()
    
    // Initialize the CashDrawer manually with a specific port and baud rate
    drawer, err := cd.New("/dev/ttyUSB0", 9600)
    if err != nil {
        log.Fatal(err)
    }
    defer drawer.Close()

    // Open the cash drawer
    err = drawer.Open(ctx)
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Cash drawer opened successfully!")
}
```

### Automatically Detecting and Opening the Cash Drawer (linux only)

If you prefer the library to automatically detect the correct USB port, you can use the `NewAuto` function.

```go
package main

import (
    "log"
    "github.com/Nimon77/cd"
)

func main() {
    // Automatically detect the cash drawer and initialize it
    drawer, err := cd.NewAuto()
    if err != nil {
        log.Fatal(err)
    }
    defer drawer.Close()

    // Open the cash drawer
    err = drawer.Open(context.Background())
    if err != nil {
        log.Fatal(err)
    }

    log.Println("Cash drawer opened successfully!")
}
```

## Functions

### `New(port string, baud int) (*CashDrawer, error)`

Creates a new `CashDrawer` instance with the specified serial port and baud rate.

- `port`: The serial port to which the cash drawer is connected (e.g., `/dev/ttyUSB0`).
- `baud`: The baud rate for the serial communication (e.g., 9600).

### `NewAuto() (*CashDrawer, error)` (linux only)

Automatically detects the USB port connected to the cash drawer and initializes a new `CashDrawer` instance. This function scans the available USB devices and matches the correct one using predefined rules.

### `Open(ctx context.Context) error`

Sends the command to open the cash drawer. It writes the necessary bytes to the serial port to trigger the drawer to open.

- `ctx`: The context for managing request deadlines and cancellation signals.

### `Close() error`

Closes the serial port connection to the cash drawer.

## Dependencies

This library depends on the following packages:

- `github.com/citilinkru/libudev` - For scanning and detecting USB devices. (linux only)
- `github.com/jacobsa/go-serial/serial` - For handling serial port communications.

Make sure to include these dependencies in your project.

## License

This project is licensed under the MIT License. See the LICENSE file for more details.

## Contributions

Contributions are welcome! Please fork the repository and submit a pull request with your changes.

## Support

If you encounter any issues or have questions, feel free to open an issue on GitHub.
