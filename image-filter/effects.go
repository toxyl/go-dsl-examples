package main

import (
	"fmt"
	"image"
	"image/color"

	"github.com/toxyl/math"
)

// @Name: invert
// @Desc: Inverts an image
// @Param:      img     - -   -   The image to invert
// @Returns:    result  - -   -   The inverted image
func invert(img *image.NRGBA) (any, error) {
	bounds := img.Bounds()
	inverted := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			inverted.Set(x, y, color.NRGBA{
				R: 255 - c.R,
				G: 255 - c.G,
				B: 255 - c.B,
				A: c.A, // Keep original alpha
			})
		}
	}
	return inverted, nil
}

// @Name: grayscale
// @Desc: Grayscales an image
// @Param:      img     - -   -   The image to grayscale
// @Returns:    result  - -   -   The grayscaled image
func grayscale(img *image.NRGBA) (any, error) {
	bounds := img.Bounds()
	grayscaled := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			// Using luminosity method: 0.21 R + 0.72 G + 0.07 B
			gray := uint8(float64(c.R)*0.21 + float64(c.G)*0.72 + float64(c.B)*0.07)
			grayscaled.Set(x, y, color.NRGBA{
				R: gray,
				G: gray,
				B: gray,
				A: c.A,
			})
		}
	}
	return grayscaled, nil
}

// @Name: sepia
// @Desc: Changes the tone of an image to sepia
// @Param:      img     - -   -   The image to change to sepia tone
// @Returns:    result  - -   -   The sepia-toned image
func sepia(img *image.NRGBA) (any, error) {
	bounds := img.Bounds()
	sepia := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			r := float64(c.R)
			g := float64(c.G)
			b := float64(c.B)

			newR := uint8(math.Min((r*0.393)+(g*0.769)+(b*0.189), 255))
			newG := uint8(math.Min((r*0.349)+(g*0.686)+(b*0.168), 255))
			newB := uint8(math.Min((r*0.272)+(g*0.534)+(b*0.131), 255))

			sepia.Set(x, y, color.NRGBA{
				R: newR,
				G: newG,
				B: newB,
				A: c.A,
			})
		}
	}
	return sepia, nil
}

// @Name: brightness
// @Desc: Changes the brightness of an image
// @Param:      img     - -   	-   The image to change brightness of
// @Param:      factor  - 0..2  0   The change factor
// @Returns:    result  - -   	-   The image with brightness changed
func brightness(img *image.NRGBA, factor float64) (any, error) {
	if factor < 0.0 || factor > 2.0 {
		return nil, fmt.Errorf("brightness factor must be between 0.0 and 2.0")
	}

	bounds := img.Bounds()
	adjusted := image.NewNRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.NRGBAAt(x, y)
			adjusted.Set(x, y, color.NRGBA{
				R: uint8(math.Min(float64(c.R)*factor, 255)),
				G: uint8(math.Min(float64(c.G)*factor, 255)),
				B: uint8(math.Min(float64(c.B)*factor, 255)),
				A: c.A,
			})
		}
	}
	return adjusted, nil
}

// @Name: fill
// @Desc: Fills the image
// @Param:      img     - - -   The image to fill
// @Param:      col  	- - -   The fill color
// @Returns:    result  - - -	The filled image
func fill(img *image.NRGBA, col color.RGBA64) (*image.NRGBA, error) {
	bounds := img.Bounds()
	filled := image.NewNRGBA(bounds)

	// Convert RGBA64 to NRGBA color
	nrgbaCol := color.NRGBA{
		R: uint8(col.R >> 8),
		G: uint8(col.G >> 8),
		B: uint8(col.B >> 8),
		A: uint8(col.A >> 8),
	}

	// Fill the entire image with the given color
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			filled.Set(x, y, nrgbaCol)
		}
	}
	return filled, nil
}
