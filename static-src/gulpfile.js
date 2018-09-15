var gulp = require('gulp');
var uncss = require('gulp-uncss');
var minifyCSS = require('gulp-minify-css');
var concat = require('gulp-concat');
var jade = require('gulp-jade');
var uglify = require('gulp-uglify');
var crypto = require('crypto');
var fs = require('fs');
//var zopfli = require("gulp-zopfli");
var merge = require('merge-stream');
var util = require('gulp-util');
var htmlmin = require('gulp-htmlmin');

gulp.task('js', function() {
  return gulp.src(['js/*.js'])
  .pipe(concat('app.js'))
  .pipe(uglify())
  .pipe(gulp.dest('../build'));
});

gulp.task('css', ['static'], function() {
    var base = gulp.src('css/*.css')
    .pipe(concat('app.css'))
    .pipe(uncss({
        ignore: [
          /* bootstrap UI */
          ".pull-right",
        ],
        html: ['../build/index.html']
    }))
    .pipe(minifyCSS({'keepSpecialComments': 0}))
    .pipe(gulp.dest('../build'));
    return base;
});

gulp.task('static', function() {
  var base = gulp.src([
    'index.html'
  ]).pipe(htmlmin({collapseWhitespace: true, removeComments: true}))
  .pipe(gulp.dest('../build'));

  var imgs = gulp.src([
    'images/**/*'
  ]).pipe(gulp.dest('../build/images'));
  return merge(base, imgs);
});

gulp.task('compress', ['static', 'js', 'css'], function() {
  return gulp.src(['../build/**/*']).pipe(zopfli()).pipe(gulp.dest('../build'));
});

var tasks = ['static', 'js', 'css'];
if (util.env.deploy) {
  //tasks.push('compress');
}
gulp.task('default', tasks);
