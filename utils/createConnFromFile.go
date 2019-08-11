package utils

import (
	"log"
	"net"
	"os"
	"time"
)

// FileConn Structure to imitate file is incoming connection
type FileConn struct {
	inAddress  FileAddress
	outAddress FileAddress

	amountRead int64
	sizeOfFile int64

	inFile  *os.File
	outFile *os.File
}

// FileAddress File Address
type FileAddress struct {
	filename string
}

// Network return networking type
func (FileAddress) Network() string {
	return "file"
}

func (f FileAddress) String() string {
	return f.filename
}

// OpenFileConn Opens file to imitate net.Conn interface
func OpenFileConn(inputfn, outputfn string) (FileConn, error) {
	var inputFile, outputFile *os.File
	var err error

	inputFile, err = os.Open(inputfn)
	if err != nil {
		return FileConn{}, err
	}

	var sizeOfFile int64
	fi, err := inputFile.Stat()
	if err != nil {
		log.Println("Problem with reading file info, setting size to -1")
		sizeOfFile = -1
	} else {
		sizeOfFile = fi.Size()
	}

	if outputfn != "" {
		outputFile, err = os.OpenFile(outputfn, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return FileConn{}, err
		}
	}

	return FileConn{
		inAddress: FileAddress{
			inputfn,
		},
		outAddress: FileAddress{
			outputfn,
		},
		sizeOfFile: sizeOfFile,
		amountRead: 0,
		inFile:     inputFile,
		outFile:    outputFile,
	}, nil
}

func (f FileConn) Percent() float64 {
	return float64(f.amountRead) / float64(f.sizeOfFile)
}

func (f FileConn) IsRead() bool {
	return f.amountRead == f.sizeOfFile
}

func (f FileConn) AmountRead() int64 {
	return f.amountRead
}

func (f *FileConn) Seek(offset int64, whence int) (n int64, err error) {
	f.amountRead += offset
	return f.inFile.Seek(offset, whence)
}

func (f *FileConn) Read(b []byte) (n int, err error) {
	f.amountRead += int64(len(b))
	return f.inFile.Read(b)
}

func (f FileConn) Write(b []byte) (n int, err error) {
	return f.outFile.Write(b)
}

// Close closing file
func (f FileConn) Close() error {
	err := f.inFile.Close()
	err2 := f.outFile.Close()
	if err != nil {
		return err
	}
	if err2 != nil {
		return err2
	}
	return nil
}

// LocalAddr return local address
func (f FileConn) LocalAddr() net.Addr {
	return f.outAddress
}

// RemoteAddr return addr with src filename
func (f FileConn) RemoteAddr() net.Addr {
	return f.inAddress
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
