package model

import "image"

// Edge is the interface that all edge types will implement
type Edge interface {
	IsStrong(x, y, size int, img *image.Gray) bool
}

type TopEdge struct{}

func (te *TopEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if x+i < img.Bounds().Max.X && img.GrayAt(x+i, y).Y == 0 {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

type BottomEdge struct{}

func (be *BottomEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if x+i < img.Bounds().Max.X && (img.GrayAt(x+i, y+size-1).Y == 0 || img.GrayAt(x+i, y+size-2).Y == 0) {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

type LeftEdge struct{}

func (le *LeftEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if y+i < img.Bounds().Max.Y && img.GrayAt(x, y+i).Y == 0 {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

type RightEdge struct{}

func (re *RightEdge) IsStrong(x, y, size int, img *image.Gray) bool {
	totalPixels, blackPixels := 0, 0
	for i := 0; i < size; i++ {
		if y+i < img.Bounds().Max.Y && (img.GrayAt(x+size-1, y+i).Y == 0 || img.GrayAt(x+size-2, y+i).Y == 0) {
			blackPixels++
		}
		totalPixels++
	}
	if totalPixels == 0 {
		return false
	}
	return isStrongBoundary(float64(blackPixels) / float64(totalPixels) * 100)
}

func isStrongBoundary(strength float64) bool {
	return strength > 82
}
