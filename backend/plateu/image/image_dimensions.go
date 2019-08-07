package image

import (
	"image"
	"log"
	"os"
)

func GetImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return 0, 0
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		log.Fatalf("%s", err.Error())
		return 0, 0
	}

	return image.Width, image.Height
}
