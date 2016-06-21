package main

import (
	"image"
	"time"

	"github.com/disintegration/imaging"
	"github.com/xiam/exif"
)

// FixOrientation modifies image in-place to match exif orientation data
// http://sylvana.net/jpegcrop/exif_orientation.html
func FixOrientation(path string, img *image.Image) (string, error) {
	orientation := ""
	reader := exif.New()
	err := reader.Open(path)
	if err == nil {
		if orientation, ok := reader.Tags["Orientation"]; ok {
			switch orientation {
			case "Top-left":
				// 1
			case "Top-right":
				// 2 flop!
				*img = imaging.FlipH(*img)
			case "Bottom-right":
				// 3 rotate!(180)
				*img = imaging.Rotate180(*img)
			case "Bottom-left":
				// 4 flip!
				*img = imaging.FlipV(*img)
			case "Left-top":
				// 5 transpose!
				*img = imaging.Rotate270(*img)
				*img = imaging.FlipH(*img)
			case "Right-top":
				// 6 rotate!(90)
				*img = imaging.Rotate270(*img)
			case "Right-bottom":
				// 7 transverse!
				*img = imaging.Rotate90(*img)
				*img = imaging.FlipH(*img)
			case "Left-bottom":
				// 8 rotate!(270)
				*img = imaging.Rotate90(*img)
			}
		}
	}
	if orientation == "" {
		orientation = "Top-left"
	}

	return orientation, err
}

func ImageTimeTaken(path string) (time.Time, error) {
	reader := exif.New()
	err := reader.Open(path)
	if err == nil {
		// DateTimeOriginal format 2011:06:04 08:56:22
		if dateTimeOriginal, ok := reader.Tags["DateTimeOriginal"]; ok {
		}
	}
}
