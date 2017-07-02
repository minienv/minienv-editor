var editor = undefined;
require(['vs/editor/editor.main'], function() {
	editor = monaco.editor.create(document.getElementById('editor'), {
		theme: 'vs',
		automaticLayout: true
	});
});

var getLanguageForExtension = function(ext) {
	var language = "text";
    if (ext == "css") {
        language = "css";
    }
    else if (ext == "htm" || ext == "html" || ext == "xhtml") {
        language = "html";
    }
    else if (ext == "js") {
        language = "javascript";
    }
    else if (ext == "json") {
        language = "json";
    }
    else if (ext == "md") {
        language = "markdown";
    }
    else if (ext == "swift") {
        language = "swift";
    }
    else if (ext == "yml" || ext == "yaml") {
        language = "yaml";
    }
    return language;
};