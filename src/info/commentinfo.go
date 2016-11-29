package info

const (
	CommentType_Blog   = 0
	CommentType_Plugin = 1
)

type CommentInfo struct {
	CommentID       int
	Type            int
	TypeID          int
	UserID          int64
	ParentCommentID int
	Content         string
	Time            int64
	Praise          int
	Dissent         int
	Address         string
}
