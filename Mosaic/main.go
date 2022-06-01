package main

import (
	"flag"
)

func main() {
	img := flag.String("img", "cat.jpg", "The image you want to convert")
	number := flag.Int("no", 10, "The number of images you want to use")

	flag.Parse()

	images := returnAverage(*number)

	generateImages(*img, images)
}
