package main

import (
	"fmt"
	"path"
	"time"
)

type Photo struct {
	InPath         string
	Md5sum         string
	OriginalPath   string
	OriginalWidth  int
	OriginalHeight int
	SlidePath      string
	SlideWidth     int
	SlideHeight    int
	ThumbPath      string
	ThumbWidth     int
	ThumbHeight    int
	Caption        string
	Author         string
	CreatedAt      time.Time
}

func (photo *Photo) Filename() string {
	return path.Base(photo.InPath)
}

func (photo1 *Photo) Update(photo2 *Photo) {
	if photo1.OriginalWidth == 0 && photo2.OriginalWidth != 0 {
		photo1.OriginalWidth = photo2.OriginalWidth
	}
	if photo1.OriginalHeight == 0 && photo2.OriginalHeight != 0 {
		photo1.OriginalHeight = photo2.OriginalHeight
	}
	if photo1.SlideWidth == 0 && photo2.SlideWidth != 0 {
		photo1.SlideWidth = photo2.SlideWidth
	}
	if photo1.SlideHeight == 0 && photo2.SlideHeight != 0 {
		photo1.SlideHeight = photo2.SlideHeight
	}
	if photo1.ThumbWidth == 0 && photo2.ThumbWidth != 0 {
		photo1.ThumbWidth = photo2.ThumbWidth
	}
	if photo1.ThumbHeight == 0 && photo2.ThumbHeight != 0 {
		photo1.ThumbHeight = photo2.ThumbHeight
	}
	if photo1.Caption == "" && photo2.Caption != "" {
		photo1.Caption = photo2.Caption
	}
	if photo1.Author == "" && photo2.Author != "" {
		photo1.Author = photo2.Author
	}
}

func (photo *Photo) SetDefaultCaption() {
	if photo.Caption == "" {
		photo.Caption = fmt.Sprintf("%s: %s", photo.Filename(), photo.CreatedAt.Format("Monday, January 2, 2006 at 3:04pm"))
	}
}

type ByCreatedAt []*Photo

func (p ByCreatedAt) Len() int {
	return len(p)
}

func (p ByCreatedAt) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p ByCreatedAt) Less(i, j int) bool {
	return p[i].CreatedAt.Before(p[j].CreatedAt)
}
