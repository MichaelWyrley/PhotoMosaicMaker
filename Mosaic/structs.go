package main

import (
	"image"
	"log"
)

type ColourSum struct {
	r uint64
	g uint64
	b uint64
}

type AverageImage struct {
	image   image.Image
	average ColourSum
}

type dimensions struct {
	scaleX int
	scaleY int
	width int
	height int
}

const NO_SLICES = 10
const QUALITY = 50

func handleError(err error, info string) {
	if err != nil {
		log.Fatal("An error occurred ", info, " error: ", err)
	}
}
