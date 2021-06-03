var gulp       = require('gulp');
var livereload = require('gulp-livereload');
var less       = require('gulp-less');
var path       = require('path');
var cleanCSS   = require('gulp-clean-css');

gulp.task('watch', function()
{
	livereload.listen();

	gulp.watch(['static/css/*.css', 'static/js/*.js', 'app/templates/*.html', 'app/templates/admin/*.html', 'static/less/*.less']).on('change', function(e)
	{
		return gulp.src(e.path)
			.pipe(livereload());
	});

	gulp.watch('static/less/*.less').on('change', function(e)
	{
		return gulp.src(e.path)
			.pipe(less())
			.pipe(cleanCSS())
			.pipe(gulp.dest('static/css'))
	});
});