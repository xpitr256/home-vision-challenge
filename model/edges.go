package model

import "image"

type Edges struct {
	Top    Edge
	Bottom Edge
	Left   Edge
	Right  Edge
}

func (e *Edges) IsStrong(x, y, size int, img *image.Gray) bool {
	return e.Top.IsStrong(x, y, size, img) && e.Bottom.IsStrong(x, y, size, img) && e.Left.IsStrong(x, y, size, img) && e.Right.IsStrong(x, y, size, img)
}
