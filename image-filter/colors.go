package main

import "image/color"

// @Name: hsla
// @Desc: Creates a color from HSLA values
// @Param:      h      	"°" 0..360   	0.0   	The color's hue
// @Param:      s     	"%" 0.0..1.0   	0.5   	The color's saturation
// @Param:      l     	"%" 0.0..1.0   	0.5   	The color's luminosity
// @Param:      alpha  	"%" 0.0..1.0   	1.0   	The color's alpha
// @Returns:    result  - 	-   		-   	The color as color.RGBA64
func hsla(h float64, s float64, l float64, alpha float64) (color.RGBA64, error) {
	// Convert HSLA to RGB
	var r, g, b float64
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

		r = hueToRGB(p, q, h/360+1.0/3.0)
		g = hueToRGB(p, q, h/360)
		b = hueToRGB(p, q, h/360-1.0/3.0)
	}

	// Convert to 16-bit color channels
	R := uint16(r * 65535)
	G := uint16(g * 65535)
	B := uint16(b * 65535)
	A := uint16(alpha * 65535)

	return color.RGBA64{
		R: R,
		G: G,
		B: B,
		A: A,
	}, nil
}
