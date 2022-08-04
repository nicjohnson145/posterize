package main
/*
Total Image; 2000W x 1600H
*/

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"image"
	"image/png"
	"math"
	"os"
)

var (
	ppi       int
	landscape bool
	prefix    string
)

const (
	//maxHeightInches = 10.5
	//maxWidthInches = 8.0

	// Assume landscape
	maxHeightInches float64 = 8.0
	maxWidthInches  float64 = 10.5
)

type Point struct {
	X int
	Y int
}

type BoundingBox struct {
	TopLeft Point
	//TopRight    Point
	//BottomLeft  Point
	BottomRight Point
}

func main() {
	if err := buildCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}

func buildCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "posterize",
		Short: "Split images",
		Args:  cobra.ExactArgs(1),
		Long:  "Split images to be printed across multiple sheets",
		Run: func(cmd *cobra.Command, args []string) {
			run(args[0])
		},
	}
	rootCmd.Flags().IntVar(&ppi, "ppi", 100, "Pixels per inch")
	rootCmd.Flags().StringVar(&prefix, "prefix", "img", "Prefix to place in front of split images")

	return rootCmd
}

func run(path string) {
	img, err := openImage(path)
	if err != nil {
		log.Fatal(err)
	}
	boxes := getBoundingBoxes(float64(img.Bounds().Dx()), float64(img.Bounds().Dy()))
	parts := sliceImage(img, boxes)
	if err := writeImages(parts); err != nil {
		log.Fatal(err)
	}
}

func openImage(path string) (image.Image, error) {
	origImage, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer origImage.Close()

	_, origType, err := image.Decode(origImage)
	if err != nil {
		return nil, fmt.Errorf("error determining image type: %w", err)
	}
	if origType != "png" {
		return nil, fmt.Errorf("image is not a png")
	}
	origImage.Seek(0, 0)

	return png.Decode(origImage)
}

func getBoundingBoxes(width float64, height float64) []BoundingBox {
	fPpi := float64(ppi)

	maxPixelHigh := maxHeightInches * fPpi
	maxPixelWide := maxWidthInches * fPpi

	numPagesWide := math.Ceil(width / maxPixelWide)
	numPagesHigh := math.Ceil(height / maxPixelHigh)

	boxes := []BoundingBox{}

	var x, y float64

	for y = 0; y < numPagesHigh; y++ {
		for x = 0; x < numPagesWide; x++ {
			boxes = append(boxes, BoundingBox{
				TopLeft: Point{
					X: int(x * maxPixelWide),
					Y: int(y * maxPixelHigh),
				},
				BottomRight: Point{
					X: int(math.Min(
						width,
						(x + 1) * maxPixelWide,
					)),
					Y: int(math.Min(
						height,
						(y + 1) * maxPixelHigh,
					)),
				},
			})
		}
	}

	return boxes
}

func sliceImage(img image.Image, boundingBoxes []BoundingBox) []image.Image {
	newImages := []image.Image{}

	for _, bb := range boundingBoxes {
		newImg := image.NewRGBA(image.Rectangle{
			image.Point{X: bb.TopLeft.X, Y: bb.TopLeft.Y},
			image.Point{X: bb.BottomRight.X, Y: bb.BottomRight.Y},
		})

		for x := bb.TopLeft.X; x <= bb.BottomRight.X; x++ {
			for y := bb.TopLeft.Y; y <= bb.BottomRight.Y; y++ {
				newImg.Set(x, y, img.At(x, y))
			}
		}

		newImages = append(newImages, newImg)
	}

	return newImages
}

func writeImages(images []image.Image) error {
	for i, img := range images {
		name := fmt.Sprintf("%v-%v.png", prefix, i)
		f, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("error creating %v: %w", name, err)
		}
		defer f.Close()

		err = png.Encode(f, img)
		if err != nil {
			return fmt.Errorf("error encoding %v as png: %w", name, err)
		}
	}

	return nil
}
