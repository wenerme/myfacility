package proto
import (
	"io"
	"bytes"
	"encoding/binary"
	"compress/zlib"
	"fmt"
)

type compressedReader struct {
	r          io.Reader
	zlibReader io.ReadCloser
	buf        *bytes.Buffer
}

func NewCompressedReader(r io.Reader) io.Reader {
	c := &compressedReader{r:r}
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

		fmt.Printf("CL: %d BL: %d SEQ: %d \n", compressedLen, beforeLen, seq)

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