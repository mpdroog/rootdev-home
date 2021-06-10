const { series, parallel, src, dest } = require('gulp');

const minifyCSS = require('gulp-minify-css');
const concat = require('gulp-concat');
const uglify = require('gulp-uglify');
const htmlmin = require('gulp-htmlmin');
const purgecss = require('gulp-purgecss');

function js() {
  return src('js/*.js')
  .pipe(concat('app.js'))
  .pipe(uglify())
  .pipe(dest('../build'));
};
function css() {
    var base = src('css/*.css')
    .pipe(concat('app.css'))
    .pipe(purgecss({
        ignore: [
          /* bootstrap UI */
          ".pull-right",
        ],
        content: ['../build/**/*.html']
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
    parallel(series(html, static, css), js),
    compress
);
