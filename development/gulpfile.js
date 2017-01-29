"use strict";

var browserify = require('browserify');
var gulp = require('gulp');
var autoprefixer = require("gulp-autoprefixer");
var plumber = require("gulp-plumber");
var rename = require('gulp-rename');
var sass = require("gulp-sass");
var uglify = require("gulp-uglify");
var reactify = require('reactify');
var through  = require('through2');
var buffer = require('vinyl-buffer');

gulp.task("react", function() {
  var browserified = through.obj(function(file, enc, next) {
      browserify(file.path)
          .transform(reactify)
          .bundle(function(err, res) {
              file.contents = res;
              next(null, file);
          }
      );
  });
  gulp.src('/monitor/app/assets/js/**/*.jsx')
      .pipe(browserified)
      .pipe(rename({extname: '.js'}))
      .pipe(uglify())
      .pipe(gulp.dest("/monitor/app/assets/js/min"));
});

gulp.task("sass", function() {
  gulp.src("/monitor/app/assets/scss/**/*.scss")
      .pipe(plumber())
      .pipe(sass())
      .pipe(autoprefixer())
      .pipe(gulp.dest("/monitor/app/assets/css"));
});

gulp.task("default", function() {
    gulp.watch("/monitor/app/assets/js/**/*.jsx", ["react"]);
    gulp.watch("/monitor/app/assets/scss/**/*.scss", ["sass"]);
});
