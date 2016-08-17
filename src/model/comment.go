package model

import (
	"container/list"
	"fmt"
	"framework/database"
	"info"
	"sync"
	"time"
)

const (
	kCommentTableName = "comment"
	kCommentBlogId    = "blog_id"
	kCommentId        = "comment_id"
	kCommentParentId  = "parent_id"
	kCommentUserId    = "user_id"
	kCommentContent   = "content"
	kCommentTime      = "time"
	kCommentPraise    = "praise"
	kCommentDissent   = "dissent"
	kCommentAddress   = "address"
)

type commentModel struct {
}

var commentModelInstance *commentModel = nil

var commentOnce sync.Once

func ShareCommentModel() *commentModel {
	commentOnce.Do(func() {
		commentModelInstance = &commentModel{}
	})
	return commentModelInstance
}

func (c *commentModel) CreateTable() error {
	if database.DatabaseInstance().DoesTableExist(kCommentTableName) {
		return nil
	}
	sql := fmt.Sprintf(`
	CREATE TABLE %s (
		%s int(32) unsigned NOT NULL AUTO_INCREMENT,
		%s int(32) NOT NULL,
		%s int(32) NOT NULL DEFAULT '-1',
		%s int(32) NOT NULL,
		%s varchar(1024) NOT NULL,
		%s int(64) NULL DEFAULT '0',
		%s int(32) NULL DEFAULT '0',
		%s int(32) NULL DEFAULT '0',
		%s varchar(1024) DEFAULT '',
		PRIMARY KEY (%s)
	) CHARSET=utf8;`, kCommentTableName, kCommentId,
		kCommentBlogId, kCommentParentId, kCommentUserId, kCommentContent, kCommentTime,
		kCommentPraise, kCommentDissent, kCommentAddress, kCommentId)
	_, err := database.DatabaseInstance().DB.Exec(sql)
	return err
}

func (c *commentModel) AddComment(userId int, blogId int, commentId int, commentContent string) (int, error) {
	sql := fmt.Sprintf("insert into %s(%s, %s, %s, %s, %s) values(?, ?, ?, ?, ?)",
		kCommentTableName, kCommentUserId, kCommentBlogId, kCommentParentId,
		kCommentContent, kCommentTime)
	stat, err := database.DatabaseInstance().DB.Prepare(sql)
	if err == nil {
		defer stat.Close()
		result, err := stat.Exec(userId, blogId, commentId, commentContent, time.Now().Unix())
		if err == nil {
			insertId, err := result.LastInsertId()
			return int(insertId), err
		}
	}
	return 0, err
}

func (c *commentModel) DeleteAllBlogComment(blogId int) error {
	sql := fmt.Sprintf("delete from %s where %s = ?", kCommentTableName, kCommentId)
	_, err := database.DatabaseInstance().DB.Exec(sql, blogId)
	return err
}

func (c *commentModel) FetchCommentByCommentId(commentId int) (*info.CommentInfo, error) {
	sql := fmt.Sprintf("select * from %s where %s = ?", kCommentTableName, kCommentId)
	rows, err := database.DatabaseInstance().DB.Query(sql, commentId)
	if err == nil {
		defer rows.Close()
		if rows.Next() {
			var commentInfo info.CommentInfo
			err = rows.Scan(&commentInfo.CommentID, &commentInfo.BlogID, &commentInfo.ParentCommentID,
				&commentInfo.UserID, &commentInfo.Content, &commentInfo.Time,
				&commentInfo.Praise, &commentInfo.Dissent, &commentInfo.Address)
			return &commentInfo, nil
		}
	}
	return nil, err
}

func (c *commentModel) FetchAllCommentByBlogId(blogId int) (*list.List, error) {
	sql := fmt.Sprintf("select * from %s where %s = ? order by %s desc", kCommentTableName, kCommentBlogId, kCommentId)
	rows, err := database.DatabaseInstance().DB.Query(sql, blogId)
	var blogList *list.List = list.New()
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var commentInfo info.CommentInfo
			err = rows.Scan(&commentInfo.CommentID, &commentInfo.BlogID, &commentInfo.ParentCommentID,
				&commentInfo.UserID, &commentInfo.Content, &commentInfo.Time,
				&commentInfo.Praise, &commentInfo.Dissent, &commentInfo.Address)
			if err == nil {
				blogList.PushBack(commentInfo)
			} else {
				return nil, err
			}
		}
	}
	return blogList, err
}

func (b *commentModel) FetchCommentCount(blogId int) (int, error) {
	sql := fmt.Sprintf("select count(*) from %s where %s = ?", kCommentTableName, kCommentBlogId)
	rows, err := database.DatabaseInstance().DB.Query(sql, blogId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var count int
			err = rows.Scan(&count)
			if err == nil {
				return count, nil
			}
		}
	}
	return 0, err
}

func (b *commentModel) FetchCommentPeopleCount(blogId int) (int, error) {
	sql := fmt.Sprintf("select count(distinct(%s)) from %s where %s = ?",
		kCommentUserId, kCommentTableName, kCommentBlogId)
	rows, err := database.DatabaseInstance().DB.Query(sql, blogId)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var count int64
			err = rows.Scan(&count)
			if err == nil {
				return int(count), nil
			}
		}
	}
	return 0, err
}
