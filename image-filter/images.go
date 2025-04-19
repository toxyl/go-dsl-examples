package main

import (
	"image"
	"image/png"
	"os"
)

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
