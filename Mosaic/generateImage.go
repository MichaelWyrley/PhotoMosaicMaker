package main

import (
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"os"
	"sync"
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

// Places the images in slices (called as a goroutine)
func slicedImage(img *image.RGBA, indexes [][]int, images []AverageImage, dim dimensions, startY, endY int, mutex *sync.Mutex, done chan bool) {
	// takes each index and places the image within the final image
	for i := startY; i < endY; i++ {
		for j:= 0; j < len(indexes[0]); j++ {
			mutex.Lock() // the mutexes are not really needed as the values won't overwrite each other but better safe than sorry
			setImage(img, images[indexes[i][j]].image, dim, i, j)
			mutex.Unlock()
		}
	}

	done <- true

}

func createImage(indexes [][]int, images []AverageImage, dim dimensions, location string) {
	upLeft := image.Point{0, 0}
	lowRight := image.Point{dim.width, dim.height}

	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	step := len(indexes) / NO_SLICES
	mutex := &sync.Mutex{}
	var returns []chan bool

	// Create the bands for the final image
	for i := 0; i < len(indexes); i += step {
		ret := make(chan bool)
		// technically there would be a bit of cut off from the fact that the step might not match up with the length of the final image
		// however this is such a small amount it doesn't really make a difference for the image
		go slicedImage(img, indexes, images, dim, i, i+step, mutex, ret)
		returns = append(returns, ret)
	}
	for _, e := range returns {
		<-e
	}

	// export the final image
	f, _ := os.Create(location + "out.jpg")
	err := jpeg.Encode(f, img, &jpeg.Options{QUALITY})
	handleError(err, "Error encoding final image")
}

func generateImages(img string, images []AverageImage, imageShrink int, location string) {
	// read in the image
	file, err := os.Open(img)
	handleError(err, "Opening image")

	imData, err := jpeg.Decode(file)
	handleError(err, "decoding image")
	err = file.Close()
	handleError(err, "closing file")

	imData = resize.Resize(uint(imData.Bounds().Max.X/imageShrink), uint(imData.Bounds().Max.Y/imageShrink), imData, resize.Lanczos3)

	indexs := calculateImage(imData, images)
	// get the width and height of the final image
	width := imData.Bounds().Max.X * images[0].image.Bounds().Max.X
	height := imData.Bounds().Max.Y * images[0].image.Bounds().Max.Y

	dim := dimensions{scaleX: images[0].image.Bounds().Max.X, scaleY: images[0].image.Bounds().Max.Y, width: width, height: height}

	createImage(indexs, images, dim, location)

}
