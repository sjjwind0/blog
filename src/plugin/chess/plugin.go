package chess

import (
	"plugin/chess/impl"
	"time"
)

type chessPlugin struct {
}

func NewChessPlugin() *chessPlugin {
	return &chessPlugin{}
}

func (c *chessPlugin) GetPluginName() string {
	return "基于alpha-beta算法的中国象棋"
}

func (c *chessPlugin) GetPluginUUID() string {
	return "16b99f47-669a-11e6-9bc1-7831c1c81ccc"
}

func (c *chessPlugin) GetPluginCoverURL() string {
	return "plugin/chess/res/img/plugin-chess-cover.jpg"
}

func (c *chessPlugin) GetPluginDisplayPath() string {
	return "plugin/chess"
}

func (c *chessPlugin) GetPluginDescription() string {
	return "棋类游戏在桌面游戏中已经非常成熟，中国象棋的版本也非常多。今天这款基于HTML5技术的中国象棋游戏非常有特色，我们不仅可以选择中国象棋的游戏难度，而且可以切换棋盘的样式。程序写累了，喝上一杯咖啡，和电脑对弈几把吧，相信这HTML5中国象棋游戏的实现算法你比较清楚，可以打开源码来研究一下这款中国象棋游戏。"
}

func (c *chessPlugin) GetPluginDownloadURL() string {
	return "plugin/chess/download"
}

func (c *chessPlugin) GetPluginPublicTime() int64 {
	return time.Now().Unix()
}

func (c *chessPlugin) ResourceHandler() []string {
	return []string{
		"/plugin/chess/res/js",
		"/plugin/chess/res/img",
		"/plugin/chess/res/css",
	}
}

func (c *chessPlugin) NormalHanlder() []interface{} {
	return []interface{}{impl.NewHomeController()}
}

func (c *chessPlugin) WebSocketHandler() []interface{} {
	return []interface{}{impl.NewMessageController()}
}
