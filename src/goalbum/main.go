package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/namsral/flag"
)

// cli args
var (
	inFlag          = flag.String("in", "", "The input directory where images can be found")
	outFlag         = flag.String("out", "", "The output directory where the static gallery will be generated")
	updateFlag      = flag.Bool("update", false, "Update an existing gallery")
	maxThumbFlag    = flag.Int("max-thumb", 300, "Maximum pixel dimension of thumbnail images")
	maxSlideFlag    = flag.Int("max-slide", 1200, "Maximum pixel dimension of slide images")
	headContentFlag = flag.String("head-content", "", "Path to file whose content should be included prior to the closing of the head element")
	bodyContentFlag = flag.String("body-content", "", "Path to file whose content should be included prior to the closing of the body element")
	includeFlag     = flag.String("include", "", "Comma separated list of files to include in document root of gallery")
)

var (
	photos []*Photo
)

type Photo struct {
	Id        string
	Path      string
	Width     int
	Height    int
	CreatedAt time.Time
}

func main() {
	flag.Parse()

	if *inFlag == "" {
		fmt.Println("in directory is required")
		os.Exit(0)
	}

	if *outFlag == "" {
		fmt.Println("out directory is required")
		os.Exit(0)
	}

	err := filepath.Walk(*inFlag, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if ext != ".jpg" {
			return nil
		}
		photo, err := newPhotoFromPath(path)
		if err != nil {
			return err
		}
		photos = append(photos, photo)
		return nil
	})

	if err != nil {
		fmt.Printf("Error getting image list: %s\n", err.Error())
		os.Exit(0)
	}

	if len(photos) == 0 {
		fmt.Println("No photos found")
		os.Exit(0)
	}

	err = os.MkdirAll(*outFlag, 0755)
	if err != nil {
		fmt.Printf("Error creating out directory: %s\n", err.Error())
		os.Exit(0)
	}

	originalsDir := path.Join(*outFlag, "originals")
	slidesDir := path.Join(*outFlag, "slides")
	thumbsDir := path.Join(*outFlag, "thumbs")

	for _, dir := range []string{originalsDir, slidesDir, thumbsDir} {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating image directory %s: %s\n", dir, err.Error())
			os.Exit(0)
		}
	}

	for _, photo := range photos {
		file, err := os.Open(photo.Path)
		if err != nil {
			fmt.Printf("Error opening image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}

		// decode jpeg into image.Image
		img, err := jpeg.Decode(file)
		if err != nil {
			fmt.Printf("Error decoding image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}
		file.Close()

		// fix orientation
		_, err = FixOrientation(photo.Path, &img)

		// write original image
		original, err := os.Create(path.Join(originalsDir, photo.Id+".jpg"))
		if err != nil {
			fmt.Printf("Error creating original image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}

		err = jpeg.Encode(original, img, nil)
		if err != nil {
			fmt.Printf("Error encoding original image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}
		err = original.Close()
		if err != nil {
			fmt.Printf("Error closing original image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}

		// write slide image
		slideImg := imaging.Fit(img, *maxSlideFlag, *maxSlideFlag, imaging.Lanczos)

		slide, err := os.Create(path.Join(slidesDir, photo.Id+".jpg"))
		if err != nil {
			fmt.Printf("Error creating slide image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}

		err = jpeg.Encode(slide, slideImg, nil)
		if err != nil {
			fmt.Printf("Error encoding slide image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}
		err = slide.Close()
		if err != nil {
			fmt.Printf("Error closing slide image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}

		// write thumb image
		thumbImg := imaging.Fit(img, *maxThumbFlag, *maxThumbFlag, imaging.Lanczos)

		thumb, err := os.Create(path.Join(thumbsDir, photo.Id+".jpg"))
		if err != nil {
			fmt.Printf("Error creating thumb image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}

		err = jpeg.Encode(thumb, thumbImg, nil)
		if err != nil {
			fmt.Printf("Error encoding thumb image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}
		err = thumb.Close()
		if err != nil {
			fmt.Printf("Error closing thumb image: %+v, %s\n", photo, err.Error())
			os.Exit(0)
		}
	}
}

func imageDims(path string) (int, int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, 0, err
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		return 0, 0, err
	}
	return image.Width, image.Height, nil
}

func newPhotoFromPath(inPath string) (*Photo, error) {
	absPath, err := filepath.Abs(inPath)
	if err != nil {
		return nil, err
	}

	w, h, err := imageDims(absPath)
	if err != nil {
		return nil, err
	}

	base := path.Base(inPath)
	ext := filepath.Ext(base)
	id := base[0 : len(base)-len(ext)]

	return &Photo{
		Id:        id,
		Path:      absPath,
		Width:     w,
		Height:    h,
		CreatedAt: time.Now(), // TODO: exif
	}, nil
}
