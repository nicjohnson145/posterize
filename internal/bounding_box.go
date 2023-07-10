package internal

import (
	"image"
)

type Point struct {
	X int
	Y int
}

type BoundingBox struct {
	TopLeft     Point
	BottomRight Point
}

func (b *BoundingBox) ToRectangle() *image.RGBA {
	return image.NewRGBA(image.Rectangle{
		image.Point{X: b.TopLeft.X, Y: b.TopLeft.Y},
		image.Point{X: b.BottomRight.X, Y: b.BottomRight.Y},
	})
}
