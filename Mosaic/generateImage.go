package main

import (
	"image"
	"image/jpeg"
	"os"
	"github.com/nfnt/resize"
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

// might be worth making this multithreaded using mutexes at some point
// Take the image and the location and places each colour in the final image
func setImage(img *image.RGBA, place image.Image, dim dimensions, y, x int){
	for i := 0; i < place.Bounds().Max.Y; i++ {
		for j := 0; j < place.Bounds().Max.X; j++ {
			img.Set(x*dim.scaleX + j, y*dim.scaleY + i, place.At(j, i))
		}
	}
}

func createImage(indexes [][]int, images []AverageImage, dim dimensions, location string) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{dim.width, dim.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// takes each index and places the image within the final image
	for i := 0; i < len(indexes); i++ {
		for j:= 0; j < len(indexes[0]); j++ {
			setImage(img, images[indexes[i][j]].image, dim, i, j)
		}
	}

	// export the final image
	f, _ := os.Create(location + ".jpg")
	err := jpeg.Encode(f, img, &jpeg.Options{QUALITY})
	handleError(err, "Error encoding final image")
}

func generateImages(img string, images []AverageImage, scaleX, scaleY int, imageShrink int, location string) {
	// read in the image
	file, err := os.Open(img)
	handleError(err, "Opening image")

	imData, err := jpeg.Decode(file)
	handleError(err, "decoding image")
	err = file.Close()
	handleError(err, "closing file")

	imData = resize.Resize(uint(imData.Bounds().Max.X/imageShrink), uint(imData.Bounds().Max.Y/imageShrink), imData, resize.Lanczos3)

	indexs := calculateImage(imData, images)
	width := imData.Bounds().Max.X * scaleX
	height := imData.Bounds().Max.Y * scaleY

	dim := dimensions{scaleX: scaleX, scaleY: scaleY, width: width, height: height}

	createImage(indexs, images, dim, location)

}
