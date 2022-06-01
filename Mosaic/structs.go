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

const NO_SLICES = 10

func handleError(err error) {
	if err != nil {
		log.Fatal("An error occurred ", err)
	}
}
