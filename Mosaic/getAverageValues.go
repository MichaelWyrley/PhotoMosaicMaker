package main

import (
	"image"
	"image/jpeg"
	"os"
	"strconv"
)

// Sum a slice of the image
func sumValues(img image.Image, xMin, yMin, xMax, yMax int, ret chan ColourSum) {
	sum := ColourSum{0, 0, 0}
	no_values := uint64((xMax - xMin) * (yMax - yMin))

	// Calculate the sum of all the colours for the slice
	for i := yMin; i < yMax; i++ {
		for j := xMin; j < xMax; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			sum.r += uint64(r)
			sum.g += uint64(g)
			sum.b += uint64(b)
		}
	}

	// Calculate the average colour of the slice
	sum.r /= no_values
	sum.g /= no_values
	sum.b /= no_values
	ret <- sum

}

// Calculate the average colour for the image
func getAverageValues(img image.Image, r chan AverageImage) {

	step := img.Bounds().Max.Y / NO_SLICES
	var returns []chan ColourSum

	// Making the bands for the image
	for i := img.Bounds().Min.Y; i <= img.Bounds().Max.Y; i += step {
		ret := make(chan ColourSum)
		go sumValues(img, img.Bounds().Min.X, i, img.Bounds().Max.X, i+step, ret)
		returns = append(returns, ret)
	}

	// collect the average for each slice
	sum := ColourSum{0, 0, 0}
	for i := 0; i < NO_SLICES; i++ {
		v := <-returns[i]
		sum.r += v.r
		sum.b += v.b
		sum.g += v.g
	}
	// calculate the average for the image
	sum.r /= NO_SLICES
	sum.g /= NO_SLICES
	sum.b /= NO_SLICES

	r <- AverageImage{image: img, average: sum}
}

func returnAverage(no_images int, imageLocation string) []AverageImage {

	images := make([]AverageImage, no_images)
	averageReturns := make([]chan AverageImage, no_images)

	// Create go routines that calculate the average colour value of each image
	for i := 0; i < no_images; i++ {
		file, err := os.Open(imageLocation + strconv.Itoa(i) + ".jpg")
		handleError(err, "reading image from file")

		imData, err := jpeg.Decode(file)
		handleError(err, "decoding image")
		err = file.Close()
		handleError(err, "closing file")

		r := make(chan AverageImage)
		go getAverageValues(imData, r)
		averageReturns[i] = r
	}

	// Collect the average pixle values
	for i := 0; i < no_images; i++ {
		images[i] = <-averageReturns[i]
	}

	return images
}
