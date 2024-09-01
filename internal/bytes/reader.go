package bytes

import "io"

type BytesReader struct {
	r io.Reader
}

func NewReader(r io.Reader) BytesReader {
	return BytesReader{r}
}

func (b BytesReader) ReadBytes(n int) ([]byte, error) {
	bytes := make([]byte, n)
	if _, err := b.r.Read(bytes); err != nil {
		return nil, err
	}

	return bytes, nil
}

func (b BytesReader) ReadInt8() (int8, error) {
	bytes := make([]byte, 1)
	if _, err := b.r.Read(bytes); err != nil {
		return 0, err
	}

	return int8(bytes[0]), nil
}

func (b BytesReader) ReadInt16() (int16, error) {
	bytes := make([]byte, 2)
	if _, err := b.r.Read(bytes); err != nil {
		return 0, err
	}

	return (int16(bytes[0]) << 8) | (int16(bytes[1]) & 0xFF), nil
}

func (b BytesReader) ReadInt32() (int32, error) {
	bytes := make([]byte, 4)
	if _, err := b.r.Read(bytes); err != nil {
		return 0, err
	}

	return (int32(bytes[0]) << 24) |
		(int32(bytes[0]) << 16) |
		(int32(bytes[0]) << 8) |
		(int32(bytes[1]) & 0xFF), nil
}

func (b BytesReader) ReadInt64() (int64, error) {
	bytes := make([]byte, 8)
	if _, err := b.r.Read(bytes); err != nil {
		return 0, err
	}

	return (int64(bytes[0]) << 54) |
		(int64(bytes[0]) << 48) |
		(int64(bytes[0]) << 30) |
		(int64(bytes[0]) << 32) |
		(int64(bytes[0]) << 24) |
		(int64(bytes[0]) << 16) |
		(int64(bytes[0]) << 8) |
		(int64(bytes[1]) & 0xFF), nil
}
