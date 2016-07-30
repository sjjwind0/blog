window.onload = function() {
	$(".ds-qq").click(function() {
		Account.loginByQQ();
	});

	$(".btn-send").click(function() {
		var content = $(".wrap-text-f").val();
		Talk.sendTalk(Blog.getBlogId(), -1, content, function() {
			console.log("评论成功");
			$(".wrap-text-f").val("");
			setReply();
		});
	});


	var setReply = function() {
		$(".evt-reply").click(function() {
			var element = $(this).parent().parent().parent().parent().children(".module-cmt-box");
			if (element.css("display") == "none") {
				element.css("display", "block");
				$(this).text("取消回复");
			} else {
				element.css("display", "none");
				$(this).text("回复");
			}
		});

		$(".btn-reply").click(function() {
			var parentCommentId = $(this).parents(".comment-node");
			if (parentCommentId == null) {
				var id = -1;
			} else {
				var id = parseInt(parentCommentId.attr("cid"));
			}
			var content = $(this).parents(".post-wrap-w").children(".wrap-area-w").children(".area-textarea-w").children().val()
			Talk.sendTalk(Blog.getBlogId(), id, content, function() {
				setReply();
				$(".module-cmt-box").css("display", "none");
				$(".textarea-fw, .textarea-bf").val("");
				$(".wrap-text-f").val("");
			});
		});
	}
	setReply();
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
// Create Base64 Object
var Base64={_keyStr:"ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/=",encode:function(e){var t="";var n,r,i,s,o,u,a;var f=0;e=Base64._utf8_encode(e);while(f<e.length){n=e.charCodeAt(f++);r=e.charCodeAt(f++);i=e.charCodeAt(f++);s=n>>2;o=(n&3)<<4|r>>4;u=(r&15)<<2|i>>6;a=i&63;if(isNaN(r)){u=a=64}else if(isNaN(i)){a=64}t=t+this._keyStr.charAt(s)+this._keyStr.charAt(o)+this._keyStr.charAt(u)+this._keyStr.charAt(a)}return t},decode:function(e){var t="";var n,r,i;var s,o,u,a;var f=0;e=e.replace(/[^A-Za-z0-9+/=]/g,"");while(f<e.length){s=this._keyStr.indexOf(e.charAt(f++));o=this._keyStr.indexOf(e.charAt(f++));u=this._keyStr.indexOf(e.charAt(f++));a=this._keyStr.indexOf(e.charAt(f++));n=s<<2|o>>4;r=(o&15)<<4|u>>2;i=(u&3)<<6|a;t=t+String.fromCharCode(n);if(u!=64){t=t+String.fromCharCode(r)}if(a!=64){t=t+String.fromCharCode(i)}}t=Base64._utf8_decode(t);return t},_utf8_encode:function(e){e=e.replace(/rn/g,"n");var t="";for(var n=0;n<e.length;n++){var r=e.charCodeAt(n);if(r<128){t+=String.fromCharCode(r)}else if(r>127&&r<2048){t+=String.fromCharCode(r>>6|192);t+=String.fromCharCode(r&63|128)}else{t+=String.fromCharCode(r>>12|224);t+=String.fromCharCode(r>>6&63|128);t+=String.fromCharCode(r&63|128)}}return t},_utf8_decode:function(e){var t="";var n=0;var r=c1=c2=0;while(n<e.length){r=e.charCodeAt(n);if(r<128){t+=String.fromCharCode(r);n++}else if(r>191&&r<224){c2=e.charCodeAt(n+1);t+=String.fromCharCode((r&31)<<6|c2&63);n+=2}else{c2=e.charCodeAt(n+1);c3=e.charCodeAt(n+2);t+=String.fromCharCode((r&15)<<12|(c2&63)<<6|c3&63);n+=3}}return t}}

var Talk = Talk || {}

Talk.sendTalk = function (blogId, commentId, content, callback) {
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
				var m = $(".list-newest-b").children();
				$(m[0]).after(Base64.decode(result.data.comment))
				if (callback != null) {
					callback();
				}
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

Account.waitingLogin = function () {
	// 创建长连接
	var listener = function(event) {
		clearInterval(Account.timer);
		if (event.data == "login") {
			// 登录成功
			var url = "http://blog.windy.live/api";
			$.ajax({
				url: url,
				type: "POST",
				data: JSON.stringify({"type": "getUserInfo"}),
				contentType: "application/json; charset=utf-8",
				dataType: "json",
				success: function(data) {
					if (data.code == 0) {
						var nickName = data.data.name;
						var pic = data.data.pic;
						$(".ds-login-buttons").css("display", "none");
						$(".no-user-name-login").css("display", "none");
						$(".user-name-login").css("display", "block");
						$(".user-name-login + .wrap-name-nick").html(nickName);
						$(".userPic").attr("src", pic);
						window.removeEventListener('login', this);
					} else {
						alert("login failed: " + data.msg);
					}
				}
			});
		}
	};
	window.addEventListener('message', listener, false);
}

Account.loginByQQ = function(clientId, redirectUri) {
	Account.waitingLogin();
	var url = "https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=101324961&redirect_uri=http%3A%2F%2Fblog.windy.live%2Flogin%3Ftype%3Dqq&scope=get_user_info";
	var child = window.open(url);
	Account.timer = setInterval(function() {
		  var message = "helo";
			child.postMessage(message, "/");
	}, 200);
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
