package bytes

import (
	"io"
)

type BytesWriter struct {
	w io.Writer
}

func NewWriter(w io.Writer) *BytesWriter {
	return &BytesWriter{w}
}

func (b *BytesWriter) WriteInt8(data int8) {
	b.w.Write([]byte{
		byte(data),
	})
}

func (b *BytesWriter) WriteInt16(data int16) {
	b.w.Write([]byte{
		byte((data >> 8) & 0xFF),
		byte(data & 0xFF),
	})
}

func (b *BytesWriter) WriteInt32(data int32) {
	b.w.Write([]byte{
		byte((data >> 24) & 0xFF),
		byte((data >> 16) & 0xFF),
		byte((data >> 8) & 0xFF),
		byte((data) & 0xFF),
	})
}

func (b *BytesWriter) WriteInt64(data int64) {
	b.w.Write([]byte{
		byte((data >> 54) & 0xFF),
		byte((data >> 48) & 0xFF),
		byte((data >> 40) & 0xFF),
		byte((data >> 32) & 0xFF),
		byte((data >> 24) & 0xFF),
		byte((data >> 16) & 0xFF),
		byte((data >> 8) & 0xFF),
		byte(data & 0xFF),
	})
}

func (b *BytesWriter) WriteString(data string) {
	for _, c := range data {
		b.w.Write([]byte{byte(c)})
	}
}
