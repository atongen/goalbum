package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/disintegration/imaging"
	"github.com/namsral/flag"
)

// cli args
var (
	inFlag          = flag.String("in", "", "The input directory where images can be found")
	outFlag         = flag.String("out", "", "The output directory where the static gallery will be generated")
	maxThumbFlag    = flag.Int("max-thumb", 300, "Maximum pixel dimension of thumbnail images")
	maxSlideFlag    = flag.Int("max-slide", 1200, "Maximum pixel dimension of slide images")
	titleFlag       = flag.String("title", "", "Title of album")
	subtitleFlag    = flag.String("subtitle", "", "Subtitle of album")
	colorFlag       = flag.String("color", "blue", "CSS colors to use")
	headContentFlag = flag.String("head-content", "", "Path to file whose content should be included prior to the closing of the head element")
	bodyContentFlag = flag.String("body-content", "", "Path to file whose content should be included prior to the closing of the body element")
	includeFlag     strslice
	updateFlag      = flag.Bool("update", false, "If output directory is existing gallery, update instead of replace")
)

var (
	indexTmpl  *template.Template
	indexCtmpl = MustAsset("index.html.ctmpl")

	originalsDirName = "originals"
	slidesDirName    = "slides"
	thumbsDirName    = "thumbs"
	assetsDirName    = "assets"

	originalsDir string
	slidesDir    string
	thumbsDir    string
	assetsDir    string

	concurrency = runtime.NumCPU()
)

type strslice []string

func (s *strslice) String() string {
	return fmt.Sprintf("%s", *s)
}

func (s *strslice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type photoResult struct {
	photo *Photo
	err   error
}

type Page struct {
	Title        string
	Subtitle     string
	Photos       []*Photo
	CreatedAt    string
	Color        string
	HeadContent  string
	BodyContent  string
	Tags         map[string]string
	BuildVersion string
	BuildTime    string
	BuildHash    string
}

func init() {
	var err error
	indexTmpl, err = template.New("index").Parse(string(indexCtmpl))

	if err != nil {
		fmt.Printf("Invalid index template: %s\n", err.Error())
		os.Exit(1)
	}

	flag.Var(&includeFlag, "include", "File to include in document root of gallery")
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

	originalsDir = path.Join(*outFlag, originalsDirName)
	slidesDir = path.Join(*outFlag, slidesDirName)
	thumbsDir = path.Join(*outFlag, thumbsDirName)
	assetsDir = path.Join(*outFlag, assetsDirName)

	// attempt to parse existing photos.json file
	var existingPhotos []*Photo
	photoJsonPath := path.Join(*outFlag, "photos.json")
	photosBlob, err := ioutil.ReadFile(photoJsonPath)
	if os.IsNotExist(err) {
		existingPhotos = []*Photo{}
	} else {
		err = json.Unmarshal(photosBlob, &existingPhotos)
		if err != nil {
			fmt.Printf("Error parsing existing photo data: %s\n", err.Error())
			os.Exit(1)
		}
	}

	var photos []*Photo
	photos, err = IndexPhotos(*inFlag)

	if err != nil {
		fmt.Printf("Error indexing photos: %s\n", err.Error())
		os.Exit(1)
	}

	if len(existingPhotos) > 0 {
		PhotoUpdate(photos, existingPhotos)
	}

	photosToAdd := PhotoSliceSubtract(photos, existingPhotos)
	var photosToRm []*Photo
	if *updateFlag {
		// we are updating an existing gallery
		photos = PhotoUnion(photos, existingPhotos)
		photosToRm = []*Photo{}
	} else {
		// we are replacing existing gallery
		photos = PhotoRemoveDuplicates(photos)
		photosToRm = PhotoSliceSubtract(existingPhotos, photos)
	}

	sort.Sort(ByCreatedAt(photos))
	tags := PhotoTags(photos)
	SetTagNames(photos, tags)
	err = SetPhotoIds(photos)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	for _, photo := range photos {
		photo.SetDefaultCaption()
	}

	// create out directories
	for _, dir := range []string{originalsDir, slidesDir, thumbsDir, assetsDir} {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			fmt.Printf("Error creating image directory %s: %s\n", dir, err.Error())
			os.Exit(1)
		}
	}

	// remove obsolete photos in out directories
	for _, photo := range photosToRm {
		for _, dir := range []string{originalsDir, slidesDir, thumbsDir} {
			photoPath := path.Join(dir, photo.Filename())
			err = os.Remove(photoPath)
			if err != nil {
				fmt.Printf("Error removing old photo %s: %s\n", photoPath, err.Error())
				os.Exit(1)
			}
		}
	}

	ResizePhotos(photosToAdd)

	f, err := os.Create(path.Join(*outFlag, "index.html"))
	if err != nil {
		fmt.Printf("Error opening html: %s\n", err.Error())
		os.Exit(1)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	indexTmpl.Execute(w, Page{
		Title:        *titleFlag,
		Subtitle:     *subtitleFlag,
		Photos:       photos,
		CreatedAt:    time.Now().Format("Monday, January 2, 2006"),
		Color:        *colorFlag,
		HeadContent:  *headContentFlag,
		BodyContent:  *bodyContentFlag,
		Tags:         tags,
		BuildVersion: buildVersion,
		BuildTime:    buildTime,
		BuildHash:    buildHash,
	})
	w.Flush()

	for _, staticAsset := range staticAssets {
		err = writeStaticAsset(assetsDir, staticAsset)
		if err != nil {
			fmt.Printf("Error writing static asset %s: %s\n", staticAsset, err.Error())
			os.Exit(1)
		}
	}

	data, err := json.MarshalIndent(photos, "", "    ")
	if err != nil {
		fmt.Printf("Error converting photos json: %s\n", err.Error())
		os.Exit(1)
	}

	err = ioutil.WriteFile(photoJsonPath, data, 0644)
	if err != nil {
		fmt.Printf("Error writing photos json: %s\n", err.Error())
		os.Exit(1)
	}

	for _, includePath := range includeFlag {
		dst := path.Join(*outFlag, path.Base(includePath))
		err = CopyFile(dst, includePath)
		if err != nil {
			fmt.Printf("Error writing photos json: %s\n", err.Error())
			os.Exit(1)
		}
	}
}

func IndexPhotos(path string) ([]*Photo, error) {
	photos := []*Photo{}
	var wg sync.WaitGroup
	photoCh := make(chan photoResult)

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
		wg.Add(1)
		go func(photoPath string) {
			photo, err := IndexPhoto(photoPath)
			photoCh <- photoResult{photo, err}
		}(path)
		return nil
	})

	quit := make(chan bool)

	go func(inCh <-chan photoResult, doneCh <-chan bool) {
		for {
			select {
			case result := <-inCh:
				if result.err != nil {
					fmt.Printf("Error indexing photo: %s\n", err.Error())
					os.Exit(1)
				}
				photos = append(photos, result.photo)
				wg.Done()
			case <-doneCh:
				return
			}
		}
	}(photoCh, quit)

	wg.Wait()
	quit <- true
	close(photoCh)
	close(quit)

	return photos, err
}

func ResizePhotos(photos []*Photo) {
	photoCh := make(chan *Photo, concurrency)
	doneCh := make(chan bool)
	errCh := make(chan error)
	progCh := make(chan string)
	var wg sync.WaitGroup

	numPhotos := len(photos)

	// start the workers
	for i := 0; i < concurrency; i++ {
		go ResizeWorker(i, photoCh, doneCh, errCh, progCh, &wg)
	}

	quit := make(chan bool)

	// read from error channel
	go func(myErrCh <-chan error, myDoneCh <-chan bool) {
		for {
			select {
			case err := <-myErrCh:
				fmt.Printf("Error resizing photo: %s\n", err.Error())
			case <-myDoneCh:
				return
			}
		}
	}(errCh, quit)

	// read from progress channel
	go func(myProgCh <-chan string, myDoneCh <-chan bool) {
		var i int = 0
		for {
			select {
			case photo := <-myProgCh:
				i += 1
				fmt.Printf("%d / %d - %s\n", i, numPhotos, photo)
			case <-myDoneCh:
				return
			}
		}
	}(progCh, quit)

	// send photos into worker pool
	for _, photo := range photos {
		wg.Add(1)
		photoCh <- photo
	}

	wg.Wait()
	for i := 0; i < concurrency; i++ {
		doneCh <- true
	}
	// errCh
	quit <- true
	// progCh
	quit <- true
	close(photoCh)
	close(doneCh)
	close(errCh)
	close(progCh)
}

func ResizeWorker(id int, photoCh <-chan *Photo, doneCh <-chan bool, errCh chan<- error, progCh chan<- string, wg *sync.WaitGroup) {
	for {
		select {
		case photo := <-photoCh:
			err := ResizePhoto(photo)
			if err != nil {
				errCh <- err
			}
			progCh <- photo.Filename()
			wg.Done()
		case <-doneCh:
			return
		}
	}
}

func ResizePhoto(photo *Photo) error {
	fmt.Println(photo.Filename())
	file, err := os.Open(photo.InPath)
	if err != nil {
		return err
	}

	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}
	file.Close()

	// fix orientation
	orientation, err := GetOrientation(photo.InPath)
	if err == nil {
		FixOrientation(&img, orientation)
	}

	// write original image
	original, err := os.Create(path.Join(originalsDir, photo.Filename()))
	if err != nil {
		return err
	}

	err = jpeg.Encode(original, img, nil)
	if err != nil {
		return err
	}
	err = original.Close()
	if err != nil {
		return err
	}
	photo.OriginalWidth = img.Bounds().Dx()
	photo.OriginalHeight = img.Bounds().Dy()

	// write slide image
	slideImg := imaging.Fit(img, *maxSlideFlag, *maxSlideFlag, imaging.Lanczos)

	slide, err := os.Create(path.Join(slidesDir, photo.Filename()))
	if err != nil {
		return err
	}

	err = jpeg.Encode(slide, slideImg, nil)
	if err != nil {
		return err
	}
	err = slide.Close()
	if err != nil {
		return err
	}
	photo.SlideWidth = slideImg.Bounds().Dx()
	photo.SlideHeight = slideImg.Bounds().Dy()

	// write thumb image
	thumbImg := imaging.Fit(img, *maxThumbFlag, *maxThumbFlag, imaging.Lanczos)

	thumb, err := os.Create(path.Join(thumbsDir, photo.Filename()))
	if err != nil {
		return err
	}

	err = jpeg.Encode(thumb, thumbImg, nil)
	if err != nil {
		return err
	}
	err = thumb.Close()
	if err != nil {
		return err
	}
	photo.ThumbWidth = thumbImg.Bounds().Dx()
	photo.ThumbHeight = thumbImg.Bounds().Dy()

	return nil
}

func IndexPhoto(inPath string) (*Photo, error) {
	absPath, err := filepath.Abs(inPath)
	if err != nil {
		return nil, err
	}

	filename := path.Base(absPath)

	md5sum, err := Md5sumFromPath(absPath)
	if err != nil {
		return nil, err
	}

	return &Photo{
		InPath:       absPath,
		Md5sum:       md5sum,
		OriginalPath: path.Join(originalsDirName, filename),
		SlidePath:    path.Join(slidesDirName, filename),
		ThumbPath:    path.Join(thumbsDirName, filename),
		CreatedAt:    ImageTimeTaken(absPath),
	}, nil
}
