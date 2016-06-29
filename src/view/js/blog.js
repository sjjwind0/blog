console.log("1111");

window.onload = function() {
	$(".btn-fw + .btn-bf + .single-btn-bf").click(function() {
		var content = $(".wrap-text-f").val();
		Talk.sendTalk(1, 1, -1, content);
	});

	$(".evt-reply").click(function() {
		console.log(this);
		console.log($(this));
		var element = $(this).parent().parent().parent().parent().children(".module-cmt-box");
		if (element.css("display") == "none") {
			element.css("display", "block");
			$(this).text("取消回复");
		} else {
			element.css("display", "none");
			$(this).text("回复");
		}
	});
	$(".btn-fw").click(function() {
		var parentCommentId = $(this).parents(".comment-node");
		if (parentCommentId == null) {
			var id = -1;
		} else {
			var id = parseInt(parentCommentId.attr("cid"));
		}
		var content = $(this).parents(".post-wrap-w").children(".wrap-area-w").children(".area-textarea-w").children().val()
		Talk.sendTalk(1, 1, id, content);
	});
	var allHref = $(window.frames["blog"].document).find("a");
	for (var i = 0; i < allHref.length; i++) {
		var url = $(allHref[i]).attr("href");
		if (typeof url != 'undefined') {
			if (url[0] != "#") {
				$(allHref[i]).attr("target", "_parent");
			}
		}
	}
}

var Talk = Talk || {}

Talk.sendTalk = function (userId, blogId, commentId, content) {
	var url = "../api";
	var content = {
		"type": "talk",
		"userId": userId,
		"blogId": blogId,
		"commentId": commentId,
		"content": content
	}
	$.ajax({
		url: url,
		type: "POST",
		data: JSON.stringify(content),
		contentType: "application/json; charset=utf-8",
		dataType: "json",
		success: function(result) {
			if (result.code == 0) {
				console.log("send talk success");
			}
			console.log("failed");
		}
	});
}