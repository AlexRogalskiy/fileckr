package codec

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"math"
	"os"

	fileckrmath "github.com/squat/fileckr/pkg/math"
)

// EncodeFile converts a file to a PNG and writes the bytes to the given Writer.
func EncodeFile(file string, w io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}
	s := fi.Size()
	return Encode(bufio.NewReader(f), uint64(s), w)
}

// DecodeFile converts a PNG to a file and writes the bytes to the given Writer.
func DecodeFile(file string, w io.Writer) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	return Decode(bufio.NewReader(f), w)
}

// Encode converts a file to a PNG and writes the bytes to the given Writer.
func Encode(r io.Reader, size uint64, w io.Writer) error {
	lenBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(lenBytes, size)
	r = io.MultiReader(bytes.NewReader(lenBytes), r)

	width, height := fileckrmath.Squarest(int(math.Ceil(float64(size+8) / 4)))
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	b := make([]byte, 4)
	z := []byte{0, 0, 0, 0}
Loop:
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			n, err := r.Read(b)
			copy(b[n:], z[n:])
			img.Set(x, y, color.NRGBA{
				R: b[0],
				G: b[1],
				B: b[2],
				A: b[3],
			})
			if err != nil {
				if err != io.EOF {
					return fmt.Errorf("failed to read from file: %v", err)
				}
				break Loop
			}
		}
	}

	p := png.Encoder{CompressionLevel: png.BestCompression}
	if err := p.Encode(w, img); err != nil {
		return fmt.Errorf("failed to encode png: %v", err)
	}
	return nil
}

// Decode converts a PNG to a file and writes the bytes to the given Writer.
func Decode(r io.Reader, w io.Writer) error {
	img, err := png.Decode(r)
	if err != nil {
		return fmt.Errorf("failed to decode png: %v", err)
	}

	var nb []uint8
	var n uint64
	buf := make([]byte, 4)
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			im := img.At(x, y).(color.NRGBA)
			buf = []byte{im.R, im.G, im.B, im.A}
			if len(nb) < 8 {
				nb = append(nb, buf...)
				if len(nb) == 8 {
					n = binary.LittleEndian.Uint64(nb)
				}
			} else {
				if n < 4 {
					_, err = w.Write(buf[:n])
					if err != nil {
						return fmt.Errorf("failed to write final data: %v", err)
					}
					return nil
				}
				_, err = w.Write(buf)
				if err != nil {
					return fmt.Errorf("failed to write data: %v", err)
				}
				n -= 4
			}
			// Special case: if file length is zero, return.
			if len(nb) == 8 && n == 0 {
				return nil
			}
		}
	}
	fmt.Println("length: ", n)
	return errors.New("failed to read all image data")
}
