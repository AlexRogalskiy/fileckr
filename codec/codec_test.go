package codec

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"
)

func TestCodec(t *testing.T) {
	type testCase struct {
		data []byte
		mode os.FileMode
	}
	cases := []testCase{
		{
			data: []byte(""),
			mode: 0600,
		},
		{
			data: []byte("1"),
			mode: 0777,
		},
		{
			data: []byte("foobarbaz"),
			mode: 0555,
		},
		{
			data: []byte("lorempsum"),
			mode: 0644,
		},
	}

	for i, c := range cases {
		img := new(bytes.Buffer)
		out := new(bytes.Buffer)
		err := Encode(bytes.NewBuffer(c.data), uint64(len(c.data)), c.mode, img)
		if err != nil {
			t.Errorf("test case %d: expected no error while encoding, got: %v", i, err)
		}
		mode, err := Decode(img, out)
		if err != nil {
			t.Errorf("test case %d: expected no error while decoding, got: %v", i, err)
		}
		if bytes.Compare(c.data, out.Bytes()) != 0 {
			t.Errorf("test case %d: expected %v, got: %v", i, c, out.Bytes())
		}
		if mode != c.mode {
			t.Errorf("test case %d: expected mode to be %o, got: %o", i, c.mode, mode)
		}
	}
}

func TestCodecFile(t *testing.T) {
	type testCase struct {
		data []byte
		mode os.FileMode
	}
	cases := []testCase{
		{
			data: []byte("temporary file's content"),
			mode: 0600,
		},
		{
			data: []byte("1"),
			mode: 0777,
		},
		{
			data: []byte(""),
			mode: 0555,
		},
	}

	for i, c := range cases {
		out := new(bytes.Buffer)
		f, err := ioutil.TempFile("", "fileckr-file")
		if err != nil {
			t.Errorf("test case %d: failed to create temporary file: %v", i, err)
		}
		err = os.Chmod(f.Name(), c.mode)
		if err != nil {
			t.Errorf("test case %d: failed to chmod temporary file: %v", i, err)
		}
		defer os.Remove(f.Name())
		if _, err = f.Write(c.data); err != nil {
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
		mode, err := DecodeFile(p.Name(), out)
		if err != nil {
			t.Errorf("test case %d: expected no error while decoding, got: %v", i, err)
		}
		if bytes.Compare(c.data, out.Bytes()) != 0 {
			t.Errorf("test case %d: expected %v, got: %v", i, c, out.Bytes())
		}
		if mode != c.mode {
			t.Errorf("test case %d: expected mode to be %o, got: %o", i, c.mode, mode)
		}
	}
}
