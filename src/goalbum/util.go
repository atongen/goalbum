package main

import (
	"crypto/md5"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"os"
	"strconv"
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

func PhotoTags(photos []*Photo) map[string]string {
	tags := make(map[string]string)

	var i int = 0
	for _, photo := range photos {
		for _, tag := range photo.Tags {
			if _, ok := tags[tag]; !ok {
				tags[tag] = "tag-" + strconv.Itoa(i)
				i += 1
			}
		}
	}

	return tags
}

func SetTagNames(photos []*Photo, tags map[string]string) {
	for _, photo := range photos {
		if len(photo.Tags) > 0 {
			tagNames := make([]string, len(photo.Tags))
			for i, tag := range photo.Tags {
				tagNames[i] = tags[tag]
			}
			photo.TagNames = tagNames
		}
	}
}

func MapKeys(myMap map[string]string) []string {
	keys := make([]string, len(myMap))

	var i int = 0
	for key, _ := range myMap {
		keys[i] = key
		i += 1
	}

	return keys
}

func FindPhotoByMd5sum(photos []*Photo, md5sum string) *Photo {
	for _, photo := range photos {
		if photo.Md5sum == md5sum {
			return photo
		}
	}
	return nil
}

func SetPhotoIds(photos []*Photo) error {
	md5sums := []string{}
	ids := []string{}

	for _, photo := range photos {
		if SliceContainsString(md5sums, photo.Md5sum) {
			photo2 := FindPhotoByMd5sum(photos, photo.Md5sum)
			photo.Id = photo2.Id
			continue
		}
		for i := 1; i < 32; i++ {
			str := photo.Md5sum[0:i]
			if !SliceContainsString(ids, str) {
				photo.Id = fmt.Sprintf("photo-%s", str)
				ids = append(ids, str)
				break
			}
		}
		if photo.Id == "" {
			return fmt.Errorf("Unable to set id for photo %s\n", photo.Filename())
		}
	}

	return nil
}
