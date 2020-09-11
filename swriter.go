package swriter

import (
	"fmt"
)

const initBuffSize = 64
const maxInt = int(^uint(0) >> 1)

// SWriter ...
type SWriter struct {
	buf []byte
	off int // next write at buf[off]
}

// NewSWriter ...
func NewSWriter(n int) *SWriter {
	buf := makeSlice(n)
	buffer := &SWriter{buf: buf}
	buffer.Reset()
	return buffer
}

// Reset ...
func (b *SWriter) Reset() {
	b.buf = b.buf[:0]
	b.off = 0
}

// Bytes ...
func (b *SWriter) Bytes() []byte {
	return b.buf[:]
}

// CopyTo ...
func (b *SWriter) CopyTo(dst []byte) {
	copy(dst, b.buf)
}

// Len ...
func (b *SWriter) Len() int {
	return len(b.buf)
}

// Cap ...
func (b *SWriter) Cap() int {
	return cap(b.buf)
}

// Seek ...
func (b *SWriter) Seek(n int) {
	b.off = n
}

// Tell ...
func (b *SWriter) Tell() int {
	return b.off
}

// Grow ...
func (b *SWriter) Grow(n int) {
	if n < 0 {
		panic(fmt.Errorf("SWriter.Grow: negative count %d", n))
	}
	b.grow(n)
}

// Write ...
func (b *SWriter) Write(p []byte) (n int, err error) {
	if growSize := b.neededGrowSize(b.off, len(p)); 0 < growSize {
		if ok := b.tryGrowByReslice(growSize); !ok {
			b.grow(growSize)
		}
	}

	n = b.writeToBuff(p)

	return n, nil
}

// WriteAt ...
func (b *SWriter) WriteAt(p []byte, off int64) (n int, err error) {
	b.off = int(off)
	return b.Write(p)
}

func ceil(v, unit int) int {
	return (v + (unit - 1)) / unit * unit
}

func (b *SWriter) grow(n int) {
	/*
	 * ex.)
	 * len(b.buf) is 2, cap(b.buf) is 6, n is 3.
	 * before : xxoooo
	 * after  : xxxxxo
	 * x is used index, o is free index.
	 */

	newLen := len(b.buf) + n

	if ok := b.tryGrowByReslice(n); ok {
		return
	}

	oldCap := cap(b.buf)
	var unit int
	if oldCap == 0 {
		unit = initBuffSize
	} else {
		unit = oldCap
	}
	newCap := ceil(newLen, unit)

	newBuf := makeSlice(newCap)
	copy(newBuf, b.buf)

	b.buf = newBuf[:newLen]
	return
}

func (b *SWriter) tryGrowByReslice(n int) bool {
	if l := len(b.buf); n <= cap(b.buf)-l { // avoid over flow
		b.buf = b.buf[:l+n]
		return true
	}
	return false
}

func (b *SWriter) neededGrowSize(off, writeSize int) int {
	bufLen := len(b.buf)

	if off+writeSize <= bufLen {
		return 0
	}

	n := writeSize
	if bufLen <= off {
		n += int(off - bufLen)
	} else {
		n -= int(bufLen - off) // (bufLen - off) is free space in buffer
	}

	if n < 0 {
		n = 0
	}

	return n
}

func (b *SWriter) writeToBuff(p []byte) (n int) {
	n = copy(b.buf[b.off:], p)
	b.off += n
	return n
}

func makeSlice(n int) []byte {
	defer func() {
		if recover() != nil {
			panic(fmt.Errorf("SWriter: makeSlice too large %x", n))
		}
	}()
	return make([]byte, n)
}
