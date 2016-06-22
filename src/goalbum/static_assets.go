package main

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

//go:generate go-bindata -pkg $GOPACKAGE -o assets.go -prefix "templates/" templates/ templates/js/ templates/css/ templates/css/default-skin/ templates/font/ templates/fonts/ templates/font/material-design-icons/ templates/font/roboto/ templates/fonts/roboto/

var (
	staticAssets = []string{
		"js/app.js",
		"css/app.css",
		"css/default-skin/default-skin.css",
		"css/default-skin/preloader.gif",
		"css/default-skin/default-skin.svg",
		"css/default-skin/default-skin.png",
		"font/material-design-icons/Material-Design-Icons.eot",
		"font/material-design-icons/Material-Design-Icons.ttf",
		"font/material-design-icons/Material-Design-Icons.woff",
		"font/material-design-icons/Material-Design-Icons.woff2",
		"font/material-design-icons/Material-Design-Icons.svg",
		"font/roboto/Roboto-Medium.eot",
		"font/roboto/Roboto-Regular.woff",
		"font/roboto/Roboto-Regular.ttf",
		"font/roboto/Roboto-Medium.woff",
		"font/roboto/Roboto-Medium.ttf",
		"font/roboto/Roboto-Bold.eot",
		"font/roboto/Roboto-Bold.woff2",
		"font/roboto/Roboto-Light.ttf",
		"font/roboto/Roboto-Bold.woff",
		"font/roboto/Roboto-Thin.eot",
		"font/roboto/Roboto-Light.woff",
		"font/roboto/Roboto-Thin.ttf",
		"font/roboto/Roboto-Thin.woff2",
		"font/roboto/Roboto-Light.woff2",
		"font/roboto/Roboto-Regular.eot",
		"font/roboto/Roboto-Light.eot",
		"font/roboto/Roboto-Thin.woff",
		"font/roboto/Roboto-Regular.woff2",
		"font/roboto/Roboto-Bold.ttf",
		"font/roboto/Roboto-Medium.woff2",
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
