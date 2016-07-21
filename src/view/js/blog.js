window.onload = function() {
	$(".ds-qq").click(function() {
		Account.loginByQQ();
	});

	$(".single-btn-bf").click(function() {
		var content = $(".wrap-text-f").val();
		Talk.sendTalk(Blog.getBlogId(), -1, content);
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
		Talk.sendTalk(Blog.getBlogId(), id, content);
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

Talk.sendTalk = function (blogId, commentId, content) {
	var url = "../api";
	var content = {
		"type": "talk",
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
			} else {
				alert("评论失败, err: " + result.msg)
			}
		}
	});
}

var Account = Account || {}

Account.loginInfo = {
	userName: "",
	sex: "",
	pic: "",
}

Account.loginByQQ = function() {
	var url = "https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=101324961&redirect_uri=http%3A%2F%2Fblog.windy.live%2Flogin%3Ftype%3Dqq&scope=get_user_info";
	window.open(url);
	// $.get("http://blog.windy.live/login?type=connect", function (data) {
	// 	console.log("data: " + data);
	// });
}

Account.loginByWeixin = function() {
	// not implement
}

Account.loginByWeibo = function() {
	// not implement
}

var Blog = Blog || {}

Blog.blogInfo = {
	blogId: "",
}

Blog.getBlogId = function() {
	return parseInt($(".blog").attr("bid"));
}

