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

type BlogInfo struct {
	BlogID           int
	BlogUUID         string
	BlogTitle        string
	BlogSortType     string
	BlogTagList      []string
	BlogTime         int64
	BlogVisitCount   int
	BlogPraiseCount  int
	BlogDissentCount int
}
