var indexOfObjAttr = function(arr, attr, value) {
    for (var i = 0; i < arr.length; i += 1) {
        if (arr[i][attr] === value) {
            return i;
        }
    }
}

var openGallery = function(id, selector) {
    var gallery;
    var items = parseItems(selector);
    var index = indexOfObjAttr(items, 'pid', id);
    var $pswp = $('.pswp')[0];

    var options = {
        index: index,
        bgOpacity: 0.7,
        showHideOpacity: true,
        history: true,
        galleryPIDs: true,
        getImageURLForShare: function() {
            var slideSrc = gallery.currItem.src;
            var filename = slideSrc.split('/')[1];
            return 'originals/' + filename;
        },
        addCaptionHTMLFn: function(item, captionEl, isFake) {
            if (!item.title) {
                captionEl.children[0].innerText = '';
                return false;
            }
            var caption = item.title;
            if (item.author) {
                caption.concat('<br/><small>Photo: ' + item.author + '</small>');
            }
            captionEl.children[0].innerHTML = caption;
            return true;
        },
        getThumbBoundsFn: function(index) {
            // See Options->getThumbBoundsFn section of docs for more info
            //var thumbnail = items[index].el.children[0],
            console.log(items[index].el);
            var thumbnail = items[index].el.children[0],
            pageYScroll = window.pageYOffset || document.documentElement.scrollTop,
            rect = thumbnail.getBoundingClientRect(); 

            return {x:rect.left, y:rect.top + pageYScroll, w:rect.width};
        },
    }

    // Initialize PhotoSwipe
    gallery = new PhotoSwipe($pswp, PhotoSwipeUI_Default, items, options);
    gallery.init();
}

var parseItems = function(selector) {
    var items = [];
    $(selector).each(function() {
        var $size = $(this).data('size').split('x');
        items.push({
            pid: $(this).data('photo-id'),
            src: $(this).find('a').attr('href'),
            msrc: $(this).data('msrc'),
            w: $size[0],
            h: $size[1],
            title: $(this).data('caption'),
            author: $(this).data('author'),
            el: $(this)[0]
        })
    });
	return items;
}

var photoswipeParseHash = function() {
	var hash = window.location.hash.substring(1),
	params = {};

	if(hash.length < 5) { // pid=1
		return params;
	}

	var vars = hash.split('&');
	for (var i = 0; i < vars.length; i++) {
		if(!vars[i]) {
			continue;
		}
		var pair = vars[i].split('=');  
		if(pair.length < 2) {
			continue;
		}
		params[pair[0]] = pair[1];
	}

	if(params.gid) {
		params.gid = parseInt(params.gid, 10);
	}

	return params;
};

var wall;
var cellSelector;
var loadWall = function() {
    wall = new Freewall(".gallery");
    wall.reset({
        animate: true,
        onResize: function() {
            wall.fitWidth();
        }
    })
    wall.fitWidth();
    $(window).trigger("resize");
}

$(function() {
    loadWall();

    cellSelector = '.cell';

    $('.tag-check').on('click', function(event) {
        var checkedTags = $('.tag-check:checked').map(function() {
            return $(this).val();
        }).get();

        if (checkedTags.length == 0) {
            // no tags checked, show all
            cellSelector = '.cell';
            wall.unFilter();
        } else {
            cellSelector = '.cell.' + checkedTags.join(', .cell.');
            wall.filter(cellSelector);
        }
    });

    $('.gallery').on('click', '.cell', function(event) {
        event.preventDefault();
        openGallery($(this).data('photo-id'), cellSelector);
    });

    // Parse URL and open gallery if it contains #&pid=3&gid=1
    var hashData = photoswipeParseHash();
    if(hashData.pid) {
        openGallery(hashData.pid, cellSelector);
    }
});
