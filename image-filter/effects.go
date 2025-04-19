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

// @Name: colorize
// @Desc: Colorizes the image
// @Param:      img     - - -   The image to colorize
// @Param:      col  	- - -   The color that determines the hue to use for colorization
// @Returns:    result  - - -	The colorized image
func colorize(img *image.NRGBA, col color.RGBA64) (*image.NRGBA, error) {
	bounds := img.Bounds()
	colorized := image.NewNRGBA(bounds)

	// Convert target color to normalized RGB and get alpha
	targetR := float64(col.R) / 65535.0
	targetG := float64(col.G) / 65535.0
	targetB := float64(col.B) / 65535.0
	alpha := float64(col.A) / 65535.0

	// Convert target color to HSL to get hue and saturation
	targetH, targetS, targetL := rgbToHsl(targetR, targetG, targetB)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.NRGBAAt(x, y)

			// Convert pixel to normalized RGB
			r := float64(c.R) / 255.0
			g := float64(c.G) / 255.0
			b := float64(c.B) / 255.0

			// Convert original pixel to HSL
			_, _, originalL := rgbToHsl(r, g, b)

			// Calculate new luminance by blending original and target luminance
			// This preserves the image's contrast while allowing some influence from target luminance
			newL := originalL*(1-alpha*0.5) + targetL*(alpha*0.5)

			// Calculate new saturation by blending original and target saturation
			// Original saturation is calculated from the RGB values
			originalS := calculateSaturation(r, g, b)
			newS := originalS*(1-alpha) + targetS*alpha

			// Convert back to RGB using the new HSL values
			newR, newG, newB := hslToRgb(targetH, newS, newL)

			// Blend with original color based on alpha
			finalR := r*(1-alpha) + newR*alpha
			finalG := g*(1-alpha) + newG*alpha
			finalB := b*(1-alpha) + newB*alpha

			// Ensure values are in valid range
			finalR = math.Min(math.Max(finalR, 0), 1)
			finalG = math.Min(math.Max(finalG, 0), 1)
			finalB = math.Min(math.Max(finalB, 0), 1)

			colorized.Set(x, y, color.NRGBA{
				R: uint8(finalR * 255),
				G: uint8(finalG * 255),
				B: uint8(finalB * 255),
				A: c.A,
			})
		}
	}
	return colorized, nil
}

// Helper function to calculate saturation from RGB
func calculateSaturation(r, g, b float64) float64 {
	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)

	if max == min {
		return 0
	}

	l := (max + min) / 2
	if l > 0.5 {
		return (max - min) / (2 - max - min)
	}
	return (max - min) / (max + min)
}

// Helper function to convert RGB to HSL
func rgbToHsl(r, g, b float64) (h, s, l float64) {
	max := math.Max(math.Max(r, g), b)
	min := math.Min(math.Min(r, g), b)

	l = (max + min) / 2

	if max == min {
		h = 0
		s = 0
	} else {
		d := max - min
		if l > 0.5 {
			s = d / (2 - max - min)
		} else {
			s = d / (max + min)
		}

		switch max {
		case r:
			h = (g - b) / d
			if g < b {
				h += 6
			}
		case g:
			h = (b-r)/d + 2
		case b:
			h = (r-g)/d + 4
		}
		h /= 6
	}

	return h, s, l
}

// Helper function to convert HSL to RGB
func hslToRgb(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		r, g, b = l, l, l
	} else {
		var q float64
		if l < 0.5 {
			q = l * (1 + s)
		} else {
			q = l + s - l*s
		}
		p := 2*l - q

		r = hueToRgb(p, q, h+1.0/3.0)
		g = hueToRgb(p, q, h)
		b = hueToRgb(p, q, h-1.0/3.0)
	}
	return r, g, b
}

// Helper function for HSL to RGB conversion
func hueToRgb(p, q, t float64) float64 {
	if t < 0 {
		t += 1
	}
	if t > 1 {
		t -= 1
	}
	if t < 1.0/6.0 {
		return p + (q-p)*6*t
	}
	if t < 1.0/2.0 {
		return q
	}
	if t < 2.0/3.0 {
		return p + (q-p)*(2.0/3.0-t)*6
	}
	return p
}
