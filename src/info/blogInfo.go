package info

type CommentInfo struct {
	CommentID       int
	BlogID          int
	UserID          int
	ParentCommentID int
	Content         string
	Time            int64
	Praise          int
	Dissent         int
	Address         string
}
