package tcp_conn

import (
	"errors"
	"io"
)

var (
	ErrReadFailed = errors.New("error to read")
	ErrNotEnough  = errors.New("have no space")
)

type buffer struct {
	buf   []byte
	start int
	end   int
}

func newBuffer(bytes []byte) buffer {
	return buffer{
		buf:   bytes,
		start: 0,
		end:   0,
	}
}

func (b *buffer) len() int {
	return b.end - b.start
}

func (b *buffer) grow() {
	if b.start == 0 {
		return
	}
	copy(b.buf, b.buf[b.start:b.end])
	b.end -= b.start
	b.start = 0
}

func (b *buffer) loadFromReader(reader io.Reader) (int, error) {
	b.grow()
	n, err := reader.Read(b.buf[b.end:])
	if err != nil {
		return 0, ErrReadFailed
	}
	b.end += n
	return n, nil
}

func (b *buffer) seek(start, end int) ([]byte, error) {
	if b.end-b.start >= end-start {
		buf := b.buf[b.start+start : b.start+end]
		return buf, nil
	}
	return nil, ErrNotEnough
}

func (b *buffer) read(offest, limit int) ([]byte, error) {
	if b.len() < offest+limit {
		return nil, ErrNotEnough
	}
	buf := b.buf[b.start+offest : b.start+limit]
	b.start += offest + limit
	return buf, nil
}

