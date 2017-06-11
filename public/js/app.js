var app = {

    filesTimerMillis: 5000,
    files: undefined,
    filePath: undefined,

    toggleTree: function(e) {
        var target = e.target || e.srcElement;
        if (target.tagName != 'LI') {
            return;
        }
        else if (target.className == 'tree-file') {
            app.loadFile(target.getAttribute('data-file-path'));
        }
        else if (target.children && target.children.length > 0 && target.children[0].tagName == "UL") {
            target = target.children[0];
            target.classList.toggle('tree-expand');
            target.classList.toggle('tree-collapse');
        }
    },

    getTreeElement: function(file) {
        if (file.fileName[0] == '.') {
            return null;
        }
        var li = document.createElement('li');
        if (file.isDir) {
            li.className = 'tree-dir';
            li.innerText = file.fileName;
            var ul = document.createElement('ul');
            ul.className = 'tree-collapse';
            if (file.children && file.children.length > 0) {
                for (var i = 0; i < file.children.length; i++) {
                    var li2 = app.getTreeElement(file.children[i]);
                    if (li2) {
                        ul.append(li2);
                    }
                }
            }
            li.appendChild(ul);
        }
        else {
            li.className = 'tree-file';
            li.innerText = file.fileName;
            li.setAttribute('data-file-path', file.filePath);
        }
        return li;
    },

    saveFile: function() {
        var request = new XMLHttpRequest();
        request.onload = function () {
            if (this.status >= 200 && this.status < 400) {
                console.log("Saved.");
            }
            else {
                console.log('Error loading file.');
            }
        };
        request.open('PUT', '/api/file', true);
        var data = 'fp=' + encodeURIComponent(app.filePath);
		var srcDir = app.getParameterByName('src');
		if (srcDir && srcDir.length > 0) {
			data += "&src=" + encodeURIComponent(srcDir);
		}
		data += '&contents=' + encodeURIComponent(editor.getValue());
        request.send(data);
    },

    loadFile: function(filePath) {
        var request = new XMLHttpRequest();
        request.onload = function () {
            if (this.status >= 200 && this.status < 400) {
                app.filePath = filePath;
                ext = filePath.substring(filePath.lastIndexOf('.')+1);
                editor.setValue(this.responseText);
                monaco.editor.setModelLanguage(editor.getModel(), getLanguageForExtension(ext));
            }
            else {
                console.log('Error loading file.');
            }
        };
        var path = '/api/file';
		var srcDir = app.getParameterByName('src');
		if (srcDir && srcDir.length > 0) {
			path += "?src=" + encodeURIComponent(srcDir) + '&';
		}
		else {
		    path += '?';
        }
        path += 'fp=' + encodeURIComponent(filePath);
        request.open('GET', path, true);
        request.send();
    },

    loadFiles: function () {
        var request = new XMLHttpRequest();
        request.onload = function () {
            if (this.status >= 200 && this.status < 400) {
                var fileResponse = JSON.parse(this.responseText);
                var treeRoot = document.getElementById('tree-root');
                while (treeRoot.firstChild) {
                    treeRoot.removeChild(treeRoot.firstChild);
                }
                var li = document.createElement('li');
                li.innerText = fileResponse.fileName;
                var ul = document.createElement('ul');
                ul.className = 'tree-expand';
                ul.addEventListener('click', app.toggleTree);
                for (var i=0; i<fileResponse.children.length; i++) {
                    var li2 = app.getTreeElement(fileResponse.children[i])
                    if (li2) {
                        ul.append(li2);
                    }
                }
                li.appendChild(ul);
                treeRoot.appendChild(li);
            }
            else {
                console.log('Error loading files.');
                setTimeout(app.loadFiles(), app.filesTimerMillis);
            }
        };
        var path = '/api/files';
        var srcDir = app.getParameterByName('src');
        if (srcDir && srcDir.length > 0) {
			path += "?src=" + encodeURIComponent(srcDir);
        }
        request.open('GET', path, true);
        request.send();
    },

	getParameterByName: function(name, url) {
		if (!url) {
		    url = window.location.href;
		}
		name = name.replace(/[\[\]]/g, "\\$&");
		var regex = new RegExp("[?&]" + name + "(=([^&#]*)|&|#|$)");
        var results = regex.exec(url);
		if (!results) {
		    return null;
		}
		if (!results[2]) {
		    return '';
		}
		return decodeURIComponent(results[2].replace(/\+/g, " "));
	},

    init: function () {
        document.getElementById('save-btn').addEventListener('click', function() {
            app.saveFile();
        });
        setTimeout(app.loadFiles(), 1);
    }

};

(function () {
    app.init();
})();