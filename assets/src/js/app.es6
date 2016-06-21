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

    var items = [];
    $('.cell').each(function() {
        var $size = $(this).data('size').split('x');
        items.push({
            src: $(this).find('a').attr('href'),
            w: $size[0],
            h: $size[1],
            title: $(this).data('caption')
        })
    });

    var $pswp = $('.pswp')[0];
    $('.gallery').on('click', '.cell', function(event) {
        event.preventDefault();

        var gallery;
        var $index = $(this).index();
        var options = {
            index: $index,
            bgOpacity: 0.7,
            showHideOpacity: true,
            getImageURLForShare: function() {
                var slideSrc = gallery.currItem.src;
                var filename = slideSrc.split('/')[1];
                return 'originals/' + filename;
            }
        }

        // Initialize PhotoSwipe
        gallery = new PhotoSwipe($pswp, PhotoSwipeUI_Default, items, options);
        gallery.init();
    });
});
