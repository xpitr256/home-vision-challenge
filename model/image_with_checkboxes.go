package model

import (
	"image"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
)

type ImageWithBoxes struct {
	ImageUrl    string
	sourceImage *image.Gray
	checkboxes  []Checkbox
}

const imageFileName = "image_with_checkboxes.jpg"
const imageStorageDir = "./response"

func NewImageWithBoxes(sourceImage *image.Gray, boxes []Checkbox) (ImageWithBoxes, error) {
	imageWithBoxes := ImageWithBoxes{
		ImageUrl:    filepath.Join("/response", imageFileName),
		sourceImage: sourceImage,
		checkboxes:  boxes,
	}
	err := imageWithBoxes.generateImageWithCheckboxes()
	return imageWithBoxes, err
}

func (iwb *ImageWithBoxes) generateImageWithCheckboxes() error {
	coloredImage := iwb.createColoredImage()
	iwb.paintCheckboxes(coloredImage)
	// TODO: To avoid unnecessary cloud storage usage, this implementation overwrites
	// the response image on each request. In a real implementation, a UUID or another
	// unique identifier should be generated for each request to store the image with
	// a unique name per user or session.
	return iwb.saveImage(coloredImage)
}

func (iwb *ImageWithBoxes) createColoredImage() *image.RGBA {
	imageBounds := iwb.sourceImage.Bounds()
	coloredImage := image.NewRGBA(imageBounds)
	draw.Draw(coloredImage, imageBounds, iwb.sourceImage, imageBounds.Min, draw.Src)
	return coloredImage
}

func (iwb *ImageWithBoxes) paintCheckboxes(coloredImage *image.RGBA) {
	for _, checkbox := range iwb.checkboxes {
		checkboxColor := checkbox.getColor()
		// Add 3 px borders
		for x := checkbox.Box.Min.X; x < checkbox.Box.Max.X; x++ {
			// Top border
			coloredImage.Set(x, checkbox.Box.Min.Y, checkboxColor)
			coloredImage.Set(x, checkbox.Box.Min.Y+1, checkboxColor)
			coloredImage.Set(x, checkbox.Box.Min.Y+2, checkboxColor)
			// Bottom border
			coloredImage.Set(x, checkbox.Box.Max.Y-1, checkboxColor)
			coloredImage.Set(x, checkbox.Box.Max.Y-2, checkboxColor)
			coloredImage.Set(x, checkbox.Box.Max.Y-3, checkboxColor)
		}

		for y := checkbox.Box.Min.Y; y < checkbox.Box.Max.Y; y++ {
			// Left border
			coloredImage.Set(checkbox.Box.Min.X, y, checkboxColor)
			coloredImage.Set(checkbox.Box.Min.X+1, y, checkboxColor)
			coloredImage.Set(checkbox.Box.Min.X+2, y, checkboxColor)
			// Right border
			coloredImage.Set(checkbox.Box.Max.X-1, y, checkboxColor)
			coloredImage.Set(checkbox.Box.Max.X-2, y, checkboxColor)
			coloredImage.Set(checkbox.Box.Max.X-3, y, checkboxColor)
		}
	}
}

func (iwb *ImageWithBoxes) saveImage(coloredImage *image.RGBA) error {
	filePath := filepath.Join(imageStorageDir, imageFileName)
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	options := &jpeg.Options{Quality: 100}
	return jpeg.Encode(outputFile, coloredImage, options)
}
