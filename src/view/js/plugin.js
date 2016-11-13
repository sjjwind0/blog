function SetCwinHeight(){
	var iframeid = document.getElementById("plugin"); //iframe id
	if (document.getElementById) {
		if (iframeid && !window.opera) {
			if (iframeid.contentDocument && iframeid.contentDocument.body.offsetHeight) {
				 iframeid.height = iframeid.contentDocument.body.offsetHeight + 30;
			} else if (iframeid.Document && iframeid.Document.body.scrollHeight) {
				 iframeid.height = iframeid.document.body.scrollHeight + 30;
			}
		}
	}
}

function SetCwinHeightAndBase(baseURL) {
	var iframeid = document.getElementById("plugin"); //iframe id
	if (document.getElementById) {
		if (iframeid && !window.opera) {
			if (iframeid.contentDocument && iframeid.contentDocument.body.offsetHeight) {
				 iframeid.height = iframeid.contentDocument.body.offsetHeight + 30;
			} else if (iframeid.Document && iframeid.Document.body.scrollHeight) {
				 iframeid.height = iframeid.document.body.scrollHeight + 30;
			}
		}
	}
	var pElement = document.createElement("base");
	pElement.href = "baseURL";
	iframeid.contentWindow.document.head.append(pElement)
}