package cd

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/kr/pty"
)

// FixedSizeBuffer is a custom implementation of io.ReadWriteCloser
// that does not grow and results in a short write if there is not enough space.
type FixedSizeBuffer struct {
	buffer bytes.Buffer
	size   int
	closed bool
}

// NewFixedSizeBuffer initializes a new FixedSizeBuffer with a given size.
func NewFixedSizeBuffer(size int) *FixedSizeBuffer {
	return &FixedSizeBuffer{
		size: size,
	}
}

// Read reads data from the FixedSizeBuffer into p.
// It returns the number of bytes read and any error encountered.
func (b *FixedSizeBuffer) Read(p []byte) (n int, err error) {
	return b.buffer.Read(p)
}

// Write writes data to the FixedSizeBuffer from p.
// If there is not enough space to write all data, it writes as much as it can
// and returns io.ErrShortWrite.
func (b *FixedSizeBuffer) Write(p []byte) (n int, err error) {
	if b.buffer.Len()+len(p) > b.size {
		n, _ = b.buffer.Write(p[:b.size-b.buffer.Len()])
		return n, io.ErrShortWrite
	}
	return b.buffer.Write(p)
}

// Close closes the FixedSizeBuffer, returning an error if it is already closed.
func (b *FixedSizeBuffer) Close() error {
	if b.closed {
		return errors.New("already closed")
	}
	b.closed = true
	return nil
}

// createMockPTY creates a mock PTY device and returns its paths.
func createMockPTY(t *testing.T) (string, *os.File, *os.File) {
	master, slave, err := pty.Open()
	if err != nil {
		t.Fatalf("failed to create PTY: %v", err)
	}
	return slave.Name(), master, slave
}

// TestCashDrawer_Open tests the CashDrawer Open method for successful operation.
// It initializes a FixedSizeBuffer with sufficient size and verifies the output.
func TestCashDrawer_Open(t *testing.T) {
	mockSerialPort := NewFixedSizeBuffer(64) // buffer size larger than the needed write size (64 is the real default buffer size)
	drawer := &CashDrawer{
		Context:    context.Background(),
		serialPort: mockSerialPort,
	}

	err := drawer.Open(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedOutput := "\x1B\x70\x00\x30"
	if mockSerialPort.buffer.String() != expectedOutput {
		t.Errorf("expected %v, got %v", expectedOutput, mockSerialPort.buffer.String())
	}
}

// TestCashDrawer_Open_ShortWrite tests the CashDrawer Open method for short write error.
// It initializes a FixedSizeBuffer with insufficient size and verifies the error.
func TestCashDrawer_Open_ShortWrite(t *testing.T) {
	mockSerialPort := NewFixedSizeBuffer(3) // buffer size smaller than the needed write size
	drawer := &CashDrawer{
		Context:    context.Background(),
		serialPort: mockSerialPort,
	}

	err := drawer.Open(context.Background())
	if err != io.ErrShortWrite {
		t.Fatalf("expected io.ErrShortWrite, got %v", err)
	}
}

// TestCashDrawer_Close tests the CashDrawer Close method.
func TestCashDrawer_Close(t *testing.T) {
	ptyPath, masterFile, slaveFile := createMockPTY(t)
	defer masterFile.Close()
	defer slaveFile.Close()

	drawer, err := New(ptyPath, 9600)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	err = drawer.Close()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Try closing again to check for error
	err = drawer.Close()
	if err == nil {
		t.Fatal("expected an error on second close, got nil")
	}
}

// TestNew tests the New function.
func TestNew(t *testing.T) {
	ptyPath, masterFile, slaveFile := createMockPTY(t)
	defer masterFile.Close()
	defer slaveFile.Close()

	drawer, err := New(ptyPath, 9600)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if drawer.serialPort == nil {
		t.Fatalf("expected serialPort to be initialized, got nil")
	}
}

// TestNewInvalidPTY tests the New function.
func TestNewInvalidPTY(t *testing.T) {
	_, err := New("/dev/notExistingFile_3b3", 1337)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
