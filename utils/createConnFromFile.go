package utils

import (
	"net"
	"os"
	"time"
)

// FileConn Structure to imitate file is incoming connection
type FileConn struct {
	inputadd   FileAddres
	outputadd  FileAddres
	inputfile  *os.File
	outputfile *os.File
}

// FileAddres File Address
type FileAddres struct {
	filename string
}

// Network return networking type
func (FileAddres) Network() string {
	return "file"
}

func (f FileAddres) String() string {
	return f.filename
}

// OpenFileConn Opens file to imitate net.Conn interface
func OpenFileConn(inputfilename, outputfilename string) (FileConn, error) {
	inputfile, err := os.Open(inputfilename)
	if err != nil {
		return FileConn{}, err
	}
	outputfile, err := os.OpenFile(outputfilename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return FileConn{}, err
	}
	return FileConn{
		inputadd: FileAddres{
			inputfilename,
		},
		outputadd: FileAddres{
			outputfilename,
		},
		inputfile:  inputfile,
		outputfile: outputfile,
	}, nil
}

func (f FileConn) Read(b []byte) (n int, err error) {
	return f.inputfile.Read(b)
}

func (f FileConn) Write(b []byte) (n int, err error) {
	return f.outputfile.Write(b)
}

// Close closing file
func (f FileConn) Close() error {
	err := f.inputfile.Close()
	err2 := f.outputfile.Close()
	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// LocalAddr return localaddress
func (f FileConn) LocalAddr() net.Addr {
	return f.outputadd
}

// RemoteAddr return addr with src filename
func (f FileConn) RemoteAddr() net.Addr {
	return f.inputadd
}

// SetDeadline do nothing
func (FileConn) SetDeadline(t time.Time) error {
	return nil
}

// SetReadDeadline do nothing
func (FileConn) SetReadDeadline(t time.Time) error {
	return nil
}

// SetWriteDeadline do nothing
func (FileConn) SetWriteDeadline(t time.Time) error {
	return nil
}
