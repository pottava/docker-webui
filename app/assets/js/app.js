var app = {};

(function($) {
	app.func = {
		ajax : _ajax,
		trim : _trim,
		query : _query,
		byteFormat: _byteFormat,
		relativeTime: _relativeTime
	};

	function _ajax(arg) {
		var dt = arg.dataType ? arg.dataType : 'json';
		$.ajax({
			url: arg.url, type: arg.type,
			data: arg.data, dataType: dt,
			success: function (data) {
				arg.success && arg.success(data);
			},
			error: function(xhr, status, err) {
				if (arg.error) {
					arg.error(xhr, status, err);
					return;
				}
				console.error(arg.url, status, err.toString());
			}
		});
	}
	function _trim(value) {
		return value.replace(/\s/g,'').replace(/　/g,'')
	}
	function _query(key, def) {
		key = key.replace(/[\[]/, "\\[").replace(/[\]]/, "\\]");
		def = def ? def : "";
		var regex = new RegExp("[\\?&]" + key + "=([^&#]*)"),
				results = regex.exec(location.search);
		return results === null ? def : decodeURIComponent(results[1].replace(/\+/g, " "));
	}
	function _byteFormat(bytes) {
		var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB'];
		if (bytes == 0) return '0 Byte';
		var i = parseInt(Math.floor(Math.log(bytes) / Math.log(1024)));
		return (bytes / Math.pow(1024, i)).toFixed(1) + ' ' + sizes[i];
	}
	function _relativeTime(date){
		var now = new Date().getTime(),
				offset = Math.ceil((now - date.getTime()) / 1000),
				message;
		if (offset < 60) {
			message = offset + ' second' + (offset == 1 ? '' : 's') + ' ago';
		} else if (offset < (60*60)) {
			var candidate = Math.floor(offset / 60);
			message = candidate + ' minute' + (candidate == 1 ? '' : 's') + ' ago';
		} else if (offset < (24*60*60)) {
			var candidate = Math.floor(offset / 3600);
			message = candidate + ' hour' + (candidate == 1 ? '' : 's') + ' ago';
		} else if (offset < (7*24*60*60)) {
			var candidate = Math.floor(offset / 86400);
			message = candidate + ' day' + (candidate == 1 ? '' : 's') + ' ago';
		} else {
			var candidate = Math.floor(offset / 604800);
			message = candidate + ' week' + (candidate == 1 ? '' : 's') + ' ago';
		}
		return message;
	}
})(jQuery);
