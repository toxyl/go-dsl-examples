package main

import "image"

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
