package main

import (
	"image"
	"image/color"
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

// hueToRGB helper function for HSL to RGB conversion
func hueToRGB(p, q, t float64) float64 {
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

// Helper function to convert HSL to RGB
func hslToRGB(h, s, l float64) (r, g, b float64) {
	if s == 0 {
		r, g, b = l, l, l
		return
	}

	var q float64
	if l < 0.5 {
		q = l * (1 + s)
	} else {
		q = l + s - l*s
	}
	p := 2*l - q

	r = hueToRGB(p, q, h+1.0/3.0)
	g = hueToRGB(p, q, h)
	b = hueToRGB(p, q, h-1.0/3.0)

	return
}
