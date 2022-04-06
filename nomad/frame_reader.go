package nomad

import (
	"bufio"
	"bytes"
	"io"
	"sync"
	"time"

	"github.com/hashicorp/nomad/api"
)

type StreamFrame struct {
	Name string
	*api.StreamFrame
}

// FrameReader is used to convert a stream of frames into a read closer.
type FrameReader struct {
	frames   <-chan *StreamFrame
	errCh    <-chan error
	cancelCh chan struct{}

	closedLock sync.Mutex
	closed     bool

	unblockTime time.Duration

	frame       *StreamFrame
	frameBytes  []byte
	frameOffset int

	byteOffset int
}

// NewFrameReader takes a channel of frames and returns a FrameReader which
// implements io.ReadCloser
func NewFrameReader(frames <-chan *StreamFrame, errCh <-chan error, cancelCh chan struct{}) *FrameReader {
	return &FrameReader{
		frames:   frames,
		errCh:    errCh,
		cancelCh: cancelCh,
	}
}

// SetUnblockTime sets the time to unblock and return zero bytes read. If the
// duration is unset or is zero or less, the read will block until data is read.
func (f *FrameReader) SetUnblockTime(d time.Duration) {
	f.unblockTime = d
}

// Read reads the data of the incoming frames into the bytes buffer. Returns EOF
// when there are no more frames.
func (f *FrameReader) Read(p []byte) (n int, err error) {
	f.closedLock.Lock()
	closed := f.closed
	f.closedLock.Unlock()
	if closed {
		return 0, io.EOF
	}

	if f.frame == nil {
		var unblock <-chan time.Time
		if f.unblockTime.Nanoseconds() > 0 {
			unblock = time.After(f.unblockTime)
		}

		select {
		case frame, ok := <-f.frames:
			if !ok {
				return 0, io.EOF
			}

			f.frame = frame

			// Prepend every log line
			buff := bytes.NewBuffer([]byte{})
			scanner := bufio.NewScanner(bytes.NewReader(frame.Data))
			for scanner.Scan() {
				buff.Write([]byte(f.frame.Name))
				buff.Write([]byte(": "))
				buff.Write(scanner.Bytes())
				buff.Write([]byte("\n"))
			}

			f.frameBytes = buff.Bytes()

			// Store the total offset into the file
			f.byteOffset = int(f.frame.Offset)
		case <-unblock:
			return 0, nil
		case err := <-f.errCh:
			return 0, err
		case <-f.cancelCh:
			return 0, io.EOF
		}
	}

	// Copy the data out of the frame and update our offset
	n = copy(p, f.frameBytes[f.frameOffset:])
	f.frameOffset += n

	// Clear the frame and its offset once we have read everything
	if len(f.frameBytes) == f.frameOffset {
		f.frame = nil
		f.frameOffset = 0
	}

	return n, nil
}

// Close cancels the stream of frames
func (f *FrameReader) Close() error {
	f.closedLock.Lock()
	defer f.closedLock.Unlock()
	if f.closed {
		return nil
	}

	close(f.cancelCh)
	f.closed = true
	return nil
}
