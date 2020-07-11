const { series, parallel, src, dest } = require('gulp');

var uncss = require('gulp-uncss');
var minifyCSS = require('gulp-minify-css');
var concat = require('gulp-concat');
var uglify = require('gulp-uglify');
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
function html() {
  var base = src([
    'index.html'
  ]).pipe(htmlmin({collapseWhitespace: true, removeComments: true}))
  .pipe(dest('../build'));
  return base;
}
function static() {
  var base = src([
    'static/**/*'
  ]).pipe(dest('../build'));
  return base;
}

function compress() {
    return src(['../build/**/*']).pipe(dest('../build'));
}

exports.default = series(
    parallel(series(html, css), static, js),
    compress
);
