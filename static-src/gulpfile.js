var gulp = require('gulp');
var uncss = require('gulp-uncss');
var minifyCSS = require('gulp-minify-css');
var concat = require('gulp-concat');
var jade = require('gulp-jade');
var uglify = require('gulp-uglify');
var crypto = require('crypto');
var fs = require('fs');
var zopfli = require("gulp-zopfli");
var merge = require('merge-stream');
var util = require('gulp-util');

gulp.task('js', function() {
  return gulp.src(['js/*.js'])
  .pipe(concat('app.js'))
  .pipe(uglify())
  .pipe(gulp.dest('./www'));
});

gulp.task('css', function() {
    var base = gulp.src('css/*.css')
    .pipe(concat('app.css'))
    .pipe(uncss({
        ignore: [
          /* bootstrap UI */
          ".pull-right",
        ],
        html: ['www/index.html']
    }))
    .pipe(minifyCSS({'keepSpecialComments': 0}))
    .pipe(gulp.dest('www'));

    /*var extra = gulp.src(['assets/typicons/typicons.min.css', 'assets/animate.css', 'assets/extra.css'])
    .pipe(concat('extra.' + deployVersion + '.css'))
    .pipe(minifyCSS({'keepSpecialComments': 0}))
    .pipe(gulp.dest('www/assets/typicons'));

    return merge(base, extra);*/
    return base;
});

gulp.task('static', function() {
  var base = gulp.src([
    'index.html'
  ]).pipe(gulp.dest('www'));
  return base;
});

gulp.task('compress', ['static', 'js', 'css'], function() {
  return gulp.src(['www/**/*']).pipe(zopfli()).pipe(gulp.dest('www/'));
});

var tasks = ['static', 'js', 'css'];
if (util.env.deploy) {
  tasks.push('compress');
}
gulp.task('default', tasks);