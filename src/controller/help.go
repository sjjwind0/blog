package controller

import (
	"time"
)

func FormatTime(t int64) string {
	publishTime := time.Unix(t, 0)
	duration := time.Now().Sub(publishTime)
	if duration < time.Minute {
		return "刚刚"
	} else if duration < time.Minute*10 {
		return "10分钟以内"
	} else {
		if duration < time.Hour {
			return "一个小时以内"
		} else if duration < time.Hour*24 {
			return "一天内"
		} else if duration < time.Hour*48 {
			return "一天前"
		} else {
			return publishTime.Format("2006年01月02日")
		}
	}
}

func FormatRealTime(t int64) string {
	publishTime := time.Unix(t, 0)
	return publishTime.Format("2006年01月02日")
}

type BlogTime struct {
	Tag  string
	Time int64
}

type BlogTimeList []*BlogTime

func (t BlogTimeList) Len() int {
	return len(t)
}

func (t BlogTimeList) Swap(tag1, tag2 int) {
	t[tag1], t[tag2] = t[tag2], t[tag1]
}

func (t BlogTimeList) Less(tag1, tag2 int) bool {
	return t[tag1].Time < t[tag2].Time
}
