var editor = undefined;
require(['vs/editor/editor.main'], function() {
	editor = monaco.editor.create(document.getElementById('editor'), {
		theme: 'vs',
		automaticLayout: true
	});
});

var getLanguageForExtension = function(ext) {
	var language = "plaintext";
	if (ext == "cpp") {
		language = "cpp";
	}
	else if (ext == "css") {
        language = "css";
    }
	else if (ext == "go") {
		language = "go";
	}
    else if (ext == "htm" || ext == "html" || ext == "xhtml") {
        language = "html";
    }
	else if (ext == "java") {
		language = "java";
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
	else if (ext == "php") {
		language = "php";
	}
	else if (ext == "py") {
		language = "python";
	}
	else if (ext == "rb") {
		language = "ruby";
	}
	else if (ext == "sql") {
		language = "sql";
	}
    else if (ext == "swift") {
        language = "swift";
    }
	else if (ext == "ts") {
		language = "typescript";
	}
	else if (ext == "xml") {
		language = "xml";
	}
    else if (ext == "yml" || ext == "yaml") {
        language = "yaml";
    }
    return language;
};