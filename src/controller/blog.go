package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"framework/response"
	"html/template"
	"info"
	"log"
	"net/http"
	"strconv"
	"time"
)

type commentRender struct {
	UserID         string
	UserName       string
	CommentContent string
	CommentTime    string
	ChildContent   template.HTML
}

const (
	firstComment = `<div class="clear-g block-cont-gw block-cont-bg" datatype="time" cmtnum="1" cmtid="957805360"> 
	<div class="cont-head-gw"> 
		<div class="head-img-gw">
			<a href="javascript:void(0)" commhref="http://www.douban.com/people/94809862/"><img src="http://qiniu.cuiqingcai.com/wp-content/uploads/2015/05/20150525111154.jpg" onerror="SOHUCS.isImgErr(this)" width="42" height="42" alt="" uid="{{.UserID}}"></a>
        </div> 
	</div> 
	<div class="cont-msg-gw"> 
		<div class="msg-wrap-gw"> 
			<div class="wrap-user-gw global-clear-spacing"> 
				<span class="user-time-gw user-time-bg evt-time">{{.CommentTime}}</span> 
				<span class="user-name-gw" title="{{.UserName}}"><a href="javascript:void(0)" commhref="http://www.douban.com/people/94809862/" uid="-1990645212">{{.UserName}}</a></span> 
			</div> 
			{{.ChildContent}}
			<div class="wrap-issue-gw"> 
				<p class="issue-wrap-gw"> <span class="wrap-word-bg ">{{.CommentContent}}</span> </p> 
			</div> 
		<div class="clear-g wrap-action-gw"> 
			<div class="action-click-gw"> 
				<i class="gap-gw"></i> 
				<span class="click-ding-gw"><a href="javascript:void(0)" title="顶" class="evt-support"><i class="icon-gw icon-ding-bg"></i><em class="icon-name-bg"></em></a></span>
				<i class="gap-gw"></i>
				<span class="click-cai-gw"><a href="javascript:void(0)" title="踩" class="evt-opposed"><i class="icon-gw icon-cai-bg"></i><em class="icon-name-bg"></em></a></span>
				<i class="gap-gw"></i>
				<span class="click-reply-gw click-reply-eg"><a href="javascript:void(0)" class="evt-reply">回复</a></span>
				<i class="gap-gw"></i>
				<span class="click-share-gw click-reply-eg"><a href="javascript:void(0)" class="evt-share">分享</a></span> 
			</div> 
			<div class="action-from-gw action-from-bg"></div> 
			</div> 
			<div class="wrap-reply-gw"></div> 
		</div> 
	</div> 
</div>`
	secondComment = `<div class="wrap-build-gw"> 
	<div class="build-floor-gw"> 
		<div class="build-msg-gw borderbot" cmtid="957805360" floornum="1"> 
			<div class="wrap-user-gw global-clear-spacing"> 
				<span class="user-time-gw user-time-bg user-floor-gw">1</span> 
				<span class="user-name-gw"><a href="javascript:void(0)" commhref="http://www.douban.com/people/94809862/" title="{{.UserName}}" uid="{{.UserID}}">{{.UserName}}</a></span> 
			</div> 
			{{.ChildContent}}
			<div class="wrap-issue-gw"> 
				<p class="issue-wrap-gw"> <span class="wrap-word-bg ">{{.CommentContent}}</span> </p> 
			</div> 
			<div class="clear-g wrap-action-gw evt-active-wrapper" style="visibility: hidden;"> 
				<div class="action-click-gw"> 
					<i class="gap-gw"></i> 
					<span class="click-ding-gw"><a href="javascript:void(0)" title="顶" class="evt-support"><i class="icon-gw icon-ding-bg"></i><em class="icon-name-bg"></em></a></span> 
					<i class="gap-gw"></i> 
					<span class="click-cai-gw"><a href="javascript:void(0)" title="踩" class="evt-opposed"><i class="icon-gw icon-cai-bg"></i><em class="icon-name-bg"></em></a></span> 
					<i class="gap-gw"></i> 
					<span class="click-reply-gw click-reply-eg"><a href="javascript:void(0)" class="evt-reply">回复</a></span> 
					<i class="gap-gw"></i> 
					<span class="click-share-gw click-reply-eg"><a href="javascript:void(0)" class="evt-share">分享</a></span> 
				</div> 
				<div class="action-from-gw action-from-bg"></div> 
			</div> 
			<div class="wrap-reply-gw"> 
			</div> 
		</div> 
	</div> 
</div>`
)

type BlogController struct {
}

func NewBlogController() *BlogController {
	return &BlogController{}
}

func (this *BlogController) MockComment() string {
	index := 0
	creator := func(parentId int) info.CommentInfo {
		var comment info.CommentInfo
		comment.CommentID = index
		comment.BlogID = 1
		comment.UserID = 1
		comment.ParentCommentID = parentId
		comment.Content = "哈哈哈，评论"
		comment.Time = time.Now().Unix()
		comment.Praise = 0
		comment.Address = "火星上"
		index++
		return comment
	}

	buildCommentRender := func(info *info.CommentInfo, childComment *string) commentRender {
		var render commentRender
		render.ChildContent = template.HTML(*childComment)
		render.CommentContent = info.Content
		render.CommentTime = "2016年12月1日 03:22"
		render.UserID = string(info.UserID)
		render.UserName = "测试"
		return render
	}

	buildParentComment := func(childComment *string, parent *info.CommentInfo) string {
		t, err := template.New("test").Parse(firstComment)
		b := bytes.NewBuffer(make([]byte, 0))
		strIO := bufio.NewWriter(b)
		if err == nil {
			t.Execute(strIO, buildCommentRender(parent, childComment))
		} else {
			fmt.Println(err)
		}
		strIO.Flush()
		return string(b.Bytes())
	}

	buildComment := func(childComment *string, parent *info.CommentInfo) string {
		t, err := template.New("test").Parse(secondComment)
		b := bytes.NewBuffer(make([]byte, 0))
		strIO := bufio.NewWriter(b)
		if err == nil {
			t.Execute(strIO, buildCommentRender(parent, childComment))
		} else {
			fmt.Println(err)
		}
		strIO.Flush()
		return string(b.Bytes())
	}

	fmt.Println(buildComment)

	var commentList []info.CommentInfo = make([]info.CommentInfo, 5)
	rawComment := ""
	for i := 0; i < 5; i++ {
		index += 1
		commentList[i].CommentID = i
		commentList[i].BlogID = 1
		commentList[i].UserID = 1
		commentList[i].ParentCommentID = -1
		commentList[i].Content = fmt.Sprintf("哈哈哈，评论%d", i)
		commentList[i].Time = time.Now().Unix()
		commentList[i].Praise = 0
		commentList[i].Address = "火星上"
		old := ""
		childComment := creator(i)
		childChildComment := creator(childComment.CommentID)
		rawChildChildComment := buildComment(&old, &childChildComment)
		rawChildComment := buildComment(&rawChildChildComment, &childComment)
		rawComment += buildParentComment(&rawChildComment, &commentList[i])
	}
	return rawComment
}

func (this *BlogController) HandlerAction(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	fmt.Println(r.Form["id"])
	_, err := strconv.Atoi(r.Form["id"][0])
	if err != nil {
		response.JsonResponseWithMsg(w, response.ErrorParamError, "param error")
		return
	}
	t, err := template.ParseFiles("./src/view/html/blog.html")
	if err != nil {
		log.Println(err)
	}
	t.Execute(w, template.HTML(this.MockComment()))
}
