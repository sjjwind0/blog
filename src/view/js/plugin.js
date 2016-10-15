function SetCwinHeight(){
	var iframeid=document.getElementById("plugin"); //iframe id
	if (document.getElementById) {
		if (iframeid && !window.opera) {
			if (iframeid.contentDocument && iframeid.contentDocument.body.offsetHeight) {
				 iframeid.height = iframeid.contentDocument.body.offsetHeight + 30;
			} else if (iframeid.Document && iframeid.Document.body.scrollHeight) {
				 iframeid.height = iframeid.Document.body.scrollHeight + 30;
			}
		}
	}
}