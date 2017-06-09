package codec

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCodec(t *testing.T) {
	cases := [][]byte{
		[]byte(""),
		[]byte("1"),
		[]byte("foobarbaz"),
		[]byte("lorempsum"),
	}

	for i, c := range cases {
		img := new(bytes.Buffer)
		out := new(bytes.Buffer)
		err := Encode(bytes.NewBuffer(c), uint64(len(c)), img)
		if err != nil {
			t.Errorf("test case %d: expected no error while encoding, got: %v", i, err)
		}
		err = Decode(img, out)
		if err != nil {
			t.Errorf("test case %d: expected no error while decoding, got: %v", i, err)
		}
		if bytes.Compare(c, out.Bytes()) != 0 {
			t.Errorf("test case %d: expected %v, got: %v", i, c, out.Bytes())
		}
	}
}

func TestCodecFile(t *testing.T) {
	cases := [][]byte{
		[]byte("temporary file's content"),
		[]byte("1"),
		[]byte(""),
	}

	for i, c := range cases {
		out := new(bytes.Buffer)
		f, err := ioutil.TempFile("", "fileckr-file")
		if err != nil {
			t.Errorf("test case %d: failed to create temporary file: %v", i, err)
		}
		defer os.Remove(f.Name())
		if _, err = f.Write(c); err != nil {
			t.Errorf("test case %d: failed to write to temporary file: %v", i, err)
		}
		if err = f.Close(); err != nil {
			t.Errorf("test case %d: failed to close temporary file: %v", i, err)
		}
		p, err := ioutil.TempFile("", "fileckr-png")
		if err != nil {
			t.Errorf("test case %d: failed to create temporary PNG: %v", i, err)
		}
		defer os.Remove(p.Name())
		err = EncodeFile(f.Name(), p)
		if err != nil {
			t.Errorf("test case %d: expected no error while encoding, got: %v", i, err)
		}
		if err = p.Close(); err != nil {
			t.Errorf("test case %d: failed to close temporary PNG: %v", i, err)
		}
		err = DecodeFile(p.Name(), out)
		if err != nil {
			t.Errorf("test case %d: expected no error while decoding, got: %v", i, err)
		}
		if bytes.Compare(c, out.Bytes()) != 0 {
			t.Errorf("test case %d: expected %v, got: %v", i, c, out.Bytes())
		}
	}
}
