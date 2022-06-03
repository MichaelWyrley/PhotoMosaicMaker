package main

import (
	"flag"
	"fmt"
)

func main() {
	img := flag.String("img", "cat.jpg", "The image you want to convert")
	number := flag.Int("no", 10, "The number of images you want to use")
	imgShrink := flag.Int("shrink", 10, "How much the given image is scaled down by")
	location := flag.String("location", "./image", "The location of the final image")

	flag.Parse()

	fmt.Println("Getting Average Values")
	images := returnAverage(*number)
	fmt.Println("Gotten Average Values")

	fmt.Println("Generating Image")
	generateImages(*img, images, *imgShrink, *location)
	fmt.Println("Finished Generating Image")
}
