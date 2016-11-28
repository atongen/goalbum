package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -prefix "templates/" templates/ templates/js/ templates/css/ templates/css/default-skin/ templates/fonts/ templates/fonts/roboto/

var (
	staticAssets = []string{
		"js/app.js",
		"css/app.css",
		"css/default-skin/default-skin.css",
		"css/default-skin/preloader.gif",
		"css/default-skin/default-skin.svg",
		"css/default-skin/default-skin.png",
		"fonts/roboto/Roboto-Medium.eot",
		"fonts/roboto/Roboto-Regular.woff",
		"fonts/roboto/Roboto-Regular.ttf",
		"fonts/roboto/Roboto-Medium.woff",
		"fonts/roboto/Roboto-Medium.ttf",
		"fonts/roboto/Roboto-Bold.eot",
		"fonts/roboto/Roboto-Bold.woff2",
		"fonts/roboto/Roboto-Light.ttf",
		"fonts/roboto/Roboto-Bold.woff",
		"fonts/roboto/Roboto-Thin.eot",
		"fonts/roboto/Roboto-Light.woff",
		"fonts/roboto/Roboto-Thin.ttf",
		"fonts/roboto/Roboto-Thin.woff2",
		"fonts/roboto/Roboto-Light.woff2",
		"fonts/roboto/Roboto-Regular.eot",
		"fonts/roboto/Roboto-Light.eot",
		"fonts/roboto/Roboto-Thin.woff",
		"fonts/roboto/Roboto-Regular.woff2",
		"fonts/roboto/Roboto-Bold.ttf",
		"fonts/roboto/Roboto-Medium.woff2",
	}
)

func writeStaticAsset(assetsDir, staticAsset string) error {
	absPath := path.Join(assetsDir, staticAsset)
	dir := filepath.Dir(absPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(absPath, MustAsset(staticAsset), 0644)
	if err != nil {
		return err
	}
	return nil
}
