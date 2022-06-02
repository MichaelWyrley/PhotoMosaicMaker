package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
)

func main() {
	img := flag.String("img", "cat.jpg", "The image you want to convert")
	number := flag.Int("no", 10, "The number of images you want to use")
	scaleString := flag.String("scale", "100x100", "The scale of the scraped images")
	imgShrink := flag.Int("shrink", 10, "How much the given image is scaled down by")
	location := flag.String("location", "./image", "The location of the final image")

	flag.Parse()

	scale := strings.Split(*scaleString, "x")
	scaleX, err := strconv.Atoi(scale[0])
	handleError(err, "converting width from string to int")
	scaleY, err := strconv.Atoi(scale[1])
	handleError(err, "converting height from string to int")

	fmt.Println("Getting Average Values")
	images := returnAverage(*number)
	fmt.Println("Gotten Average Values")

	fmt.Println("Generating Image")
	generateImages(*img, images, scaleX, scaleY, *imgShrink, *location)
	fmt.Println("Finished Generating Image")
}
