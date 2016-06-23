var openGallery = function(index) {
    var gallery;
    var items = parseItems();
    var $pswp = $('.pswp')[0];

    var options = {
        index: index,
        bgOpacity: 0.7,
        showHideOpacity: true,
        history: true,
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

var parseItems = function() {
    var items = [];
    $('.cell').each(function() {
        var $size = $(this).data('size').split('x');
        items.push({
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

$(function() {
    var wall = new Freewall(".gallery");
    wall.reset({
        animate: true,
        onResize: function() {
            wall.fitWidth();
        }
    })
    wall.fitWidth();
    // for scroll bar appear
    $(window).trigger("resize");

    $('.gallery').on('click', '.cell', function(event) {
        event.preventDefault();
        openGallery($(this).index());
    });

    // Parse URL and open gallery if it contains #&pid=3&gid=1
    var hashData = photoswipeParseHash();
    if(hashData.pid) {
        openGallery(hashData.pid);
    }
});
