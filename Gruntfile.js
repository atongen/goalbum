module.exports = function(grunt) {
    require('load-grunt-tasks')(grunt);

    grunt.config('env', grunt.option('env') || process.env.GRUNT_ENV || 'development');

    var production = grunt.config('env') === 'production';

    grunt.loadNpmTasks('grunt-contrib-copy');
    grunt.loadNpmTasks('grunt-contrib-uglify');
    grunt.loadNpmTasks('grunt-sass');
    grunt.loadNpmTasks('grunt-contrib-cssmin');
    grunt.loadNpmTasks('grunt-contrib-clean')

    grunt.initConfig({
        clean: {
            main: {
                src: "public/assets/*"
            }
        },
        copy: {
            main: {
                files: [
                    // photoswipe
                    {src:'node_modules/photoswipe/dist/photoswipe.js', dest: 'assets/build/js/photoswipe.js'},
                    {src:'node_modules/photoswipe/dist/photoswipe-ui-default.js', dest: 'assets/build/js/photoswipe-ui-default.js'},
                    {src:'node_modules/photoswipe/dist/photoswipe.css', dest: 'assets/build/css/photoswipe.css'},
                    // photoswipe default skin
                    {src:'node_modules/photoswipe/dist/default-skin/default-skin.css', dest: 'src/goalbum/templates/css/default-skin/default-skin.css'},
                    {src:'node_modules/photoswipe/dist/default-skin/default-skin.png', dest: 'src/goalbum/templates/css/default-skin/default-skin.png'},
                    {src:'node_modules/photoswipe/dist/default-skin/default-skin.svg', dest: 'src/goalbum/templates/css/default-skin/default-skin.svg'},
                    {src:'node_modules/photoswipe/dist/default-skin/preloader.gif', dest: 'src/goalbum/templates/css/default-skin/preloader.gif'},
                    // freewall
                    {src: 'node_modules/freewall/freewall.js', dest: 'assets/build/js/freewall.js'},
                    // materialize
                    {src: 'node_modules/materialize-css/dist/css/materialize.css', dest: 'assets/build/css/materialize.css'},
                    {
                        expand: true,
                        cwd: 'node_modules/materialize-css/dist/font/',
                        src: '**/*',
                        dest: 'src/goalbum/templates/font'
                    },
                    {
                        expand: true,
                        cwd: 'node_modules/materialize-css/dist/fonts/',
                        src: '**/*',
                        dest: 'src/goalbum/templates/fonts'
                    },
                    // jquery
                    {src: 'node_modules/jquery/dist/jquery.js', dest: 'assets/build/js/jquery.js'}

                ]
            }
        },
        babel: {
            options: {
                sourceMap: true,
                presets: ['es2015']
            },
            main: {
                files: [
                    {
                        expand: true,
                        cwd: 'assets/src/js/',
                        src: '**/*.es6',
                        dest: 'assets/build/js/',
                        ext: '.js'
                    }
                ]
            }
        },
        uglify: {
            options: {
                mangle: production,
                compress: production,
                beautify: !production,
                sourceMap: !production,
                preserveComments: false,
                screwIE8: true
            },
            main: {
                files: {
                    'src/goalbum/templates/js/app.js': [
                        'assets/build/js/jquery.js',
                        'assets/build/js/photoswipe.js',
                        'assets/build/js/photoswipe-ui-default.js',
                        'assets/build/js/freewall.js',
                        'assets/build/js/app.js',
                        'assets/build/js/**/*.js'
                    ]
                }
            }
        },
        sass: {
            options: {},
            main: {
                files: [
                    {
                        expand: true,
                        cwd: 'assets/src/css/',
                        src: '**/*.scss',
                        dest: 'assets/build/css/',
                        ext: '.css'
                    }
                ]
            }
        },
        cssmin: {
            options: {
                shorthandCompacting: false,
                roundingPrecision: -1,
                keepSpecialComments: production ? 0 : 1
            },
            main: {
                files: {
                    'src/goalbum/templates/css/app.css': [
                        'assets/build/css/photoswipe.css',
                        'assets/build/css/default-skin.css',
                        'assets/build/css/materialize.css',
                        'assets/build/css/**/*.css'
                    ]
                }
            }
        }
    });

    grunt.registerTask('default', ['clean','copy','babel','uglify','sass','cssmin']);
}
