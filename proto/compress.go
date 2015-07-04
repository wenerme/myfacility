package proto
import (
	"io"
	"bytes"
	"encoding/binary"
	"compress/zlib"
)

type compressedReader struct {
	r   io.Reader
	buf *bytes.Buffer
}

type compressedWriter struct {
	w   io.Writer
	buf *bytes.Buffer
	seq uint8
}

func NewCompressedReader(r io.Reader) io.Reader {
	c := &compressedReader{r:r}
	c.buf = &bytes.Buffer{}
	return c
}
func NewCompressedWriter(w io.Writer) io.Writer {
	c := &compressedWriter{w:w}
	c.buf = &bytes.Buffer{}
	return c
}

func (r *compressedReader)Read(p []byte) (n int, err error) {
	if r.buf.Len() < len(p) {
		// read more
		tmp := make([]byte, 4)
		_, err = r.r.Read(tmp[:3])
		if err != nil {return }

		compressedLen := binary.LittleEndian.Uint32(tmp)
		_, err = r.r.Read(tmp[:1])
		if err != nil {return }
		seq := uint8(tmp[3])
		_, err = r.r.Read(tmp[:3])
		if err != nil {return }
		beforeLen := binary.LittleEndian.Uint32(tmp)

		//		fmt.Printf("CL: %d BL: %d SEQ: %d \n", compressedLen, beforeLen, seq)

		if beforeLen == 0 {
			_, err = io.CopyN(r.buf, r.r, int64(compressedLen))
		}else {
			var zr io.ReadCloser
			zr, err = zlib.NewReader(r.r)
			if err != nil {return }
			_, err = io.Copy(r.buf, zr)
		}

		if err != nil {return }
		_, _ = seq, compressedLen
	}
	n, err = r.buf.Read(p)
	return
}

func (w *compressedWriter)Write(p []byte) (written int, err error) {
	// Must write as a packet
	if p[4] == 0 { w.seq = 0 }

	tmp := make([]byte, 4)
	n := 0
	// uncompressed
	if len(p) < MIN_COMPRESS_LENGTH {
		binary.LittleEndian.PutUint32(tmp, uint32(len(p)))
		n, err = w.w.Write(tmp[:3])
		written += n
		if err != nil {return }
		tmp[0] = w.seq
		w.seq ++
		n, err = w.w.Write(tmp[:1])
		written += n
		if err != nil {return }
		binary.LittleEndian.PutUint32(tmp, 0)
		n, err = w.w.Write(tmp[:3])
		written += n
		if err != nil {return }
		n, err = w.w.Write(p)
		written += n
	}else {
		zw := zlib.NewWriter(w.buf)
		zw.Write(p)
		zw.Close()

		binary.LittleEndian.PutUint32(tmp, uint32(w.buf.Len()))
		n, err = w.w.Write(tmp[:3])
		written += n
		if err != nil {return }
		tmp[0] = w.seq
		w.seq ++
		n, err = w.w.Write(tmp[:1])
		written += n
		if err != nil {return }
		binary.LittleEndian.PutUint32(tmp, uint32(len(p)))
		n, err = w.w.Write(tmp[:3])
		written += n
		if err != nil {return }
		n, err = w.w.Write(w.buf.Bytes())
		written += n
		w.buf.Reset()
	}
	return
}