const { series, parallel, src, dest } = require('gulp');

var uncss = require('gulp-uncss');
var minifyCSS = require('gulp-minify-css');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
var crypto = require('crypto');
var fs = require('fs');
var merge = require('merge-stream');
var util = require('gulp-util');
var htmlmin = require('gulp-htmlmin');
const gulpif = require('gulp-if');

function js() {
  return src('js/*.js')
  .pipe(concat('app.js'))
  .pipe(uglify())
  .pipe(dest('../build'));
};
function css() {
    var base = src('css/*.css')
    .pipe(concat('app.css'))
    .pipe(uncss({
        ignore: [
          /* bootstrap UI */
          ".pull-right",
        ],
        html: ['../build/index.html']
    }))
    .pipe(minifyCSS({'keepSpecialComments': 0}))
    .pipe(dest('../build'));
    return base;
}
function static() {
    var base = src([
    'index.html'
  ]).pipe(htmlmin({collapseWhitespace: true, removeComments: true}))
  .pipe(dest('../build'));

  var imgs = src([
    'images/**/*'
  ]).pipe(dest('../build/images'));
  var fa = src([
    'fa/**/*'
  ]).pipe(dest('../build/fa'));
  var pub = src([
    'leaflet*.*',
  ]).pipe(dest('../build'));
  return merge(base, imgs, fa, pub);
}

function compress() {
    return src(['../build/**/*']).pipe(dest('../build'));
}

exports.default = series(
    parallel(series(static, css), js),
    compress
);
