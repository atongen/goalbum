package main

import (
	"bufio"
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/disintegration/imaging"
	"github.com/namsral/flag"
)

// cli args
var (
	inFlag       = flag.String("in", "", "The input directory where images can be found")
	outFlag      = flag.String("out", "", "The output directory where the static gallery will be generated")
	maxThumbFlag = flag.Int("max-thumb", 300, "Maximum pixel dimension of thumbnail images")
	maxSlideFlag = flag.Int("max-slide", 1200, "Maximum pixel dimension of slide images")
	titleFlag    = flag.String("title", "", "Title of album")
	subtitleFlag = flag.String("subtitle", "", "Subtitle of album")
	colorFlag    = flag.String("color", "blue", "CSS colors to use")
	// TODO: future version
	//updateFlag      = flag.Bool("update", false, "Update an existing gallery")
	//headContentFlag = flag.String("head-content", "", "Path to file whose content should be included prior to the closing of the head element")
	//bodyContentFlag = flag.String("body-content", "", "Path to file whose content should be included prior to the closing of the body element")
	//includeFlag     = flag.String("include", "", "Comma separated list of files to include in document root of gallery")
)

var (
	photos     []*Photo
	indexTmpl  *template.Template
	indexCtmpl = MustAsset("index.html.ctmpl")
)

type Photo struct {
	Id          string
	Path        string
	Width       int
	Height      int
	SlideWidth  int
	SlideHeight int
	ThumbWidth  int
	ThumbHeight int
	Caption     string
	CreatedAt   time.Time
}

type Page struct {
	Title     string
	Subtitle  string
	Photos    []*Photo
	CreatedAt string
	Color     string
}

func init() {
	// parse templates
	var err error

	indexTmpl, err = template.New("index").Parse(string(indexCtmpl))

	if err != nil {
		fmt.Printf("Invalid index template: %s\n", err.Error())
		os.Exit(1)
	}
}

func main() {
	flag.Parse()

	if *inFlag == "" {
		fmt.Println("in directory is required")
		os.Exit(1)
	}

	if *outFlag == "" {
		fmt.Println("out directory is required")
		os.Exit(1)
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
		os.Exit(1)
	}

	if len(photos) == 0 {
		fmt.Println("No photos found")
		os.Exit(1)
	}

	err = os.MkdirAll(*outFlag, 0755)
	if err != nil {
		fmt.Printf("Error creating out directory: %s\n", err.Error())
		os.Exit(1)
	}

	originalsDir := path.Join(*outFlag, "originals")
	slidesDir := path.Join(*outFlag, "slides")
	thumbsDir := path.Join(*outFlag, "thumbs")
	assetsDir := path.Join(*outFlag, "assets")

	for _, dir := range []string{originalsDir, slidesDir, thumbsDir, assetsDir} {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating image directory %s: %s\n", dir, err.Error())
			os.Exit(1)
		}
	}

	for _, photo := range photos {
		fmt.Println(photo.Path)
		file, err := os.Open(photo.Path)
		if err != nil {
			fmt.Printf("Error opening image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}

		// decode jpeg into image.Image
		img, err := jpeg.Decode(file)
		if err != nil {
			fmt.Printf("Error decoding image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}
		file.Close()

		// fix orientation
		_, err = FixOrientation(photo.Path, &img)

		// write original image
		original, err := os.Create(path.Join(originalsDir, photo.Id+".jpg"))
		if err != nil {
			fmt.Printf("Error creating original image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}

		err = jpeg.Encode(original, img, nil)
		if err != nil {
			fmt.Printf("Error encoding original image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}
		err = original.Close()
		if err != nil {
			fmt.Printf("Error closing original image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}

		// write slide image
		slideImg := imaging.Fit(img, *maxSlideFlag, *maxSlideFlag, imaging.Lanczos)

		slide, err := os.Create(path.Join(slidesDir, photo.Id+".jpg"))
		if err != nil {
			fmt.Printf("Error creating slide image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}

		err = jpeg.Encode(slide, slideImg, nil)
		if err != nil {
			fmt.Printf("Error encoding slide image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}
		err = slide.Close()
		if err != nil {
			fmt.Printf("Error closing slide image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}
		photo.SlideWidth = slideImg.Bounds().Dx()
		photo.SlideHeight = slideImg.Bounds().Dy()

		// write thumb image
		thumbImg := imaging.Fit(img, *maxThumbFlag, *maxThumbFlag, imaging.Lanczos)

		thumb, err := os.Create(path.Join(thumbsDir, photo.Id+".jpg"))
		if err != nil {
			fmt.Printf("Error creating thumb image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}

		err = jpeg.Encode(thumb, thumbImg, nil)
		if err != nil {
			fmt.Printf("Error encoding thumb image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}
		err = thumb.Close()
		if err != nil {
			fmt.Printf("Error closing thumb image: %+v, %s\n", photo, err.Error())
			os.Exit(1)
		}
		photo.ThumbWidth = thumbImg.Bounds().Dx()
		photo.ThumbHeight = thumbImg.Bounds().Dy()
	}

	f, err := os.Create(path.Join(*outFlag, "index.html"))
	if err != nil {
		fmt.Printf("Error opening html: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	indexTmpl.Execute(w, Page{
		Title:     *titleFlag,
		Subtitle:  *subtitleFlag,
		Photos:    photos,
		CreatedAt: time.Now().Format("Monday, January 2, 2006"),
		Color:     *colorFlag,
	})
	w.Flush()

	for _, staticAsset := range staticAssets {
		err = writeStaticAsset(assetsDir, staticAsset)
		if err != nil {
			fmt.Printf("Error writing static asset %s: %s\n", staticAsset, err.Error())
			os.Exit(1)
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
