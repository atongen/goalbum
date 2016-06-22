package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"io/ioutil"
	"os"
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

func ImageTimeTaken(path string) time.Time {
	reader := exif.New()
	err := reader.Open(path)
	if err == nil {
		// DateTimeOriginal format 2011:06:04 08:56:22
		if val, ok := reader.Tags["DateTimeOriginal"]; ok {
			fmt.Println(val)
		}
	}
	f, err := os.Stat(path)
	if err != nil {
		return f.ModTime()
	}
	return time.Now()
}

func Md5sumFromPath(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(data)), nil
}

func SliceContainsString(s []string, b string) bool {
	for _, a := range s {
		if a == b {
			return true
		}
	}
	return false
}

func PhotoRemoveDuplicates(photos []*Photo) []*Photo {
	result := []*Photo{}
	md5s := []string{}

	for _, photo := range photos {
		if !SliceContainsString(md5s, photo.Md5sum) {
			result = append(result, photo)
			md5s = append(md5s, photo.Md5sum)
		}
	}

	return result
}

func PhotoSliceSubtract(photos1, photos2 []*Photo) []*Photo {
	result := []*Photo{}

CheckPhotos:
	for _, photo1 := range photos1 {
		for _, photo2 := range photos2 {
			if photo1.Md5sum == photo2.Md5sum {
				continue CheckPhotos
			}
		}
		result = append(result, photo1)
	}

	return result
}

// Prefer dup from photos1
// Ensure captions are preserved
func PhotoUnion(photos1, photos2 []*Photo) []*Photo {
	result := []*Photo{}
	md5s := []string{}
	captions := make(map[string]string)

	for _, photos := range [][]*Photo{photos1, photos2} {
		for _, photo := range photos {
			if photo.Caption != "" {
				if _, ok := captions[photo.Md5sum]; !ok {
					captions[photo.Md5sum] = photo.Caption
				}
			}
			if !SliceContainsString(md5s, photo.Md5sum) {
				result = append(result, photo)
				md5s = append(md5s, photo.Md5sum)
			}
		}
	}

	for _, photo := range result {
		if val, ok := captions[photo.Md5sum]; ok {
			if photo.Caption == "" {
				photo.Caption = val
			}
		}
	}

	return result
}
