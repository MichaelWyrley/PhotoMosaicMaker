package main

import (
	"image"
	"image/jpeg"
	"os"
)

func squareDifference(r1, r2, g1, g2, b1, b2 uint64) uint64 {
	return (r1 - r2) ^ 2 + (g1 - g2) ^ 2 + (b1 - b2) ^ 2
}

func closest(r, g, b uint64, images []AverageImage) int {
	smallest := ^uint64(0)
	index := 0

	for i, img := range images {
		temp := squareDifference(r, img.average.r, g, img.average.g, b, img.average.b)
		if temp < smallest {
			smallest = temp
			index = i
		}
	}

	return index
}

func getValues(img image.Image, images []AverageImage, xMin, yMin, xMax, yMax int, ret chan [][]int) {

	var indexes [][]int

	// Calculate the closest value in images
	for i := yMin; i < yMax; i++ {
		var index []int
		for j := xMin; j < xMax; j++ {
			r, g, b, _ := img.At(j, i).RGBA()
			index = append(index, closest(uint64(r), uint64(g), uint64(b), images))
		}
		indexes = append(indexes, index)
	}

	ret <- indexes
}

// Calculates the image that corresponds to each colour in the given image
// returns an array of indexs that corresponds to each image in the average image array
func calculateImage(img image.Image, images []AverageImage) [][]int {
	step := img.Bounds().Max.Y / NO_SLICES
	var returns []chan [][]int

	// We move through the image in bands
	for i := img.Bounds().Min.Y; i <= img.Bounds().Max.Y; i += step {
		ret := make(chan [][]int, 100)
		go getValues(img, images, img.Bounds().Min.X, i, img.Bounds().Max.X, i+step, ret)
		returns = append(returns, ret)
	}

	var indexes [][]int

	// Get the go routines values
	for i := 0; i < NO_SLICES; i++ {
		indexes = append(indexes, <-returns[i]...)
	}

	return indexes
}

func createImage(indexes [][]int, images []AverageImage) {

}

func generateImages(img string, images []AverageImage) {
	// read in the image
	file, err := os.Open(img)
	handleError(err)

	imData, err := jpeg.Decode(file)
	handleError(err)

	indexs := calculateImage(imData, images)

	createImage(indexs, images)

}
