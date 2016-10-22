# goalbum

Build stand-alone html image galleries, from images on your hard drive.

This is basically a re-write of [rpa](https://github.com/atongen/rpa) in golang.

## Example

http://atongen.github.io/goalbum

## Installation

Download the latest release from the github releases page. Unzip and put
the goalbum binary somewhere on your PATH.

## Usage

Basic usage looks like this:

```shell
$ goalbum -in path/to/photo/directory -out path/to/html/output -title "My Cool Image Gallery"
```

### Tagging, Author, Caption

To add tags to photos or update the author or caption, first generate the gallery, for example:

```shell
$ goalbum -in path/to/photo/directory -out path/to/html/output -title "My Cool Image Gallery"
```

Then, edit the `photos.json` file in the output directory. Update each photo entry accordingly
by adding tags or other information. For example, this:

```json
...
{
		"Id": "photo-6",
		"InPath": "/home/atongen/tmp/example/022.jpg",
		"Md5sum": "61aa461810008e0bb50a62ae39c7c1ee",
		"OriginalPath": "originals/022.jpg",
		"OriginalWidth": 1600,
		"OriginalHeight": 1063,
		"SlidePath": "slides/022.jpg",
		"SlideWidth": 1200,
		"SlideHeight": 797,
		"ThumbPath": "thumbs/022.jpg",
		"ThumbWidth": 300,
		"ThumbHeight": 199,
		"Caption": "022.jpg: Monday, October 10, 2016 at 10:33am",
		"Author": "",
		"Tags": null,
		"TagNames": null,
		"CreatedAt": "2016-10-10T10:33:51-05:00"
},
...
```

could be updated to this:

```json
...
{
		"Id": "photo-6",
		"InPath": "/home/atongen/tmp/example/022.jpg",
		"Md5sum": "61aa461810008e0bb50a62ae39c7c1ee",
		"OriginalPath": "originals/022.jpg",
		"OriginalWidth": 1600,
		"OriginalHeight": 1063,
		"SlidePath": "slides/022.jpg",
		"SlideWidth": 1200,
		"SlideHeight": 797,
		"ThumbPath": "thumbs/022.jpg",
		"ThumbWidth": 300,
		"ThumbHeight": 199,
		"Caption": "Photo taken on Monday, October 10, 2016 at 10:33am",
		"Author": "Andrew Tongen",
		"Tags": ["Alice", "Bob"],
		"TagNames": null,
		"CreatedAt": "2016-10-10T10:33:51-05:00"
},
...
```

Then update the gallery:

```shell
$ goalbum -in path/to/photo/directory -out path/to/html/output -title "My Cool Image Gallery" -update
```

You can also quickly add and remove images from your gallery using this technique.
Keep your input directory around until your certain you like the way your gallery looks.

### Command Line Options

```shell
$ goalbum -h
Usage of goalbum:
  -body-content="": Path to file whose content should be included prior to the closing of the body element
  -color="blue": CSS colors to use (http://materializecss.com/color.html#palette)
  -exiftool="": Provide path to exiftool. If empty, PATH will be searched
  -head-content="": Path to file whose content should be included prior to the closing of the head element
  -in="": The input directory where images can be found
  -include=[]: File to include in document root of gallery
  -max-slide=1200: Maximum pixel dimension of slide images
  -max-thumb=300: Maximum pixel dimension of thumbnail images
  -out="": The output directory where the static gallery will be generated
  -subtitle="": Subtitle of album
  -title="": Title of album
  -update=false: If output directory is existing gallery, update instead of replace
  -version=false: Show the version and exit.
```

## Building

### Requirements

* golang 1.5.x or later
* [gb](https://getgb.io/)
* [npm](https://www.npmjs.com/)
* [grunt](http://gruntjs.com/)

```shell
$ git clone git@github.com/atongen/goalbum.git
$ cd goalbum
$ make
```

## Contributing

1. Fork it
1. Create your feature branch (`git checkout -b my-new-feature`)
1. Commit your changes (`git commit -am 'Add some feature'`)
1. Push to the branch (`git push origin my-new-feature`)
1. Create new Pull Request
