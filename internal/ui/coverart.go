package ui

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/nfnt/resize"
)

func decodeImage(data []byte) (image.Image, error) {
	if len(data) == 0 {
		return nil, nil
	}

	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	return img, nil
}

func toGrayscale(c color.Color) uint8 {
	r, g, b, _ := c.RGBA()
	return uint8((uint32(r)*19595 + uint32(g)*38470 + uint32(b)*7471) >> 16)
}

type CoverArtRenderer struct{}

func NewCoverArtRenderer() *CoverArtRenderer {
	return &CoverArtRenderer{}
}

func (r *CoverArtRenderer) Render(data []byte, width, height int) (string, error) {
	if len(data) == 0 || width <= 0 || height <= 0 {
		return "", nil
	}

	img, err := decodeImage(data)
	if err != nil || img == nil {
		return "", err
	}

	scaled := resize.Resize(uint(width*2), uint(height), img, resize.Lanczos3)

	charLines := make([]string, height)

	for y := 0; y < height; y++ {
		lineChars := make([]rune, width)

		for x := 0; x < width; x++ {
			topPixel := color.Gray{}
			bottomPixel := color.Gray{}

			pixelX := x * 2
			pixelYTop := y
			pixelYBottom := y + 1

			if pixelYTop < int(scaled.Bounds().Max.Y) {
				topPixel = color.Gray{Y: toGrayscale(scaled.At(int(scaled.Bounds().Min.X)+pixelX, int(scaled.Bounds().Min.Y)+pixelYTop))}
			}

			if pixelYBottom < int(scaled.Bounds().Max.Y) {
				bottomPixel = color.Gray{Y: toGrayscale(scaled.At(int(scaled.Bounds().Min.X)+pixelX, int(scaled.Bounds().Min.Y)+pixelYBottom))}
			}

			topBright := topPixel.Y > 127
			bottomBright := bottomPixel.Y > 127

			var char rune
			if topBright && bottomBright {
				char = ' '
			} else if topBright && !bottomBright {
				char = '▁'
			} else if !topBright && bottomBright {
				char = '▀'
			} else {
				char = '█'
			}

			lineChars[x] = char
		}

		charLines[y] = string(lineChars)
	}

	return strings.Join(charLines, "\n"), nil
}
