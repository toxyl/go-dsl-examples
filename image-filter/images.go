package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/toxyl/math"
)

// Helper function to create a new RGBA64 image with the same bounds as the source
func createNewRGBA64FromBounds(img *image.RGBA64) *image.RGBA64 {
	return image.NewRGBA64(img.Bounds())
}

// Helper function to get RGBA64 components as uint32 for calculations
func getRGBA64Components(c color.RGBA64) (r, g, b, a uint32) {
	return uint32(c.R), uint32(c.G), uint32(c.B), uint32(c.A)
}

// Helper function to set RGBA64 color with clamped values
func setRGBA64Color(img *image.RGBA64, x, y int, r, g, b, a uint32) {
	img.Set(x, y, color.RGBA64{
		R: uint16(min(r, 0xffff)),
		G: uint16(min(g, 0xffff)),
		B: uint16(min(b, 0xffff)),
		A: uint16(min(a, 0xffff)),
	})
}

// Helper function for Porter-Duff alpha compositing
func porterDuffAlpha(a1, a2 uint32) uint32 {
	return a1 + a2 - ((a1 * a2) / 0xffff)
}

// @Name: load
// @Desc: Loads an image
// @Param:      path    - -   -   Path to the image
// @Returns:    result  - -   -   The loaded image
func load(path string) (any, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return png.Decode(file)
}

// @Name: save
// @Desc: Saves an image
// @Param:      img     - -   -    The image to save
// @Param:      path     - -   -   Path where to save
func save(img *image.NRGBA, path string) (any, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	png.Encode(file, img) // note that we use NRGBA because storing PNGs is much faster that way
	return img, nil
}

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

// @Name: blend-multiply
// @Desc: Blends the two images using the multiply blend mode
// @Param:      imgA     - -   	-   The bottom image
// @Param:      imgB     - -   	-   The top image
// @Returns:    result  - -   	-   The blended image
func blendMultiply(imgA *image.RGBA64, imgB *image.RGBA64) (any, error) {
	result := createNewRGBA64FromBounds(imgA)
	bounds := imgA.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := imgA.RGBA64At(x, y)
			c2 := imgB.RGBA64At(x, y)

			r1, g1, b1, a1 := getRGBA64Components(c1)
			r2, g2, b2, a2 := getRGBA64Components(c2)

			aOut := porterDuffAlpha(a1, a2)

			var rOut, gOut, bOut uint32
			if aOut > 0 {
				rOut = ((r1 * r2) + (r1 * (0xffff - a2)) + (r2 * (0xffff - a1))) / 0xffff
				gOut = ((g1 * g2) + (g1 * (0xffff - a2)) + (g2 * (0xffff - a1))) / 0xffff
				bOut = ((b1 * b2) + (b1 * (0xffff - a2)) + (b2 * (0xffff - a1))) / 0xffff
			}

			setRGBA64Color(result, x, y, rOut, gOut, bOut, aOut)
		}
	}

	return result, nil
}

// @Name: blend-screen
// @Desc: Blends the two images using the screen blend mode
// @Param:      imgA     - -   	-   The bottom image
// @Param:      imgB     - -   	-   The top image
// @Returns:    result  - -   	-   The blended image
func blendScreen(imgA *image.RGBA64, imgB *image.RGBA64) (any, error) {
	result := createNewRGBA64FromBounds(imgA)
	bounds := imgA.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := imgA.RGBA64At(x, y)
			c2 := imgB.RGBA64At(x, y)

			r1, g1, b1, a1 := getRGBA64Components(c1)
			r2, g2, b2, a2 := getRGBA64Components(c2)

			// Screen blend mode in premultiplied space: 1 - (1 - a) * (1 - b)
			r := 0xffff - ((0xffff - r1) * (0xffff - r2) / 0xffff)
			g := 0xffff - ((0xffff - g1) * (0xffff - g2) / 0xffff)
			b := 0xffff - ((0xffff - b1) * (0xffff - b2) / 0xffff)
			a := porterDuffAlpha(a1, a2)

			setRGBA64Color(result, x, y, r, g, b, a)
		}
	}

	return result, nil
}

// @Name: blend-exclusion
// @Desc: Blends the two images using the exclusion blend mode
// @Param:      imgA     - -   	-   The bottom image
// @Param:      imgB     - -   	-   The top image
// @Returns:    result  - -   	-   The blended image
func blendExclusion(imgA *image.RGBA64, imgB *image.RGBA64) (any, error) {
	result := createNewRGBA64FromBounds(imgA)
	bounds := imgA.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := imgA.RGBA64At(x, y)
			c2 := imgB.RGBA64At(x, y)

			r1, g1, b1, a1 := getRGBA64Components(c1)
			r2, g2, b2, a2 := getRGBA64Components(c2)

			// Exclusion blend mode formula: a + b - 2ab
			// For each channel, we calculate: bottom + top - (2 * bottom * top / max)
			r := r1 + r2 - ((r1 * r2) >> 15)
			g := g1 + g2 - ((g1 * g2) >> 15)
			b := b1 + b2 - ((b1 * b2) >> 15)
			a := porterDuffAlpha(a1, a2)

			setRGBA64Color(result, x, y, r, g, b, a)
		}
	}

	return result, nil
}
