package main

import (
	"os"
	"os/exec"
	"path/filepath"
)

var (
	exiftoolName = "exiftool"
)

func ExiftoolPath(toolPath string) (myPath string, err error) {
	if toolPath == "" {
		// provided path is empty, search PATH, it is not an error
		// if exiftool isn't found
		myPath, _ = ExiftoolPathFind()
	} else {
		myPath, err = ExiftoolPathValidate(toolPath)
	}
	return
}

// ExiftoolPathValidate validates the provided path to exiftool exists,
// and if so, returns the absolute path to it.
func ExiftoolPathValidate(toolPath string) (myPath string, err error) {
	_, err = os.Stat(toolPath)
	if err != nil {
		return
	}
	myPath, err = filepath.Abs(toolPath)
	return
}

// ExiftoolPathFind searches PATH for exiftool and returns the absolute
// path to it if found.
func ExiftoolPathFind() (myPath string, err error) {
	myPath, err = exec.LookPath(exiftoolName)
	if err != nil {
		return
	}
	myPath, err = filepath.Abs(myPath)
	return
}

// ExifCp copies exif data from src image to dst image.
// It excludes Orientation tag because the orientation has been
// normalized in the processed images.
func ExifCp(toolPath, src, dst string) (string, error) {
	cmd := exec.Command(toolPath, "-overwrite_original_in_place", "-tagsFromFile", src, "-x", "Orientation", dst)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
