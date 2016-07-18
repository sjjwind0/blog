package info

const (
	AccountTypeQQ    = iota
	AccountTypeWeibo = iota
)

type UserInfo struct {
	UserID          int64
	UserOpenID      string
	UserAccountType int
	UserName        string
	Sex             string
	SmallFigureurl  string
	BigFigureurl    string
	LastLoginTime   int64
	RegisterTime    int64
	CurrentIP       string
}
