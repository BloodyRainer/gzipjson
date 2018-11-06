package gzipjson

import (
	"compress/gzip"
	"encoding/json"
	"github.com/juju/errors"
	"io"
	"sync"
)

// Content that has more bytes than DefaultMinSize is compressed if not configured otherwise.
const DefaultMinSize = 1400

type GzipWriteCloser struct {
	io.Writer
	gw             *gzip.Writer
	buf            []byte
	minContentSize int
	compressed     bool
	option         Option
}

type Option struct {
	MinSize              int
	compressedCallback   func()
	uncompressedCallback func()
}

var gzipWriterPool *sync.Pool

func init() {
	gzipWriterPool = &sync.Pool{
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
}

// Encodes the given reference to JSON and compresses it if the size exceeds 1400 Bytes.
func Compress(w io.Writer, v interface{}) (bool, error) {
	o := Option{
		MinSize:            DefaultMinSize,
		compressedCallback: nil,
	}
	return CompressWitOption(w, v, o)
}

func CompressWitOption(w io.Writer, v interface{}, o Option) (bool, error) {

	gwc := newGzipWriteCloser(w, o)

	if err := json.NewEncoder(gwc).Encode(v); err != nil {
		return false, err
	}

	if err := gwc.Close(); err != nil {
		return false, err
	}

	return gwc.compressed, nil
}

// Encodes the given reference to JSON and compresses it if the size exceeds the given minSize value of bytes.
func newGzipWriteCloser(w io.Writer, o Option) *GzipWriteCloser {

	ms := DefaultMinSize

	if o.MinSize > -1 {
		ms = o.MinSize
	}

	return &GzipWriteCloser{
		Writer:         w,
		minContentSize: ms,
		option:         o,
	}
}

func (gwc *GzipWriteCloser) Write(b []byte) (int, error) {

	// if it was decided to gzip the content, gw is not nil
	if gwc.gw != nil {
		return gwc.gw.Write(b)
	}

	gwc.buf = append(gwc.buf, b...)

	// if the length of the buffer does not exceed the minimum size of the content, wait for more data
	if len(gwc.buf) < gwc.minContentSize {
		return len(b), nil
	} else {
		gwc.compressed = true

		//same as gwc.gw = gzip.NewWriter(gwc.Writer), but with pool
		gwc.gw = gzipWriterPool.Get().(*gzip.Writer)
		gwc.gw.Reset(gwc.Writer)

		return gwc.gw.Write(gwc.buf)
	}

}

func (gwc *GzipWriteCloser) Close() error {
	// gzip was not triggered (content too small)
	if gwc.gw == nil {
		return gwc.doNotGzip()
	}

	if gwc.compressed && gwc.option.compressedCallback != nil{
		gwc.option.compressedCallback()
	} else if gwc.option.uncompressedCallback != nil {
		gwc.option.uncompressedCallback()
	}

	err := gwc.gw.Close()
	if err != nil {
		return errors.Wrap(err, errors.New("unable to close gzip writer"))
	}

	gzipWriterPool.Put(gwc.gw)

	return nil
}

func (gwc *GzipWriteCloser) doNotGzip() error {

	if gwc.buf != nil {
		n, err := gwc.Writer.Write(gwc.buf)

		if err == nil && n < len(gwc.buf) {
			err = errors.Wrap(io.ErrShortWrite, errors.New("wrote less bytes than the size of the buffer"))
		}

		return err
	}

	return nil
}
