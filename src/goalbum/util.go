package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/disintegration/imaging"
	"github.com/rwcarlsen/goexif/exif"
)

func GetOrientation(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}

	x, err := exif.Decode(f)
	if err != nil {
		return 0, err
	}

	orientation, err := x.Get(exif.Orientation)
	if err != nil {
		return 0, err
	}

	val, err := orientation.Int(0)
	if err != nil {
		return 0, err
	}

	return val, nil
}

// FixOrientation modifies image in-place to match exif orientation data
// http://sylvana.net/jpegcrop/exif_orientation.html
func FixOrientation(img *image.Image, orientation int) error {
	switch orientation {
	case 1:
		// do nothing
	case 2:
		// flop!
		*img = imaging.FlipH(*img)
	case 3:
		// rotate!(180)
		*img = imaging.Rotate180(*img)
	case 4:
		// flip!
		*img = imaging.FlipV(*img)
	case 5:
		// transpose!
		*img = imaging.Rotate270(*img)
		*img = imaging.FlipH(*img)
	case 6:
		// rotate!(90)
		*img = imaging.Rotate270(*img)
	case 7:
		// transverse!
		*img = imaging.Rotate90(*img)
		*img = imaging.FlipH(*img)
	case 8:
		// rotate!(270)
		*img = imaging.Rotate90(*img)
	default:
		return fmt.Errorf("Invalid orientation %d", orientation)
	}
	return nil
}

func ImageTimeTaken(path string) time.Time {
	f, err := os.Open(path)
	if err != nil {
		return time.Now()
	}

	var timeTaken time.Time

	x, err := exif.Decode(f)
	if err == nil {
		timeTaken, _ = x.DateTime()
	}

	if timeTaken.IsZero() {
		fi, err := f.Stat()
		if err == nil {
			timeTaken = fi.ModTime()
		}
	}

	if timeTaken.IsZero() {
		timeTaken = time.Now()
	}

	return timeTaken
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

func PhotoUpdate(photos1, photos2 []*Photo) {
	for _, photo1 := range photos1 {
		for _, photo2 := range photos2 {
			if photo1.Md5sum == photo2.Md5sum {
				photo1.Update(photo2)
			}
		}
	}
}

func PhotoUnion(photos1, photos2 []*Photo) []*Photo {
	result := []*Photo{}
	md5s := []string{}

	for _, photos := range [][]*Photo{photos1, photos2} {
		for _, photo := range photos {
			if !SliceContainsString(md5s, photo.Md5sum) {
				result = append(result, photo)
				md5s = append(md5s, photo.Md5sum)
			}
		}
	}

	return result
}

func CopyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err != nil {
		return err
	}
	return cerr
}
