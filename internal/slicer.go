package internal

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math"

	"github.com/rs/zerolog"
)

type SlicerConfig struct {
	Logger           zerolog.Logger
	PixelsPerInch    float64
	MaxPageLongSide  float64
	MaxPageShortSide float64
}

func NewSlicer(conf SlicerConfig) *Slicer {
	return &Slicer{
		log:              conf.Logger,
		pixelsPerInch:    conf.PixelsPerInch,
		maxPageLongSide:  conf.MaxPageLongSide,
		maxPageShortSide: conf.MaxPageShortSide,
	}
}

type Slicer struct {
	log              zerolog.Logger
	pixelsPerInch    float64
	maxPageLongSide  float64
	maxPageShortSide float64
}

func (s *Slicer) ImageFromBytes(input []byte) (image.Image, error) {
	s.log.Debug().Msg("decoding image to determine type")
	_, imgType, err := image.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, fmt.Errorf("error decoding image: %w", err)
	}

	s.log.Debug().Msgf("decoded as %v", imgType)
	switch imgType {
	case "png":
		return png.Decode(bytes.NewReader(input))
	default:
		return nil, fmt.Errorf("unable to process image type of %v", imgType)
	}
}

func (s *Slicer) getBoundingBoxes(width float64, height float64) []BoundingBox {
	// Determine if we should slice in portrait or landscape
	s.log.Debug().Float64("height", height).Float64("width", width).Msg("sizing image")
	var maxPixelHigh, maxPixelWide float64
	if height > width {
		s.log.Debug().Msg("slicing as portrait")
		maxPixelHigh = s.maxPageLongSide * s.pixelsPerInch
		maxPixelWide = s.maxPageShortSide * s.pixelsPerInch
	} else {
		s.log.Debug().Msg("slicing as landscape")
		maxPixelHigh = s.maxPageShortSide * s.pixelsPerInch
		maxPixelWide = s.maxPageLongSide * s.pixelsPerInch
	}

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
