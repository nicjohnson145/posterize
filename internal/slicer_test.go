package internal

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetBoundingBoxes(t *testing.T) {
	const maxLongSide = 10.5
	const maxShortSide = 8.0
	const ppi = 100.0

	newSlicer := func() *Slicer {
		return NewSlicer(SlicerConfig{
			PixelsPerInch:    ppi,
			MaxPageLongSide:  maxLongSide,
			MaxPageShortSide: maxShortSide,
		})
	}

	t.Run("1 page landscape smokes", func(t *testing.T) {
		slicer := newSlicer()

		boxes := slicer.getBoundingBoxes(1050.0, 800)
		require.Equal(t, 1, len(boxes))
	})

	t.Run("1 page portrait smokes", func(t *testing.T) {
		slicer := newSlicer()

		boxes := slicer.getBoundingBoxes(800, 1050.0)
		require.Equal(t, 1, len(boxes))
	})

	t.Run("semi real landscape", func(t *testing.T) {
		slicer := newSlicer()

		boxes := slicer.getBoundingBoxes(2000, 1600)
		require.Equal(t, 4, len(boxes))
		require.Equal(
			t,
			[]BoundingBox{
				{TopLeft: Point{X: 0, Y: 0}, BottomRight: Point{X: 1050, Y: 800}},
				{TopLeft: Point{X: 1050, Y: 0}, BottomRight: Point{X: 2000, Y: 800}},
				{TopLeft: Point{X: 0, Y: 800}, BottomRight: Point{X: 1050, Y: 1600}},
				{TopLeft: Point{X: 1050, Y: 800}, BottomRight: Point{X: 2000, Y: 1600}},
			},
			boxes,
		)
	})
}
