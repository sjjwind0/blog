package controller

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"info"
	"model"
	"strconv"
)

func buildCommentRender(info *info.CommentInfo, childComment *string,
	floor *int) apiCommentRender {
	var render apiCommentRender
	render.ChildContent = template.HTML(*childComment)
	render.CommentContent = info.Content
	render.CommentTime = FormatTime(info.Time)
	render.CommentID = strconv.Itoa(info.CommentID)
	render.UserID = string(info.UserID)
	render.Floor = *floor
	userInfo, err := model.ShareUserModel().GetUserInfoById(info.UserID)
	if err == nil {
		render.User = userInfo
	}
	(*floor)++
	return render
}

func buildCommentString(child *string, info *info.CommentInfo,
	step int, floor *int) string {
	var tmpl string = ""
	if step == 0 {
		tmpl = firstComment
	} else {
		tmpl = secondComment
	}
	t, err := template.New("comment").Parse(tmpl)
	buf := bytes.NewBuffer(make([]byte, 0))
	strIO := bufio.NewWriter(buf)
	if err == nil {
		t.Execute(strIO, buildCommentRender(info, child, floor))
	} else {
		fmt.Println(err)
	}
	strIO.Flush()
	return string(buf.Bytes())
}

func buildOneCommentFromCommentList(commentList *[]*info.CommentInfo) string {
	step := 0
	var floor int = 1
	var childCommentContent string = ""
	for i, commentInfo := range *commentList {
		// the first comment
		if i != len(*commentList)-1 {
			step = 1
		} else {
			step = 0
		}
		currentCommentContent := buildCommentString(&childCommentContent, commentInfo, step, &floor)
		childCommentContent = currentCommentContent
	}
	return childCommentContent
}

func buildOneCommentFromCommentTree(commentTree *map[int]*info.CommentInfo,
	currentComment *info.CommentInfo) string {
	step := 0
	var floor int = 1
	return buildOneCommentFromCommentTreeRecursion(commentTree, currentComment, step, &floor)
}

func buildOneCommentFromCommentTreeRecursion(commentTree *map[int]*info.CommentInfo,
	currentComment *info.CommentInfo, step int, floor *int) string {
	var childComment = ""
	if currentComment.ParentCommentID == -1 {
		return buildCommentString(&childComment, currentComment, step, floor)
	}
	// 首先build 子元素
	childInfo := (*commentTree)[currentComment.ParentCommentID]
	childComment = buildOneCommentFromCommentTreeRecursion(commentTree, childInfo, step+1, floor)
	return buildCommentString(&childComment, currentComment, step, floor)
}
