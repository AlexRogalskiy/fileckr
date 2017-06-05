package codec

import (
	"bytes"
	"testing"
)

func TestCodec(t *testing.T) {
	cases := [][]byte{
		[]byte(""),
		[]byte("1"),
		[]byte("foobarbaz"),
		[]byte("lorempsum"),
	}

	var err error
	for i, c := range cases {
		img := new(bytes.Buffer)
		out := new(bytes.Buffer)
		err = Encode(bytes.NewBuffer(c), uint64(len(c)), img)
		if err != nil {
			t.Errorf("test case %d: expected no error while encoding, got: %v", i, err)
		}
		err := Decode(img, out)
		if err != nil {
			t.Errorf("test case %d: expected no error while decoding, got: %v", i, err)
		}
		if bytes.Compare(c, out.Bytes()) != 0 {
			t.Errorf("test case %d: expected %v, got: %v", i, c, out.Bytes)
		}
	}
}
