package gzipjson

import (
	"compress/gzip"
	"encoding/json"
	"github.com/juju/errors"
	"io"
	"sync"
)

const defaultMinSize = 1400

type GzipWriteCloser struct {
	io.Writer
	gw             *gzip.Writer
	buf            []byte
	minContentSize int
	compressed     bool
}

var gzipWriterPool *sync.Pool

func init() {
	gzipWriterPool = &sync.Pool {
		New: func() interface{} {
			return gzip.NewWriter(nil)
		},
	}
}

func Compress(w io.Writer, v interface{}) (bool, error) {
	return CompressWitMinSize(w, v, defaultMinSize)
}

func CompressWitMinSize(w io.Writer, v interface{}, minSize int) (bool, error) {

	gwc := newGzipWriteCloser(w, minSize)

	if err := json.NewEncoder(gwc).Encode(v); err != nil {
		return false, err
	}

	if err := gwc.Close(); err != nil {
		return false, err
	}

	return gwc.compressed, nil
}

func newGzipWriteCloser(w io.Writer, minSize int) *GzipWriteCloser {

	ms := defaultMinSize

	if minSize > -1 {
		ms = minSize
	}

	return &GzipWriteCloser{
		Writer:         w,
		minContentSize: ms,
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
	// Gzip was not triggered (content too small)
	if gwc.gw == nil {
		return gwc.doNotGzip()
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
			err = errors.Wrap(io.ErrShortWrite, errors.New("doNotGzip: Writer.Write wrote less bytes than the size of the buffer"))
		}

		return err
	}

	return nil
}
