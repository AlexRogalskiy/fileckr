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

// EncodeFile converts the named file to a PNG and writes the bytes to the given Writer.
func EncodeFile(name string, w io.Writer) error {
	f, err := os.Open(name)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()
	s, err := f.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}
	return Encode(bufio.NewReader(f), uint64(s.Size()), s.Mode(), w)
}

// DecodeFile converts the named PNG to a file and writes the bytes to the given Writer.
func DecodeFile(name string, w io.Writer) (os.FileMode, error) {
	f, err := os.Open(name)
	if err != nil {
		return 0, fmt.Errorf("failed to open file: %v", err)
	}
	return Decode(bufio.NewReader(f), w)
}

// Encode converts a file to a PNG and writes the bytes to the given Writer.
func Encode(r io.Reader, size uint64, mode os.FileMode, w io.Writer) error {
	sizeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(sizeBytes, size)
	modeBytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(modeBytes, uint32(mode))
	r = io.MultiReader(bytes.NewReader(sizeBytes), bytes.NewReader(modeBytes), r)

	width, height := fileckrmath.NiceSquarest(int(math.Ceil(float64(size+8+4) / 4)))
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	b := make([]byte, 4)
	z := [4]byte{}
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
func Decode(r io.Reader, w io.Writer) (os.FileMode, error) {
	var mode os.FileMode
	img, err := png.Decode(r)
	if err != nil {
		return mode, fmt.Errorf("failed to decode png: %v", err)
	}

	var sizeBytes []uint8
	var size uint64
	var setMode bool
	b := make([]byte, 4)
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y
	var n uint64
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			im := img.At(x, y).(color.NRGBA)
			b = []byte{im.R, im.G, im.B, im.A}
			if len(sizeBytes) < 8 {
				sizeBytes = append(sizeBytes, b...)
				if len(sizeBytes) == 8 {
					size = binary.LittleEndian.Uint64(sizeBytes)
				}
				continue
			} else if !setMode {
				mode = os.FileMode(binary.LittleEndian.Uint32(b))
				setMode = true
			} else {
				if size < 4 {
					n = size
				} else {
					n = 4
				}
				_, err = w.Write(b[:n])
				if err != nil {
					return mode, fmt.Errorf("failed to write data: %v", err)
				}
				size -= n
			}
			if len(sizeBytes) == 8 && size == 0 {
				return mode, nil
			}
		}
	}
	return mode, errors.New("failed to read all image data")
}
